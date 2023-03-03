package group

import (
	"GeeCacheV0.7/protobuf"
	"fmt"
	"google.golang.org/protobuf/proto"
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
func (h *httpClient) Get(in *protobuf.Request, out *protobuf.Response) error {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()), //QueryEscape函数对s进行转码使之可以安全的用在URL查询里。
		url.QueryEscape(in.GetKey()),
	)

	res, err := http.Get(u)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	//protobuf解码
	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}
