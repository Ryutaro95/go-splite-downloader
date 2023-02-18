package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	splitNum := flag.Int("n", 10, "number of splits")
	inputUrl := flag.Arg(0)
	url, err := url.Parse(inputUrl)
	if err != nil {
		die(err)
	}
	download := &Downloader{
		url:      url,
		splitNum: *splitNum,
		ranges:   []string{"bytes=0-100", "bytes=101-200"},
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "error: ", err)
	os.Exit(1)
}
