package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"crawler/db"
	"crawler/model"

	"golang.org/x/net/html"
)

var visited = map[string]bool{}

func main() {
	visitUrl("https://aprendagolang.com.br")
}

func visitUrl(url string) {

	if checkIfVisited(url) {
		return
	}

	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintln("status diferente de 200: %d", res.StatusCode))
	}

	htmlTree, err := html.Parse(res.Body)
	if err != nil {
		panic(err)
	}

	links := buildResult(htmlTree)

	fmt.Println(len(links))
}

func checkIfVisited(url string) bool {
	if ok := visited[url]; ok {
		return true
	}

	visited[url] = true
	return false
}

func buildResult(htmlTree *html.Node) []string {
	links := []string{}

	extractLinks(htmlTree, &links)

	return links
}

func extractLinks(node *html.Node, links *[]string) {
	if node.Type == html.ElementNode && node.Data == "a" {

		for _, attr := range node.Attr {
			if attr.Key != "href" {
				continue
			}

			link, err := url.Parse(attr.Val)
			if err != nil || link.Scheme == "" {
				continue
			}

			url := link.String()

			visitedLink := model.VisitedLink{
				Website:     link.Host,
				Link:        link.String(),
				VisitedDate: time.Now(),
			}

			*links = append(*links, url)
			db.Insert("links", visitedLink)

			visitUrl(url)
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		extractLinks(child, links)
	}
}
