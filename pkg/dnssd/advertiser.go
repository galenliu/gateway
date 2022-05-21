package dnssd

const (
	CommssionAdvertiseModeCommissionableNode = iota
	CommssionAdvertiseModeCommissioner
)

const (
	CommissioningModeDisabled        = iota // Commissioning Mode is disabled, CM=0 in DNS-SD key/value pairs
	CommissioningModeEnabledBasic           // Basic Commissioning Mode, CM=1 in DNS-SD key/value pairs
	CommissioningModeEnabledEnhanced        // Enhanced Commissioning Mode, CM=2 in DNS-SD key/value pairs
)
