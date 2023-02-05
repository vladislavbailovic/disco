package storage

type Storer interface {
	Fetch(*Key) (Valuer, error)
	Put(*Key, string) error
	Delete(*Key) error
	Stats() Valuer
}

type Valuer interface {
	Value() string
	MIME() ContentType
}

func Default() Storer {
	return NewPlainStore()
}

type ContentType uint

const (
	ContentTypeText ContentType = iota
	ContentTypeJSON
)

func (x ContentType) String() string {
	switch x {
	case ContentTypeText:
		return ""
	case ContentTypeJSON:
		return ""
	}
	panic("Unknown content type")
}
