package addon

//
//import (
//	"gateway/db"
//	"gateway/plugin"
//)
//
//type Addon struct {
//	Manifest     *PacketManifest
//	PluginServer *plugin.PluginsServer
//}
//
//type PacketManifest struct {
//	ID          string `yaml:"id"`
//	Author      string `yaml:"author"`
//	Name        string `yaml:"name"`
//	ShortName   string `yaml:"short_name"`
//	Description string `yaml:"description"`
//	HomepageUrl string `yaml:"homepage_url"`
//	Version     string `yaml:"version"`
//	License     string `yaml:"license"`
//	Enabled     bool
//
//	ManifestVersion int `yaml:"manifest_version"`
//}
//
//func (p *PacketManifest) toSetting() *db.Setting {
//	return &db.Setting{
//		ID:              p.ID,
//		Author:          p.Author,
//		Name:            p.Name,
//		ShortName:       p.ShortName,
//		Description:     p.Description,
//		Version:         p.Version,
//		License:         p.License,
//		ManifestVersion: p.ManifestVersion,
//		Enabled:         false,
//	}
//}
//*/
