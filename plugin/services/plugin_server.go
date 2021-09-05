package services

import "sync"

type PluginsServer struct {
	plugins sync.Map
}

func (s *PluginsServer) getPlugin(id string) *Plugin {
	p, ok := s.plugins.Load(id)
	plugin, ok := p.(*Plugin)
	if !ok {
		return nil
	}
	return plugin
}

func (s *PluginsServer) getPlugins() (plugins []*Plugin) {
	s.plugins.Range(func(key, value interface{}) bool {
		p, ok := value.(*Plugin)
		if ok {
			plugins = append(plugins, p)
		}
		return true
	})
	return
}
