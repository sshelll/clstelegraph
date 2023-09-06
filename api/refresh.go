package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type RefreshTelegraphListResponse struct {
	L map[string]*RefreshTelegraph `json:"l"`
}

type RefreshTelegraph struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	CTime    int64  `json:"ctime"`
	Subjects []struct {
		SubjectName string `json:"subjectName"`
	} `json:"subjects"`
}

func RefreshTelegraphList() (*RefreshTelegraphListResponse, error) {
	url := "https://www.cls.cn/nodeapi/refreshTelegraphList"

	req, _ := http.NewRequest("GET", url, nil)
	addQueryForRefreshTelegramList(req)
	addHeaderForRefreshTelegramList(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	resp := &RefreshTelegraphListResponse{}
	err = json.Unmarshal(body, resp)
	return resp, err
}

func addQueryForRefreshTelegramList(req *http.Request) {
	query := req.URL.Query()
	query.Add("app", "CailianpressWeb")
	query.Add("lastTime", fmt.Sprintf("%d", time.Now().Truncate(time.Minute).Unix()))
	query.Add("os", "web")
	query.Add("sv", "7.7.5")
	req.URL.RawQuery = query.Encode()
}

func addHeaderForRefreshTelegramList(req *http.Request) {
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
