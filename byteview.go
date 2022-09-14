package gocache

// ByteView 只读的数据结构 避免缓存被修改
type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {

	return string(v.b)

}
func cloneBytes(bs []byte) []byte {

	c := make([]byte, len(bs))
	copy(c, bs)
	return c
}
