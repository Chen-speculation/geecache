package gee_cache

import (
	"GeeCacheV0.4/utils"
	"fmt"
	"net/http"
	"strings"
)

type IHTTPServer interface {
	Run() error
}

const defaultBasePath = "/geeCache/"

type HTTPServer struct {
	self     string //记录自己的地址，包括	”主机名/IP“ 和 ”端口“。
	basePath string //作为节点间通讯地址的前缀，默认是 defaultBasePath
}

func NewHTTPServer(self string) IHTTPServer {
	return &HTTPServer{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPServer) log(format string, v ...any) {
	//v一定要拆包
	utils.Logger().Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	if !strings.HasPrefix(req.URL.Path, p.basePath) {
		utils.Logger().Errorln("HTTPPoll serving unexpected path: " + req.URL.Path)
		http.Error(resp, "HTTPPoll serving unexpected path: "+req.URL.Path, http.StatusBadRequest)
		return
	}

	p.log("%s %s", req.Method, req.URL.Path)

	// basePath/GroupName/<key>
	parts := strings.SplitN(req.URL.Path[len(decodeBasePath(p.basePath)):], "/", 2)
	if len(parts) != 2 {
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(resp, "no such gee_cache: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/octet-stream")
	resp.Write(view.ByteSlice())
}

func (p *HTTPServer) Run() error {
	utils.Logger().Infof("you request url should %s%sGroupName/cacheKey", p.self, p.basePath)
	return http.ListenAndServe(p.self, p)
}
