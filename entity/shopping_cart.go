package entity

type Shopping_Cart struct {
	ID        int64   `gorm:"primary_key" json:"id"`
	UserID    uint    `json:"-"`
	ProductID uint    `json:"-"`
	User      User    `gorm:"foreignKey:UserID" json:"user"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
}
