package network

type ApiKey string

func NewApiKey(src string) *ApiKey {
	k := ApiKey(src)
	return &k
}

func (x *ApiKey) String() string {
	return string(*x)
}

func (x *ApiKey) Equals(y *ApiKey) bool {
	lenx := len(*x)
	leny := len(*y)
	if lenx == 0 || leny == 0 {
		return false
	}
	if lenx != leny {
		return false
	}
	for i := 0; i < lenx; i++ {
		if (*x)[i] != (*y)[i] {
			return false
		}
	}
	return true
}
