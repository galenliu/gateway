package gateway

import "fmt"

var (
	majorVersion = 1
	minorVersion = 0
	patchVersion = 0

	commit string // automatically set git commit hash

	shortVersion = fmt.Sprintf("%v.%v", majorVersion, minorVersion)
	fullVersion  = fmt.Sprintf("%v.%v", shortVersion, patchVersion)

	Version = func() string {
		if commit != "" {
			return fullVersion + "-" + commit
		}
		return fullVersion + "-dev"
	}()
)
