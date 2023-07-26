package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const defaultAddr = "https://www.google.com"

func main() {
	var addr string
	flag.StringVar(&addr, "addr", defaultAddr, "Check HTTP GET Response is OK")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		tk := time.NewTicker(3 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				if err := IsRespOK(addr); err != nil {
					fmt.Printf("%s address http get request failed: %v\n", addr, err)
				}
			}
		}
	}(ctx)

	time.Sleep(1 * time.Minute)
	cancel()
}

func IsRespOK(addr string) error {
	resp, err := http.Get(addr)
	if err != nil {
		return errors.Wrap(err, "response is not ok")
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("response is not ok")
	}
	return nil
}
