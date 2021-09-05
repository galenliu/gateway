package services

import "plugin"

type Service struct {
	p *plugin.Plugin
}

func (s *Service) Start(path string) error {
	var err error
	s.p, err = plugin.Open(path)
	if err != nil {
		return err
	}
	return nil
}
