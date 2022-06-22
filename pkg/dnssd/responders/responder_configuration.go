package responders

import "github.com/galenliu/gateway/pkg/dnssd/record"

type ResponseConfiguration struct {
	mTtlSecondsOverride *uint32
}

func (c ResponseConfiguration) Adjust(r record.ResourceRecord) {
	if c.mTtlSecondsOverride != nil {
		r.SetTtl(*c.mTtlSecondsOverride)
	}
}

func (c *ResponseConfiguration) SetTtlSecondsOverride(i uint32) {
	c.mTtlSecondsOverride = &i
}
