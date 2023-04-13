package models

type User struct {
	ID int `json:"id"`
	// Role      string    `json:"role" gorm:"type:string"`
	Name     string `json:"name" gorm:"type: varchar(255)"`
	Email    string `json:"email" gorm:"type: varchar(255)"`
	Password string `json:"-" gorm:"type: varchar(255)"`
	// CreatedAt time.Time `json:"-"`
	// UpdatedAt time.Time `json:"-"`
}

type UserResponse struct {
	ID int `json:"id"`
	// Role string `json:"string"`
	Name string `json:"name"`
}

func (UserResponse) TableName() string {
	return "users"
}
