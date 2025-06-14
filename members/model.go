package members

import(
	"time"
	"gorm.io/gorm"
)

type Status string

const(
	StatusComplete Status = "COMPLETE"
	StatusIncomplete Status = "INCOMPLETE"
)

type Member struct {
	
	ID uint `gorm:"primaryKey;autoIncrement"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName string `gorm:"not null" json:"last_name"`
	Email string    `gorm:"not null;unique" json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Cohort string     `gorm:"not null" json:"cohort"`
	Year int   `gorm:"not null" json:"year"`
	Session1 *bool `gorm:"default:false" json:"session1"`
	Session2 *bool `gorm:"default:false" json:"session2"`
	Session3 *bool `gorm:"default:false" json:"session3"`
	Session4 *bool `gorm:"default:false" json:"session4"`
	Session1CompletionTime time.Time
	Session2CompletionTime time.Time
	Session3CompletionTime time.Time
	Session4CompletionTime time.Time
	Status Status `gorm:"type:status;default:INCOMPLETE;not null" json:"status"`
	CreatedAt time.Time
}

func (m *Member) BeforeSave(tx *gorm.DB) error {
	// Check if all sessions are complete
	if (m.Session1 != nil && *m.Session1) &&
	(m.Session2 != nil && *m.Session2) &&
	(m.Session3 != nil && *m.Session3) &&
	(m.Session4 != nil && *m.Session4) {
    	m.Status = Status("COMPLETE")
	} else {
		m.Status = Status("INCOMPLETE")
	}
	return nil

}