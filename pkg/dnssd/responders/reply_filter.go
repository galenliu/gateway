package responders

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QClass"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
)

type ReplyFilter interface {
	Accept(QType.QType, QClass.QClass, *core.FullQName) bool
}
