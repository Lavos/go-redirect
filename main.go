package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	MAX_REQUESTS = 5
)

type Redirect struct {
	URL string
	RedirectCount int
}

func main () {
	list := make([]string, 0)
	cursor := 0
	response := make(chan Redirect)
	redirects := make([]Redirect, 0)

	for {
		var url string
		_, err := fmt.Scanln(&url)

		if err != nil {
			break
		}

		list = append(list, url)
	}

	log.Printf("got list: %v", list)

	for x := 0; x < MAX_REQUESTS; x++ {
		go doRequest(list[cursor], response)
		cursor++
	}

	for n := 0; n < len(list); n++ {
		redirects = append(redirects, <-response)

		if cursor < len(list) {
			go doRequest(list[cursor], response)
			cursor++
		}
	}

	for _, r := range redirects {
		fmt.Printf("%v\t%v\n", r.URL, r.RedirectCount)
	}
}

func doRequest (url string, response chan Redirect) {
	redirects := 0

	c := &http.Client{
		CheckRedirect: func (req *http.Request, via []*http.Request) error {
			redirects++

			return nil
		},
	}

	log.Printf("checking url: %v", url)
	c.Get(url)
	response <- Redirect{
		URL: url,
		RedirectCount: redirects,
	}
}
