package cache

// ByteView 只读(所以b保存的指针不会发生改变，所以下面的方法没毛笔1)
type ByteView struct {
	b []byte
}

func NewByteView(b []byte) ByteView {
	return ByteView{b: CloneBytes(b)}
}

func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 为了保证只读，选择复制缓存
func (v ByteView) ByteSlice() []byte {
	return CloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}
