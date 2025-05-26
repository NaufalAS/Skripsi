package userservice

import (
	"errors"
	"fmt"
	"skripsi/helper"
	"skripsi/model/domain"
	"skripsi/model/entity"
	"skripsi/model/web"
	userrepo "skripsi/repository/user"
	"os"
    "path/filepath"
    "io"
    "mime/multipart"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userrepo userrepo.UserRepository
}

func NewSektorUsahaService(	userrepo userrepo.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userrepo: userrepo,
	}
}

func (service *UserServiceImpl) SaveUser(request web.LoginUserRequest) (map[string]interface{}, error) {
	passHash, errHash := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if errHash != nil {
		return nil, errHash
	}

	request.Password = string(passHash)
	newUser := domain.AppUser{
		Name:     request.Name,
		Password: request.Password, // Tanpa hashing untuk sementara
	}


	saveUser, errUser := service.userrepo.SaveUser(newUser)
	if errUser != nil {
		fmt.Println("Database error:", errUser)
		return nil, errUser
	}

	data := helper.ResponseToJson{
		"id":    saveUser.ID,
		"name": saveUser.Name,
	}
	return data, nil
}

func (service *UserServiceImpl) Login(name, password string) (helper.ResponseToJson, error) {
	buyer, getBuyerErr := service.userrepo.LoginUser(name)
	if getBuyerErr != nil {
		return nil, errors.New("wrong name or password")
	}

	if checkPasswordErr := bcrypt.CompareHashAndPassword([]byte(buyer.Password), []byte(password)); checkPasswordErr != nil {
		return nil, errors.New("wrong name or password")
	}

	loginResponse, loginErr := helper.Login(buyer.ID, buyer.Name)
	if loginErr != nil {
		return nil, loginErr
	}

	return helper.ResponseToJson{
		"token":      loginResponse["token"],
		"expires_at": loginResponse["expires_at"],
	}, nil
}

func (service *UserServiceImpl) GetUser() ([]entity.UserEntity, error) {
	getUser, err := service.userrepo.GetListUser()
	if err != nil {
		return nil, err
	}

	userEntities := entity.ToUserListEntity(getUser)

	return userEntities, nil
}

func(service *UserServiceImpl) GetUserById(id int)(entity.UserEntity, error){
	getUserService, err := service.userrepo.GetUserById(id)
	if err != nil {
		return entity.UserEntity{}, err
	}

	return entity.ToUserEntity(getUserService), nil
}


// Fungsi untuk menyimpan file foto
func saveFile(file multipart.File, destination string) error {
	// Membuat file baru di lokasi tujuan
	out, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// Salin konten dari file upload ke file baru
	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	return nil
}

// Fungsi untuk memperbarui data pengguna
func (service *UserServiceImpl) UpdateUserId(Id int, req web.UpdateUserRequest, file multipart.File) (map[string]interface{}, error) {
	// Cek apakah user ada
	_, errUser := service.userrepo.GetUserById(Id)
	if errUser != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Jika ada foto baru yang diupload
	if file != nil {
		// Ambil data user lama
		oldUser, _ := service.userrepo.GetUserById(Id)
		oldProfilePath := filepath.Join("public", "profile", filepath.Base(oldUser.Profile))

		// Hapus foto lama jika ada
		if _, err := os.Stat(oldProfilePath); err == nil {
			if err := os.Remove(oldProfilePath); err != nil {
				return nil, fmt.Errorf("failed to delete old profile picture: %v", err)
			}
		}

		// Tentukan nama file baru
		newProfileName := fmt.Sprintf("%d-%s", Id, filepath.Base(req.Foto))
		newProfilePath := filepath.Join("public", "profile", newProfileName)

		// Simpan file
		if err := saveFile(file, newProfilePath); err != nil {
			return nil, fmt.Errorf("failed to save profile picture: %v", err)
		}

		// Set path untuk disimpan ke database
		req.Foto = "/public/profile/" + newProfileName
	}

	// Persiapkan data untuk update
	updatedUser := domain.AppUser{
		ID:        Id,
		Name:      req.Name,
		Profile:   req.Foto,
		Email:     req.Email,
		NoTelepon: req.NoTelepon,
		Alamat:    req.Alamat,
	}

	// Update ke repo
	result, errUpdate := service.userrepo.UpdateId(Id, updatedUser)
	if errUpdate != nil {
		return nil, errUpdate
	}

	// Response akhir
	response := helper.ResponseToJson{
		"name":       result.Name,
		"profile":    result.Profile,
		"email":      result.Email,
		"no_telepon": result.NoTelepon,
		"alamat":     result.Alamat,
	}
	return response, nil
}




func (service *UserServiceImpl) DeleteProduk(id int) error {
	// Cari produk berdasarkan ID
	_, err := service.userrepo.GetUserById(id)
	if err != nil {
		return err
	}

	// Hapus produk dari database
	return service.userrepo.DeleteId(id)
}

func (service *UserServiceImpl) UpdatePassword(Id int, oldpassword string, newpassword string, req web.UpdatePasswordRequest) (map[string]interface{}, error) {
    // 1. Ambil user berdasarkan ID untuk verifikasi password lama
    existingUser, err := service.userrepo.GetUserById(Id)
    if err != nil {
        return nil, fmt.Errorf("user not found")
    }

    // 2. Verifikasi password lama dengan bcrypt
    err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(oldpassword))
    if err != nil {
        return nil, fmt.Errorf("password lama salah")
    }

    // 3. Verifikasi jika password lama dan password baru sama
    if oldpassword == newpassword {
        return nil, fmt.Errorf("password baru tidak boleh sama dengan password lama")
    }

    // 4. Enkripsi password baru
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("gagal mengenkripsi password baru")
    }

    // 5. Update password baru
    updatedUser := domain.AppUser{
        ID:       Id,
        Password: string(hashedPassword),
    }

    // 6. Update data user
    result, errUpdate := service.userrepo.UpdatePassword(Id, updatedUser)
    if errUpdate != nil {
        return nil, errUpdate
    }

    // 7. Membuat response JSON
    response := helper.ResponseToJson{
        "name":     result.Name,
        "password": result.Password,
    }

    return response, nil
}
