package userrepo

import (
	"skripsi/model/domain"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{
		DB: db,
	}
}

func (repo *AuthRepositoryImpl) SaveUser(user domain.AppUser) (domain.AppUser, error) {
	err := repo.DB.Create(&user).Error
	if err != nil {
		return domain.AppUser{}, err
	}
	return user, nil
}

func(repo *AuthRepositoryImpl) LoginUser(name string)(*domain.AppUser, error){
	user := new(domain.AppUser)
	if err := repo.DB.Where("name = ?", name).Take(user).Error; err != nil{
		return nil, err
	}

	return user, nil
}
func (repo *AuthRepositoryImpl) GetListUser() ([]domain.AppUser, error) {
	var user []domain.AppUser
	err := repo.DB.Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func(repo *AuthRepositoryImpl) GetUserById(id int)(domain.AppUser, error){
	var user domain.AppUser

	if err := repo.DB.Find(&user,"id = ?", id).Error; err != nil{
		return user, err
	}

	return user, nil
}

func (repo *AuthRepositoryImpl) UpdateId(idUser int, user domain.AppUser) (domain.AppUser, error) {
    if err := repo.DB.Model(&domain.AppUser{}).Where("id = ?", idUser).Updates(user).Error; err != nil {
        return domain.AppUser{}, err
    }
    return user, nil
}

func (repo *AuthRepositoryImpl) DeleteId(iduser int)error{
	if err := repo.DB.Delete(&domain.AppUser{}, iduser).Error; err != nil {
        return err
    }
    return  nil
}

func (repo *AuthRepositoryImpl) UpdatePassword(id int, updatedUser domain.AppUser) (domain.AppUser, error) {
    if err := repo.DB.Model(&domain.AppUser{}).Where("id = ?", id).Update("password", updatedUser.Password).Error; err != nil {
        return domain.AppUser{}, err
    }
    var updated domain.AppUser
    if err := repo.DB.Where("id = ?", id).First(&updated).Error; err != nil {
        return domain.AppUser{}, err
    }
    return updated, nil
}