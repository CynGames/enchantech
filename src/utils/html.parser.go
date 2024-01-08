package utils

import "golang.org/x/net/html"

func GetThumbnailProperties(n *html.Node, properties map[string]string) {
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style" || n.Data == "body") {
		return
	}

	if n.Type == html.ElementNode && n.Data == "meta" {
		//fmt.Println("meta: ", n.Attr)
		var property, content string

		for _, attr := range n.Attr {

			if attr.Key == "property" && (attr.Val == "og:title" || attr.Val == "og:description" || attr.Val == "og:image") {
				property = attr.Val
			}

			if attr.Key == "content" {
				content = attr.Val
			}

		}

		if property != "" && content != "" {
			properties[property] = content
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		GetThumbnailProperties(c, properties)
	}
}
