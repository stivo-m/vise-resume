package domain

import (
	"time"

	"github.com/lib/pq"
)

type Resume struct {
	Base
	UserId      string         `gorm:"type:uuid;not null;index;"`
	Score       int            `gorm:"default:0"`
	DocumentUrl string         `gorm:"type:text"`
	Skills      pq.StringArray `gorm:"type:text[]"`
	Experiences []WorkExperience
	Education   []WorkExperience
}

type WorkExperience struct {
	Base
	ResumeId    string `gorm:"type:uuid;not null;index;"`
	CompanyName string `gorm:"size:255; not null"`
	Role        string `gorm:"size:255; not null"`
	StartDate   time.Time
	EndDate     time.Time
}

type Education struct {
	Base
	ResumeId   string `gorm:"type:uuid;not null;index;"`
	SchoolName string `gorm:"size:255; not null"`
	Course     string `gorm:"size:255; not null"`
	StartDate  time.Time
	EndDate    time.Time
}
