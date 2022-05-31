package responders

import "github.com/galenliu/gateway/pkg/dnssd/core"

type Responder struct {
	mQType core.QType
	mQName core.FullQName
}
