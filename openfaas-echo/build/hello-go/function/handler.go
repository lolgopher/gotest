package function

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func Handle(input []byte) string {
	url := string(input)
	log.Print("Hello! You said:\n" + url + "\n\n")

	// https://datahub.kakaocloud.ai/api/search/v2/metaform-contents/prod/kepGGyXzvZ4SRWIcUakWQbs0A?qnorm=%EC%84%9C%EC%9A%B8%EC%A7%80%EC%97%AD%EB%B0%B0%EC%86%A1&collection=CMS-SVC-FAQ-COM_1&entity=
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	reqDump, _ := httputil.DumpRequest(req, true)
	respDump, _ := httputil.DumpResponse(resp, true)
	return fmt.Sprintf("Request: %s\n\nResponse: %s", reqDump, respDump)
}
