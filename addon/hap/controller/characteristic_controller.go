package controller

import (
	"bytes"
	"encoding/json"
	"gateway/accessory"
	"gateway/accessory/characteristic"
	"gateway/logger"
	"gateway/services/hap"
	"gateway/services/hap/data"
	"github.com/xiam/to"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"strings"
)

type CharacteristicController struct {
	container *accessory.Container
}

func NewCharacteristicController(container *accessory.Container) *CharacteristicController {
	return &CharacteristicController{container: container}
}

// 参数 ：/characteristics?id=1.4,1.5
func (ctr *CharacteristicController) HandleGetCharacteristics(values url.Values, conn net.Conn) (io.Reader, error) {
	var b bytes.Buffer
	var chs []data.Characteristic

	// 例：id=1.4,1.5
	paths := strings.Split(values.Get("id"), ",")
	//paths[1.4,1.5]
	for _, p := range paths {
		if ids := strings.Split(p, "."); len(ids) == 2 {
			aid := to.Uint64(ids[0])
			iid := to.Uint64(ids[1])
			c := data.Characteristic{AccessoryID: aid, CharacteristicID: iid}
			if ch := ctr.GetCharacteristic(aid, iid); ch != nil {
				c.Value = ch.GetValueFromConnection(conn)
			} else {
				c.Status = hap.StatusServiceCommunicationFailure
			}
			chs = append(chs, c)
		}
	}
	result, err := json.Marshal(&data.Characteristics{Characteristics: chs})
	if err != nil {
		logger.Warning.Print(err)
	}
	b.Write(result)
	return &b, err

}

//更新 characteristc requst body like: [{aid:11,iid:11,value:10}{aid:11,iid:12,value:10,events:true}]
func (ctr *CharacteristicController) HandleUpdateCharacteristics(reader io.Reader, conn net.Conn) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	var chs data.Characteristics
	err = json.Unmarshal(b, &chs)
	for _, ch := range chs.Characteristics {
		c := ctr.GetCharacteristic(ch.AccessoryID, ch.CharacteristicID)
		if c == nil {
			logger.Info.Printf("Could not find characteristic with aid &d and iid &d", ch.AccessoryID, ch.CharacteristicID)
			continue
		}
		if ch.Value != nil {
			c.UpdateValueFromConnection(ch.Value, conn)
		}
		if events, ok := ch.Events.(bool); ok == true {
			c.Events = events
		}

	}
	return nil
}

func (c *CharacteristicController) GetCharacteristic(aid uint64, iid uint64) *characteristic.Characteristic {
	for _, a := range c.container.Accessories {
		if a.ID == aid {
			for _, s := range a.GetServices() {
				for _, c := range s.GetCharacteristics() {
					if c.ID == iid {
						return c
					}
				}
			}

		}
	}
	return nil
}
