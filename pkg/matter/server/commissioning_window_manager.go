package server

import "github.com/galenliu/gateway/pkg/dnssd"

type CommissioningWindowManager struct {
	mServer      *Server
	mAppDelegate any
}

func NewCommissioningWindowManager(s *Server) *CommissioningWindowManager {
	return &CommissioningWindowManager{
		mServer: s,
	}
}

func (m *CommissioningWindowManager) SetAppDelegate(delegate any) {
	m.mAppDelegate = delegate
}

func (m *CommissioningWindowManager) OpenBasicCommissioningWindow() error {
	return nil
}

func (m *CommissioningWindowManager) GetCommissioningMode() dnssd.CommissioningMode {
	return dnssd.KDisabled
}
