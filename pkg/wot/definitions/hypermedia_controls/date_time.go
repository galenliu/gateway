package hypermedia_controls

import (
	"time"
)

type DataTime struct {
	time.Time
}

//func (d DataTime) MarshalJSON() ([]byte, error) {
//	return []byte(d.Time.Format("2006-01-02 15:04:05")), nil
//}
//
//func (d *DataTime) UnmarshalJSON(data []byte) error {
//
//	s := strings.Trim(string(data), "\"")
//	t, err := time.Parse("2006-01-02 15:04:05", s)
//	if err != nil {
//		return err
//	}
//	d.Time = t
//	return nil
//}
