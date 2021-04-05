package models

type UserID struct {
	Id uint `uri:"id" json:"id"`
}

type UserBase struct {
	UserName           string `json:"user_name" binding:"required,alphanum,min=4,max=255"`
	FirstName          string `json:"first_name" binding:"max=1024"`
	MiddleName         string `json:"middle_name" binding:"max=1024"`
	LastName           string `json:"last_name" binding:"max=1024"`
	Email              string `json:"email" binding:"email"`
	PrimaryPhoneNumber string `json:"primary_phone_number"`
}

type UserIncoming struct {
	UserBase
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type UserOutgoing struct {
	UserID
	UserBase
}

type UserAccount struct {
	UserID
	UserBase
	PasswordHash string `json:"password_hash"`
}
