package zCache

type ByteView struct {
	b []byte //存储真实的缓存值
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

// ByteSlice 返回一个副本
// b 是只读的，使用 ByteSlice() 方法返回一个拷贝，防止缓存值被外部程序修改。
// 获取时，即调用 get() 时，不需要复制，ByteView 是只读的，不可修改。通过 ByteSlice() 或 String() 方法取到缓存值的副本。只读属性，是设计 ByteView 的主要目的之一。
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
