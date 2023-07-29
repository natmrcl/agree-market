package entity

type Product struct {
	ID         int64    `gorm:"primary_key" json:"id"`
	Name       string   `json:"name"`
	Brand      string   `json:"brand"`
	Price      float64  `json:"price"`
	Seller     string   `json:"seller"`
	ImageURL   string   `json:"image_url"`
	CategoryID uint     `json:"-"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`
}
