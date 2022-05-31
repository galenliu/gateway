package interaction_affordance

import (
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type InteractionAffordance struct {
	AtType       string                   `json:"@type,omitempty,optional"`
	Title        string                   `json:"title,omitempty,optional"`
	Titles       map[string]string        `json:"titles,omitempty,optional"`
	Description  string                   `json:"description,omitempty,optional"`
	Descriptions map[string]string        `json:"descriptions,omitempty,optional"`
	Forms        []controls.Form          `json:"forms,omitempty,mandatory" wot:"optional"`
	UriVariables map[string]schema.Schema `json:"uriVariables,omitempty"`
}

func (v *InteractionAffordance) UnmarshalJSON(data []byte) error {

	var node ast.Node
	var err error
	var nodeMap map[string]ast.Node

	node, err = sonic.Get(data, "@type")
	if node.Exists() && err == nil {
		v.AtType, _ = node.String()
	}

	node, err = sonic.Get(data, "title")
	if node.Exists() && err == nil {
		v.Title, _ = node.String()
	}

	node, err = sonic.Get(data, "description")
	if node.Exists() && err == nil {
		v.Description, _ = node.String()
	}

	node, err = sonic.Get(data, "titles")
	if node.Exists() && err == nil {
		d, _ := node.MarshalJSON()
		_ = sonic.Unmarshal(d, &v.Titles)
	}

	node, err = sonic.Get(data, "descriptions")
	if node.Exists() && err == nil {
		d, _ := node.MarshalJSON()
		_ = sonic.Unmarshal(d, &v.Descriptions)
	}

	node, err = sonic.Get(data, "forms")
	if node.Exists() && err == nil {
		d, _ := node.MarshalJSON()
		_ = sonic.Unmarshal(d, &v.Forms)
	}

	node, err = sonic.Get(data, "uriVariables")
	if node.Exists() && err == nil {
		nodeMap, err = node.MapUseNode()
		if err == nil {
			for name, value := range nodeMap {
				d, _ := value.MarshalJSON()
				s, err := schema.MarshalSchema(d)
				if err != nil {
					v.UriVariables[name] = s
				}
			}
		}
	}

	return nil
}
