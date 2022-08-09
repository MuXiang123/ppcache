// Package ppcache
// @author    : MuXiang123
// @time      : 2022/7/27 18:29
package ppcache

// ByteView 保存字节不可变的view
type ByteView struct {
	b []byte // 存储真正的缓存值，byte类型能支持任意数据类型的存储
}

// Length 返回view的长度
func (v ByteView) Length() int64 {
	// 被缓存对象必须实现Value里面的的接口
	return int64(len(v.b))
}

// ByteSlice 以字节切片的形式返回数据副本
func (v ByteView) ByteSlice() []byte {
	//防止外部程序修改变量
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

//以字符串的形式返回数据，并且再需要的的时候复制
func (v ByteView) String() string {
	return string(v.b)
}
