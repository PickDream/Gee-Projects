package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//定义一个hash函数的实现
//将数据Key映射到 2^32 的环形空间中
type Hash func(data []byte) uint32

type Map struct {
	hash Hash
	//虚拟节点倍数
	replicas int
	//hash环
	keys []int
	//虚拟节点与真实节点的映射表
	//维护从hashcode->真实节点名称的映射
	hashMap map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// @param keys 是一个或者多个真实节点的名称

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		//通过放大倍数方法replicas倍
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			//维护从
			m.hashMap[hash] = key
		}
	}
	//将keys进行排序
	sort.Ints(m.keys)
}

//返回真实节点的名称
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	//找到排序数据的大于keys的位置
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	//使用取余操作确保结果处于一个固定的区间内，最后通过具体的节点值通过节点key的映射找到最终对应的服务名称
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
