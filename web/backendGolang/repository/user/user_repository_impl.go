package userrepo

import "skripsi/model/domain"

type UserRepository interface {
	SaveUser(user domain.AppUser) (domain.AppUser, error)
	LoginUser(name string)(*domain.AppUser, error)
	GetListUser() ([]domain.AppUser, error)
	GetUserById(id int) (domain.AppUser, error)
	UpdateId(idUser int, user domain.AppUser) (domain.AppUser, error)
	DeleteId(iduser int) error
	UpdatePassword(id int, updatedUser domain.AppUser) (domain.AppUser, error) 
}
