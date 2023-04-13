package models

type Menu struct {
	ID          int    `json:"id"`
	Name        string `json:"name" gorm:"type: varchar(255)"`
	Description string `json:"desc" gorm:"type: text"`
	Price       int    `json:"price" gorm:"type: int"`
	Category    string `json:"category" gorm:"type:varchar(255)"`
}
