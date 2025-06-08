package users

import(
	"time"
)

type Role string

const(
	SuperUserRole Role = "SUPERUSER"
	UserRole Role = "USER"
)

type User struct {
	
	ID uint `gorm:"primaryKey;autoIncrement"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName string `gorm:"not null" json:"last_name"`
	Email string    `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null;unique" json:"password"`
	Role Role `gorm:"type:role;default:'USER';not null" json:"role"`
	CreatedAt time.Time
}