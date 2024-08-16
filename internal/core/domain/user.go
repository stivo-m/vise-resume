package domain

import "time"

type User struct {
	Base
	FullName        string          `gorm:"size:100;"  json:"full_name"`
	Email           string          `gorm:"size:150;not null;unique" json:"email"`
	EmailVerifiedAt *time.Time      `gorm:"default:null" json:"email_verified_at"`
	Password        Password        `json:"-"`
	Tokens          []Token         `json:"-"`
	Verifications   []Verifications `json:"-"`
}

type Password struct {
	Base
	Value  string `gorm:"size:150;not null" json:"email"`
	UserId string `gorm:"type:uuid;not null;index;"`
}

type Token struct {
	Base
	AccessToken string `gorm:"type:varchar(500);not null" json:"access_token"`
	UserId      string `gorm:"type:uuid;not null;index;"`
}

type Verifications struct {
	Base
	Type   string `gorm:"size:20;not null"`
	UserId string `gorm:"type:uuid;not null;index;"`
	Code   string `gorm:"size:20;not null;"`
}
