package dnssd

type CommissioningMode int
type CommssionAdvertiseMode = int

const (
	KDisabled           CommissioningMode      = -0 // Commissioning Mode is disabled, CM=0 in DNS-SD key/value pairs
	EnableBasic         CommissioningMode      = 1  // Basic Commissioning Mode, CM=1 in DNS-SD key/value pairs
	EnabledEnhanced     CommissioningMode      = 2  // Enhanced Commissioning Mode, CM=2 in DNS-SD key/value pairs
	KCommissionableNode CommssionAdvertiseMode = 0
	KCommissioner       CommssionAdvertiseMode = 1
)
