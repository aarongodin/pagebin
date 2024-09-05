package core

type ContextKey string

func (k ContextKey) String() string {
	return string(k)
}
