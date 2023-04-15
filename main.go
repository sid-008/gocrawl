package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/sid-008/gocrawl/queue"
)

type Page struct {
	url  string
	html string
}

func fetch(url string, queue *queue.Queue) {
	x := strings.HasPrefix(url, "https")
	if !x {
		return
	}

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close() // need to make sure there are no resource leaks
	log.Println("STATUS CODE:", res.StatusCode)

	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	page := &Page{
		url:  url,
		html: string(content),
	}

	findLinks(page, queue)
}

func findLinks(page *Page, queue *queue.Queue) {
	re := regexp.MustCompile(`href=["']([^"']+)["']`)
	body := page.html
	links := re.FindAllStringSubmatch(body, -1)
	for _, link := range links {
		fmt.Println(link[1])
		queue.Enqueue(link[1])
	}
}

func main() {
	initialUrl := "https://google.com"
	queue := &queue.Queue{}
	queue.Enqueue(initialUrl)

	for {
		next := queue.Dequeue()
		fmt.Println("The next url to be explored is:\n\n\n", next)
		fetch(next, queue)
	}
}
