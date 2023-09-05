package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	url := "https://www.cls.cn/nodeapi/refreshTelegraphList"

	req, _ := http.NewRequest("GET", url, nil)
	addQuery(req, 10)
	addHeader(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

func addQuery(req *http.Request, cnt int) {
	query := req.URL.Query()
	query.Add("app", "CailianpressWeb")
	query.Add("lastTime", fmt.Sprintf("%d", time.Now().Truncate(time.Minute*5).Unix()))
	query.Add("os", "web")
	query.Add("sv", "7.7.5")
	query.Add("rn", fmt.Sprintf("%d", cnt))
	req.URL.RawQuery = query.Encode()
}

func addHeader(req *http.Request) {
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("DNT", "1")
	req.Header.Add("Referer", "https://www.cls.cn/telegraph")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	req.Header.Add("sec-ch-ua", `"Not)A;Brand";v="24", "Chromium";v="116"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `macOS`)
}
