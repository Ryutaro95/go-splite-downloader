package main

import "net/url"

type Downloader struct {
	url      *url.URL
	splitNum int
	ranges   []string
}
