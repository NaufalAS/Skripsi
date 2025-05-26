package entity

import "skripsi/model/domain"

type UserEntity struct {
	UserID int    `json:"id"`
	Name   string `json:"name"`
	Profile string `json:"foto"`
	Email string `json:"email"`
	NoTelepon string `json:"no_telepon"`
	Alamat string `json:"alamat"`
}

func ToUserEntity(user domain.AppUser) UserEntity {
	return UserEntity{
		UserID: user.ID,
		Name: user.Name,
		Profile : user.Profile,
		Email: user.Email,
		Alamat: user.Alamat,
		NoTelepon: user.NoTelepon,
	}
}

func ToUserListEntity(admins []domain.AppUser) []UserEntity {
	var userData []UserEntity

	for _, admin := range admins {
		userData = append(userData, ToUserEntity(admin))

	}
	return userData
}