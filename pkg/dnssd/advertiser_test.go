package dnssd

import "testing"

func TestCommissionableInstanceName(t *testing.T) {
	ad := Advertiser{}
	ad.UpdateCommissionableInstanceName()
	t.Log(ad.mCommissionableInstanceName, "\t\n")
}
