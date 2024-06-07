package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func Menu(res http.ResponseWriter, req *http.Request) {
	var (
		keys       map[string]string
		barUrls    []string
		tsUrls     []string
		mapUrls    []string
		statesUrls []string
	)

	_, keys = openFile("UAS_external_event_queue.jsonl")

	barUrls = append(barUrls, "http://localhost:8080/oem-bar")
	barUrls = append(barUrls, "http://localhost:8080/test")
	tsUrls = append(tsUrls, "http://localhost:8080/oem-serie")
	mapUrls = append(mapUrls, "http://localhost:8080/map")

	for _, key := range keys {
		var (
			url_t  string
			params url.Values
		)
		params = url.Values{
			"key": {key},
		}
		url_t = "http://localhost:8080/oem-bar?" + params.Encode()
		barUrls = append(barUrls, url_t)
		url_t = "http://localhost:8080/oem-serie?" + params.Encode()
		tsUrls = append(tsUrls, url_t)
		url_t = "http://localhost:8080/map?" + params.Encode()
		mapUrls = append(mapUrls, url_t)
		url_t = "http://localhost:8080/state?" + params.Encode()
		statesUrls = append(statesUrls, url_t)
	}

	// Set the Content-Type header to text/html
	res.Header().Set("Content-Type", "text/html")

	// Write the HTML opening tags
	fmt.Fprintln(res, "<html><body>")
	fmt.Fprintln(res, "<h1>List of Bar plot </h1>")

	// Iterate over the URLs and create clickable links
	for _, url := range barUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	fmt.Fprintln(res, "<h1>List of Time Series plot </h1>")
	// Iterate over the URLs and create clickable links
	for _, url := range tsUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	fmt.Fprintln(res, "<h1>List of Map plot </h1>")
	// Iterate over the URLs and create clickable links
	for _, url := range mapUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	fmt.Fprintln(res, "<h1>List of States plot </h1>")
	// Iterate over the URLs and create clickable links
	for _, url := range statesUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	// Write the HTML closing tags
	fmt.Fprintln(res, "</body></html>")

}
