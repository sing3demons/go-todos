package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null" json:"-"`
	FirstName string
	LastName  string
	Avatar    string
	Role      string `gorm:"default:'Member';not null"`
	Todo      []Todo `gorm:"foreignKey:UserID"`
}

//Promote - update user --> editor
func (u *User) Promote() {
	u.Role = "Editor"
}

//Demote - Change user --> editor
func (u *User) Demote() {
	u.Role = "Member"
}

func (u *User) GenerateEncryptedPassword() string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	return string(hash)
}
