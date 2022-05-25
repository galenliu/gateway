package dnssd

type CommissioningMode int
type CommssionAdvertiseMode = int

const (
	Disable            CommissioningMode      = -0 // Commissioning Mode is disabled, CM=0 in DNS-SD key/value pairs
	EnableBasic        CommissioningMode      = 1  // Basic Commissioning Mode, CM=1 in DNS-SD key/value pairs
	EnabledEnhanced    CommissioningMode      = 2  // Enhanced Commissioning Mode, CM=2 in DNS-SD key/value pairs
	CommissionableNode CommssionAdvertiseMode = 0
	Commissioner       CommssionAdvertiseMode = 1
)
