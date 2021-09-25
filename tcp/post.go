package tcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HttpClient struct {
	c        http.Client
	ApiToken string
}

func (c *HttpClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err

	}

	return c.Do(req)

}

func (c *HttpClient) Post(url, contentType string, containers interface{}) (resp *http.Response, err error) {

	requestBody, err := json.Marshal(containers)

	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewBuffer(requestBody)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err

	}

	req.Header.Set("Content-Type", contentType)

	fmt.Println(req.URL.Host)

	return c.Do(req)

}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	// req.Header.Set("x-api-key", c.ApiToken)

	// req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0")
	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	// req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("Cookie", "_ga_25N7NHMSZH=GS1.1.1620105353.5.0.1620105353.0; _ga=GA1.1.861892589.1618714056")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
	// req.Header.Set("If-Modified-Since", "Wed, 14 Apr 2021 07:02:41 GMT")

	// 	GET / HTTP/1.1
	// Host: tracking.bcscdepot.com
	// User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0
	// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
	// Accept-Language: en-US,en;q=0.5
	// Accept-Encoding: gzip, deflate
	// Connection: keep-alive
	// Cookie: _ga_25N7NHMSZH=GS1.1.1620105353.5.0.1620105353.0; _ga=GA1.1.861892589.1618714056
	// Upgrade-Insecure-Requests: 1
	// If-Modified-Since: Wed, 14 Apr 2021 07:02:41 GMT
	// Cache-Control: max-age=0

	// req.Host = "tracking.bcscdepot.com"
	req.Header.Add("Host", "tracking.bcscdepot.com")
	req.Header.Add("User-Agent", "myClient")

	fmt.Println("Host:", req.Header.Get("Host"))

	return c.c.Do(req)

}
