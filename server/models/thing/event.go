package thing

type Event struct {
	ID   string `json:"-" gorm:"primaryKey"`
	Name string `json:"name"`
}
