package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	splitNum := flag.Int("n", 10, "number of splits")
	flag.Parse()
	inputUrl := flag.Arg(0)
	if err := validInputURL(inputUrl); err != nil {
		die(err)
	}

	url, err := url.Parse(inputUrl)
	if err != nil {
		die(err)
	}
	download := &Downloader{
		url:      url,
		splitNum: *splitNum,
		ranges:   []string{"bytes=0-100", "bytes=101-200"},
	}
	if err := download.Execute(); err != nil {
		die(err)
	}
}

func validInputURL(url string) error {
	if url == "" {
		return errors.New("URL is required")
	}
	return nil
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "error: ", err)
	os.Exit(1)
}
