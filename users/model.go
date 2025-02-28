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
	Name string `gorm:"not null" json:"name"`
	Email string    `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null;unique" json:"password"`
	Role Role `gorm:"type:role;default:'USER';not null" json:"role"`
	CreatedAt time.Time
}