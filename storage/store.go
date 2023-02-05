package storage

type Storer interface {
	Fetch(*Key) (Valuer, error)
	Put(*Key, string) error
}

type Valuer interface {
	Value() string
}

func Default() Storer {
	return NewPlainStore()
}
