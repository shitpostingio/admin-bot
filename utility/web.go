package utility

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

//UnshortenURL unshortens an URL performing a web request
func UnshortenURL(shortURL string) (unshortenedURL string, err error) {

	if !strings.HasPrefix(shortURL, "http") {
		shortURL = fmt.Sprintf("http://%s", shortURL)
	}

	data, err := http.Head(shortURL) // nolint: gosec
	if err != nil {
		return
	}
	defer CloseSafely(data.Body)

	unshortenedURL = data.Request.URL.String()
	return unshortenedURL, err
}

//IsGroupOrChannelHandle performs a web request to see if an handle belongs to a channel or a group
func IsGroupOrChannelHandle(handle string) bool {

	//Target URL creation
	url := fmt.Sprintf("https://t.me/%s", handle)

	//Web request
	data, err := http.Get(url) // nolint: gosec
	if err != nil {
		return false
	}
	defer CloseSafely(data.Body)

	//HTML parsing
	bodyResp, err := html.Parse(data.Body)
	if err != nil {
		return false
	}

	htmlNode := getElementByClass(bodyResp, "tgme_action_button_new")
	if htmlNode == nil {
		return false
	}

	return htmlNode.FirstChild.Data == "View in Telegram"
}

/*
 * FUNCTIONS FOR WEB REQUESTS
 * used in IsUnwantedHandle
 */
func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func checkClass(n *html.Node, class string) bool {
	if n.Type == html.ElementNode {
		s, ok := getAttribute(n, "class")
		if ok && s == class {
			return true
		}
	}
	return false
}

func traverse(n *html.Node, class string) *html.Node {
	if checkClass(n, class) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, class)
		if result != nil {
			return result
		}
	}

	return nil
}

func getElementByClass(n *html.Node, class string) *html.Node {
	return traverse(n, class)
}
