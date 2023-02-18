package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	url, err := url.Parse("https://example.com")
	if err != nil {
		die(err)
	}
	download := &Downloader{
		url:      url,
		splitNum: 8,
		ranges:   []string{"bytes=0-100", "bytes=101-200"},
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "error: ", err)
	os.Exit(1)
}
