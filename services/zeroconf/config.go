package zeroconf

func DefaultConfig(name string, port int) *DnsSdConf {
	c := &DnsSdConf{
		Name:    name,
		Type:    "_hap._tcp",
		Domain:  "local",
		Host:    "",
		Port:    port,
		Pin:     "12344321", /// default pin
		SetupId: "BBBB",

		id:           "AA AA AA AA AA AA AA",
		discoverable: true,
		mfiCompliant: false,
	}
	c.TextRecords = c.txtRecords()
	return c
}
