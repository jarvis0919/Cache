package cachegroup

//实现value接口，作为cache缓存中的最基础的数据储存结构
type ByteView struct {
	b []byte
}

//获取其中数据
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

//获取字符串数据
func (v ByteView) String() string {
	return string(v.b)
}

//获取长度
func (v ByteView) Len() int {
	return len(v.b)
}

//复制数据
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
