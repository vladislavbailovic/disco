package store

type Storer interface {
	Fetch(*Key) (string, error)
	Put(*Key, string) error
}

func Default() Storer {
	return NewPlainStore()
}
