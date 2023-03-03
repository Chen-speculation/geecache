package group

import "GeeCacheV0.7/protobuf"

// PeerClient 一个节点必须实现这个接口
type PeerClient interface {
	// Get 使用protobuf进行通信
	Get(in *protobuf.Request, out *protobuf.Response) error
}

// PeerGetter 根据传入的 缓存数据对应的key 选择相应节点
type PeerGetter interface {
	GetPeer(key string) (peer PeerClient, ok bool)
}
