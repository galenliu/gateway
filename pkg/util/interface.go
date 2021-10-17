package util

import "github.com/kardianos/service"

type Program struct {
	start func()
	stop  func()
}

func (p *Program) Start(s service.Service) error {
	// Run should not block. Do the actual work async.
	go p.start()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	p.stop()
	return nil
}
