package access

type AccessControler interface {
	Init(delegate Delegate) error
}

type AccessControl struct {
	mDelegate Delegate
}

func (c *AccessControl) Init(delegate Delegate) {
	c.mDelegate = delegate
}
