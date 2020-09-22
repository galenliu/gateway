package gateway

type RuntimeConfig struct {
	cfgDir        string
	logRotateDays int
	verbose       bool
	reset         bool
}

func NewRuntimeConfig(cfgDir string, logRotateDays int, verbose bool, reset bool) *RuntimeConfig {
	return &RuntimeConfig{
		cfgDir:        cfgDir,
		logRotateDays: logRotateDays,
		verbose:       verbose,
		reset:         reset,
	}
}
