package getter

type Getter interface {
	Get(key string) ([]byte, error)
}

// 接口型函数 实现Getter接口
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
