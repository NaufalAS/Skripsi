package web

type LoginUserRequest struct {
	// Foto string `validate:"required" jsonL:"foto"`
	Name    string `validate:"required" json:"name"`
	Password string `validate:"required" json:"password"`
}

type UpdateUserRequest struct {
	Foto string `validate:"required" json:"foto"`
	Name     string `validate:"required" json:"name"`
	Email string `validate:"required" json:"email"`
	NoTelepon string `validate:"required" json:"no_telepon"`
	Alamat string `validate:"required" json:"alamat"`
}

type UpdatePasswordRequest struct{
	Password string `validate:"required" json:"password_lama"`
	NewPassword string `validate:"required" json:"password_baru"`
}