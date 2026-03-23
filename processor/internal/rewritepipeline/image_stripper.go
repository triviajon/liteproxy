package rewritepipeline

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type ImageStripper struct{}

func (s *ImageStripper) CanHandle(ct string) bool {
	return strings.Contains(strings.ToLower(ct), "text/html")
}

func (s *ImageStripper) Rewrite(input io.Reader, ct string) (io.ReadCloser, error) {
	doc, err := html.Parse(input)
	if err != nil {
		return nil, err
	}

	s.stripImages(doc)
	pr, pw := io.Pipe()

	go func() {
		err := html.Render(pw, doc)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		pw.Close()
	}()

	return pr, nil
}

// stripImages removes <img> elements.
func (s *ImageStripper) stripImages(n *html.Node) {
	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		if c.Type == html.ElementNode && c.Data == "img" {
			n.RemoveChild(c)
		} else {
			s.stripImages(c)
		}
		c = next
	}
}
