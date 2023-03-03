package gee_cache

// PeerClient 一个节点必须实现这个接口
type PeerClient interface {
	// Get 用于从对应 gee_cache 查找缓存值。PeerClient 就对应于上述流程中的 HTTP 客户端。
	Get(group string, key string) ([]byte, error)
}

// PeerPicker 根据传入的 缓存数据对应的key 选择相应节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerClient, ok bool)
}
