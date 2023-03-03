package group

import (
	"GeeCacheV0.7/cache"
	"GeeCacheV0.7/consistentHash"
	"GeeCacheV0.7/protobuf"
	"GeeCacheV0.7/utils"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/geeCache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self     string //记录自己的地址，包括	”主机名/IP“ 和 ”端口“。
	basePath string //作为节点间通讯地址的前缀，默认是 defaultBasePath
	//下面属性是为了实现"节点选择"功能
	mu         sync.Mutex
	peers      consistentHash.IMap
	httpClient map[string]PeerClient //key by e.g. "http://47.98.251.199:8888"
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...any) {
	//v一定要拆包
	utils.Logger().Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP 服务端方法
func (p *HTTPPool) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, p.basePath) {
		utils.Logger().Errorln("HTTPPoll serving unexpected path: " + req.URL.Path)
		http.Error(resp, "HTTPPoll serving unexpected path: "+req.URL.Path, http.StatusBadRequest)
		return
	}
	p.Log("%s %s", req.Method, req.URL.Path)
	// basePath/GroupName/<key> ==> GroupName/Key
	parts := strings.SplitN(req.URL.Path[len(cache.DecodeBasePath(p.basePath)):], "/", 2)
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

	//TODO 将需要传输的数据变成a proto message.
	bytes, err := proto.Marshal(&protobuf.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/octet-stream")
	resp.Write(bytes)
}

// GetPeer 知道缓存的key,找到主机节点
func (p *HTTPPool) GetPeer(key string) (peer PeerClient, ok bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if hostName := p.peers.Get(key); hostName != "" && hostName != p.self {
		p.Log("Pick PeerHostName: %s", hostName)
		return p.httpClient[hostName], true
	}
	return nil, false
}

// Set peer是 http://ip:port 后面没有"/"
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistentHash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpClient = make(map[string]PeerClient, len(peers))
	for _, peer := range peers {
		p.httpClient[peer] = NewHttpClient(peer + p.basePath)
	}
}
