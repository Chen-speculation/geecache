package consistent_hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 哈希算法
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           //哈希算法函数
	replicas int            //虚拟节点 = replicas*真实节点
	keys     []int          //哈希环上的虚拟节点的hash值
	hashMap  map[int]string //虚拟节点与真实节点的映射表
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE //使用IEEE多项式返回数据的CRC-32校验和。
	}
	return m
}

// Add adds some keys to the hash.
func (m *Map) Add(keys ...string) { //允许传入 0 或 多个真实节点的名称。
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ { //依次计算虚拟节点的hash，并保存到hashMap中
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) //TODO 这个可能会重复吧！
			m.keys = append(m.keys, hash)                      //加入hash环中
			m.hashMap[hash] = key                              //加入集合中
		}
	}
	sort.Ints(m.keys) //对环的虚拟节点的hash值进行排序
}

func (m *Map) Get(key string) string { //我们需要寻找一个主机，获取key对应的缓存
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]] //取余的作用是达到数组尽头之后回到原点
}
