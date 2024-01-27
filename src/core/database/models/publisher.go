package models

type Publisher struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"name"`
	RSS      string    `json:"rss"`
	Articles []Article `gorm:"foreignKey:PublisherID"`
}
