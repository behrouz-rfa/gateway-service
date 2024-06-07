package domain

import "time"

type User struct {
	Base
	Name     string    `gorm:"size:255;" json:"name"`
	Email    string    `gorm:"size:255;not null;unique" json:"email"`
	Password string    `json:"-"`
	JWTToken *JWTToken `json:"token" gorm:"-"`
}

type UserInput struct {
	ID        string `json:"id" `
	Name      string `json:"username" `
	Email     string `json:"firstname"`
	Password  string `json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserLoginInput struct {
	Email    string `json:"email" `
	Password string `json:"password"`
}
