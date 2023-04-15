package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/sid-008/gocrawl/queue"
)

type Page struct {
	url  string
	html string
}

func fetch(url string, queue *queue.Queue, wg *sync.WaitGroup, urls chan string) {
	defer wg.Done()
	//log.Println(count)
	x := strings.HasPrefix(url, "http")
	if !x {
		return
	}

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close() // need to make sure there are no resource leaks
	//log.Println("STATUS CODE for", url, ":", res.StatusCode) //this was for logging purposes

	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	page := &Page{
		url:  url,
		html: string(content),
	}

	findLinks(page, queue, wg, urls)
}

func findLinks(page *Page, queue *queue.Queue, wg *sync.WaitGroup, urls chan string) {
	defer wg.Done()
	re := regexp.MustCompile(`href=["']([^"']+)["']`)
	body := page.html
	links := re.FindAllStringSubmatch(body, -1)
	for _, link := range links {
		queue.Enqueue(link[1])
		urls <- link[1]
	}
}

func consume(urls chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for link := range urls {
		fmt.Println(link)
	}
}

var urls chan string

// var count uint64

//TODO currently the program runs indefinitely and throws an OOM error, certain memory optimisations have to be made. For instance it could be written to a disk or a db rather than any sort of in memory storage
//TODO the OOM error issue rn is probably the most severe one that interferes with the crawler functioning

func main() {
	wg := &sync.WaitGroup{}
	urls = make(chan string, 100)
	initialUrl := "https://linuxhint.com"
	queue := &queue.Queue{}
	queue.Enqueue(initialUrl)

	for {
		/*count++
		if count > 100000 {
			break
		}*/
		next := queue.Dequeue()
		wg.Add(1)
		go fetch(next, queue, wg, urls)
		wg.Add(1)
		go consume(urls, wg)
	}
	// log.Println("count is:", count)
}
