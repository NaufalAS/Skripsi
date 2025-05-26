package domain

import "time"

type AppUser struct {
	ID        int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string `gorm:"column:name"`
	Profile      string `gorm:"column:profile"`
	Password  string `gorm:"column:password"`
	Email string `gorm:"column:email"`
	NoTelepon string `gorm:"column:no_telepon"`
	Alamat string `gorm:"alamat"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (AppUser) TableName() string {
	return "appuser"
}