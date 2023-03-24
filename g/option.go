package g

type Option[T any] func(*T)

func ApplyOptions[T any, F ~func(*T)](options ...F) T {
	cnt := new(T)
	for _, option := range options {
		option(cnt)
	}
	return *cnt
}

func ApplyOptionsOnDefault[T any, F ~func(*T)](cnt T, options ...F) T {
	for _, option := range options {
		option(&cnt)
	}
	return cnt
}
