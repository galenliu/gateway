package server

type ServerInitParams struct {
	operationalServicePort        int
	userDirectedCommissioningPort int
}

func DefaultServerInitParams() *ServerInitParams {
	return &ServerInitParams{
		operationalServicePort:        ChipPort,
		userDirectedCommissioningPort: ChipUdcPort,
	}
}
