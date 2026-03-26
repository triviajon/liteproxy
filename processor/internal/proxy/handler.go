package proxy

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/triviajon/liteproxy/processor/internal/cache"
	"github.com/triviajon/liteproxy/processor/internal/constant"
	"github.com/triviajon/liteproxy/processor/internal/rewritepipeline"
	"github.com/triviajon/liteproxy/processor/internal/util"
)

type ProxyServer struct {
	Pipeline rewritepipeline.Pipeline
	Cache    cache.Cache
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET supported", http.StatusMethodNotAllowed)
		return
	}

	rawTarget := r.URL.Query().Get("url")
	if rawTarget == "" {
		http.Error(w, "Missing required query parameter: url", http.StatusBadRequest)
		return
	}

	targetURL, err := url.Parse(rawTarget)
	if err != nil || targetURL.Scheme == "" || targetURL.Host == "" {
		http.Error(w, "Invalid url parameter", http.StatusBadRequest)
		return
	}

	if err := p.serveFromCache(w, r, *targetURL); err == nil {
		return
	}

	modifyResponseFn, err := util.Bind1(p.modifyResponse, *targetURL)
	if err != nil {
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

	proxy.ServeHTTP(w, r)
}

func (p *ProxyServer) serveFromCache(w http.ResponseWriter, r *http.Request, targetURL url.URL) error {
	ctx := r.Context()
	cachedData, err := p.Cache.Get(ctx, targetURL)
	if err != nil {
		return err
	}

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
			p.Cache.Set(ctx, url, data, constant.DefaultCacheTTL)
		}
	}()

	resp.Body = io.NopCloser(tr)
	resp.ContentLength = -1

	return nil
}
