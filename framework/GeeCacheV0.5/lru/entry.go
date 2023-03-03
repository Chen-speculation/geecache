package lru

type entry struct {
	key   string
	value Value
}

func newEntry(key string, value Value) *entry {
	return &entry{
		key:   key,
		value: value,
	}
}
