package models

type User struct {
	ID        uint `gorm:"primaryKey"`
	Username  string
	Password  string
	Favorites []Article `gorm:"many2many:favorites;"`
}
