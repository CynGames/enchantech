package models

type Article struct {
	ID             string `gorm:"primaryKey" json:"id"`
	PublisherID    uint   `json:"publisherId"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	ImageUrl       string `json:"imageUrl"`
	ParseAttempted bool   `json:"parseAttempted"`
}
