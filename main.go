package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type Page struct {
	url  string
	html string
}

func fetch(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close() // need to make sure there are no resource leaks
	log.Println("STATUS CODE:", res.StatusCode)

	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	page := &Page{
		url:  url,
		html: string(content),
	}

	re := regexp.MustCompile(`href=["']([^"']+)["']`)
	body := page.html
	links := re.FindAllStringSubmatch(body, -1)
	for _, link := range links {
		fmt.Println(link[1])
	}
}

func main() {
	fetch("https://www.google.com")
}
