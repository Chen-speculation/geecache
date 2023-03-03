package gee_cache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type httpClient struct {
	baseURL string //TODO baseURL最后一个字符一定要是"/"
}

func NewHttpClient(baseURL string) PeerClient {
	return &httpClient{
		baseURL: baseURL,
	}
}

// Get 客户端方法
func (h *httpClient) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group), //QueryEscape函数对s进行转码使之可以安全的用在URL查询里。
		url.QueryEscape(key),
	)
	res, err := http.Get(u)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}
