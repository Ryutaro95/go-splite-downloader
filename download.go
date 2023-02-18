package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
)

type Downloader struct {
	url      *url.URL
	splitNum int
	ranges   []string
}

func (d *Downloader) Execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := http.DefaultClient.Get(d.url.String())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err := allowedHTTPAcceptsRanges(response); err != nil {
		return err
	}

	if err := d.splitToRanges(response); err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", "partials")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	if err := d.downloadbyRanges(ctx, tempDir); err != nil {
		return err
	}

	return nil
}

func (d *Downloader) splitToRanges(resp *http.Response) error {
	if err := validSplitRanges(resp, d.splitNum); err != nil {
		return err
	}

	var rangeStart, rangeEnd int

	rangeLength := resp.ContentLength / int64(d.splitNum)

	// 分割回数分のヘッダー情報を ranges に含める
	for i := 0; i < d.splitNum; i++ {
		if i != 0 {
			rangeStart = rangeEnd + 1
		}
		rangeEnd = rangeStart + int(rangeLength)
		if i == d.splitNum-1 && rangeEnd != int(rangeLength) {
			rangeEnd = int(resp.ContentLength)
		}

		d.ranges = append(d.ranges, fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd))
	}

	return nil
}

func (d *Downloader) downloadbyRanges(ctx context.Context, tempDir string) error {
	eg, ctx := errgroup.WithContext(ctx)

	for i, r := range d.ranges {
		i, r := i, r
		eg.Go(func() error {
			req, err := http.NewRequest("GET", d.url.String(), nil)
			if err != nil {
				return err
			}
			req = req.WithContext(ctx)
			req.Header.Set("Range", r)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			partialPath := generatePartialPath(tempDir, i)
			fmt.Println(partialPath)
			fmt.Printf("Downloading range %v / %v (%v) ... \n", i+1, len(d.ranges), r)

			f, err := os.Create(partialPath)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err = io.Copy(f, resp.Body); err != nil {
				return err
			}
			return nil
		})
	}
	return eg.Wait()
}

func allowedHTTPAcceptsRanges(resp *http.Response) error {
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		return errors.New("split download is not supported in this resopnse")
	}
	return nil
}

func validSplitRanges(resp *http.Response, splitNum int) error {
	if int(resp.ContentLength) < splitNum {
		return errors.New("the number of split ranges is larger then file length")
	}
	return nil
}

func generatePartialPath(tempDir string, i int) string {
	base := strings.Join([]string{"partial", strconv.Itoa(i)}, "_")
	return strings.Join([]string{tempDir, base}, "/")
}
