package main

import (
	"fmt"
	"net/url"
)

type Downloader struct {
	url      *url.URL
	splitNum int
	ranges   []string
}

func (d *Downloader) Execute() error {
	fmt.Println(d.url)
	fmt.Println(d.splitNum)
	fmt.Println(d.ranges)
	return nil
}
