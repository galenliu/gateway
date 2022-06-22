package platform

import (
	"fmt"
	"net"
	"testing"
)

func GetPrimaryWiFiMACAddress() (i []net.Interface) {
	ifs, _ := net.Interfaces()
	return ifs
}

func TestThingMarshal(t *testing.T) {
	is := GetPrimaryWiFiMACAddress()
	for _, i := range is {
		fmt.Printf("------------------ -----------\t\n")
		fmt.Printf("Instance: %v \t\n", i.Name)
		fmt.Printf("Index: %v \t\n", i.Index)
		fmt.Printf("Flags: %v \t\n", i.Flags)
		fmt.Printf("MTU: %v \t\n", i.MTU)
		fmt.Printf("HardwareAddr: %v \t\n", i.HardwareAddr.String())
		mAdds, _ := i.MulticastAddrs()
		as, _ := i.Addrs()
		for _, a := range as {
			fmt.Printf("Addrs: %v \t\n", a.String())
		}
		for _, a := range mAdds {
			fmt.Printf("MulticastAddr: %v \t\n", a.String())
		}
		fmt.Printf("------------------ -----------\t\n")
	}
}
