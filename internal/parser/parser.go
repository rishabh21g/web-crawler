package parser

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"github.com/rishabh21g/web-crawler/internal/models"
)

func Traverse(node *html.Node, baseURL *url.URL, result *[]models.URLTask, depth int) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, att := range node.Attr {
			if att.Key == "href" {
				href := att.Val
				if href == "" {
					continue
				}
				//filter junks
				if strings.HasPrefix(href, "#") ||
					strings.HasPrefix(href, "javascript:") ||
					strings.HasPrefix(href, "mailto:") ||
					strings.HasPrefix(href, "tel:") {
					continue
				}
				ref, err := url.Parse(href)
				if err != nil {
					continue
				}
				abs := baseURL.ResolveReference(ref)
				*result = append(*result, models.URLTask{
					URL:   abs.String(),
					Depth: depth + 1,
				})
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		Traverse(c, baseURL, result, depth)
	}
}
func ParseHTML(htmlStr string, base string, depth int) ([]models.URLTask, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))

	if err != nil {
		return nil, err
	}
	baseURL, err := url.Parse(base)

	if err != nil {
		return nil, err
	}
	var result []models.URLTask

	Traverse(doc, baseURL, &result, depth)
	return result, nil

}
