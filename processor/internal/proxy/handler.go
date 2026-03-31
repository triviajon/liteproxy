package proxy

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/triviajon/liteproxy/processor/internal/cache"
	"github.com/triviajon/liteproxy/processor/internal/constant"
	"github.com/triviajon/liteproxy/processor/internal/logging"
	"github.com/triviajon/liteproxy/processor/internal/rewritepipeline"
	"github.com/triviajon/liteproxy/processor/internal/util"
)

type ProxyServer struct {
	Pipeline rewritepipeline.Pipeline
	Cache    cache.Cache
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logging.Warnf("Method not allowed - method=%s path=%s", r.Method, r.RequestURI)
		http.Error(w, "Only GET supported", http.StatusMethodNotAllowed)
		return
	}

	rawTarget := r.URL.Query().Get("url")
	if rawTarget == "" {
		logging.Warnf("Missing url parameter - path=%s", r.RequestURI)
		http.Error(w, "Missing required query parameter: url", http.StatusBadRequest)
		return
	}

	logging.Debugf("Processing request - raw_target=%s", rawTarget)
	targetURL, err := url.Parse(rawTarget)
	if err != nil || targetURL.Scheme == "" || targetURL.Host == "" {
		logging.Warnf("Invalid URL - raw_target=%s error=%v", rawTarget, err)
		http.Error(w, "Invalid url parameter", http.StatusBadRequest)
		return
	}

	if err := p.serveFromCache(w, r, *targetURL); err == nil {
		logging.Debugf("Served from cache - target=%s", targetURL.String())
		return
	}

	logging.Debugf("Cache miss, proxying upstream - target=%s", targetURL.String())
	modifyResponseFn, err := util.Bind1(p.modifyResponse, *targetURL)
	if err != nil {
		logging.Errorf("Error binding modify response - error=%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(targetURL)
			pr.Out.Host = targetURL.Host
			pr.Out.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
			pr.Out.Header.Set("X-Proxy-Processor", "liteproxy")
			pr.Out.Header.Del("Accept-Encoding")
		},
		ModifyResponse: modifyResponseFn,
	}

	logging.Debugf("Starting reverse proxy - target=%s", targetURL.String())
	proxy.ServeHTTP(w, r)
	logging.Debugf("Reverse proxy completed - target=%s", targetURL.String())
}

func (p *ProxyServer) serveFromCache(w http.ResponseWriter, r *http.Request, targetURL url.URL) error {
	logging.Debugf("Attempting cache retrieval - target=%s", targetURL.String())
	ctx := r.Context()
	cachedData, err := p.Cache.Get(ctx, targetURL)
	if err != nil {
		logging.Debugf("Cache retrieval failed - target=%s error=%v", targetURL.String(), err)
		return err
	}

	logging.Debugf("Serving cached response - target=%s bytes=%d", targetURL.String(), len(cachedData))
	w.Header().Set("X-Proxy-Cache", "HIT")
	w.Write(cachedData)
	return nil
}

func (p *ProxyServer) modifyResponse(url url.URL, resp *http.Response) error {
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}

	rewrittenBody, err := p.Pipeline.Process(resp.Body, ct)
	if err != nil {
		return err
	}

	ctx := resp.Request.Context()

	pr, pw := io.Pipe()
	tr := io.TeeReader(rewrittenBody, pw)
	go func() {
		defer rewrittenBody.Close()
		defer pw.Close()

		data, err := io.ReadAll(pr)
		if err == nil {
			logging.Debugf("Caching rewritten response - url=%s bytes=%d ttl=%v", url.String(), len(data), constant.DefaultCacheTTL)
			p.Cache.Set(ctx, url, data, constant.DefaultCacheTTL)
		} else {
			logging.Errorf("Error reading rewritten body for caching - url=%s error=%v", url.String(), err)
		}
	}()

	resp.Body = io.NopCloser(tr)
	resp.ContentLength = -1
	resp.Header.Del("Content-Length")

	return nil
}
