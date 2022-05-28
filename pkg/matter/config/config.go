package config

var (
	ChipPort                                    = 5540
	ChipUdcPort                                 = ChipPort + 10
	ConfigNetworkLayerBle                       bool
	ChipDeviceConfigEnablePairingAutostart      bool
	ChipDeviceConfigEnableCommissionerDiscovery bool
)
