package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"crawler/db"
	"crawler/model"
	"crawler/website"

	"golang.org/x/net/html"
)

var (
	link   string
	action string
)

func init() {
	flag.StringVar(&link, "url", "https://aprendagolang.com.br", "url para iniciar visitas")
	flag.StringVar(&action, "action", "website", "qual serviço iniciar")
}

func main() {
	flag.Parse()

	switch action {

	case "website":
		website.Run()
	case "webcrawler":
		done := make(chan bool)
		go visitUrl(link)
		<-done
	default:
		fmt.Printf("action '%s' não reconhecida\n", action)
	}
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
		fmt.Sprintln("status diferente de 200: %d", res.StatusCode)
		return
	}

	htmlTree, err := html.Parse(res.Body)
	if err != nil {
		panic(err)
	}

	links := buildResult(htmlTree)

	fmt.Println(len(links))
}

func checkIfVisited(url string) bool {
	fmt.Printf("Visitando link: %s\n", url)
	if db.IsLinkVisited(url) {
		return true
	}

	fmt.Printf("Link já visitado: %s\n", url)
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
			if err != nil || !strings.HasPrefix(link.Scheme, "http") {
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
