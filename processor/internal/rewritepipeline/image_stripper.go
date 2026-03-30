package rewritepipeline

import (
	"fmt"
	"io"
	"strings"

	"github.com/triviajon/liteproxy/processor/internal/logging"
	"golang.org/x/net/html"
)

type ImageStripper struct{}

func (s *ImageStripper) CanHandle(ct string) bool {
	return strings.Contains(strings.ToLower(ct), "text/html")
}

// Rewrite removes all <img> elements from HTML content.
// Requires that input is not nil and ct is not empty.
// Returns an io.ReadCloser with the modified HTML, otherwise an error from parsing or describing a constraint violation.
func (s *ImageStripper) Rewrite(input io.Reader, ct string) (io.ReadCloser, error) {
	if input == nil {
		return nil, fmt.Errorf("input must not be nil")
	}
	if ct == "" {
		return nil, fmt.Errorf("ct must not be empty")
	}

	logging.Debugf("Starting HTML parsing - content_type=%s", ct)
	doc, err := html.Parse(input)
	if err != nil {
		logging.Errorf("HTML parsing failed - error=%v", err)
		return nil, err
	}

	logging.Debugf("HTML parsed successfully, stripping images")
	imageCount := s.stripImages(doc)
	logging.Debugf("Stripped %d image element(s)", imageCount)
	pr, pw := io.Pipe()

	go func() {
		err := html.Render(pw, doc)
		if err != nil {
			logging.Errorf("HTML render error - error=%v", err)
			pw.CloseWithError(err)
			return
		}
		logging.Debugf("HTML rendered successfully")
		pw.Close()
	}()

	return pr, nil
}

// stripImages removes <img> elements and returns the count of removed images.
func (s *ImageStripper) stripImages(n *html.Node) int {
	count := 0
	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		if c.Type == html.ElementNode && c.Data == "img" {
			n.RemoveChild(c)
			count++
		} else {
			count += s.stripImages(c)
		}
		c = next
	}
	return count
}
