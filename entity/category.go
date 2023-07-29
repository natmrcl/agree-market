package entity

type Category struct {
	ID   int64  `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}
