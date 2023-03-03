package consistentHash

// Hash 哈希算法
type Hash func(data []byte) uint32
