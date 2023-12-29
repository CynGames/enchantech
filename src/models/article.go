package models

type Article struct {
	ID             string `gorm:"primaryKey"`
	PublisherID    uint
	Title          string
	Description    string
	ImageUrl       string
	ParseAttempted bool
}
