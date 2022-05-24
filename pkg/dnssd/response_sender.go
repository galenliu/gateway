package dnssd

type ResponseSender struct {
	mServer MdnsServerBase
}

func NewResponseSender() *ResponseSender {
	return &ResponseSender{}
}

func (r ResponseSender) SetServer(ser MdnsServerBase) {
	r.mServer = ser
}
