package users

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/jiquanzhong/realword-gin/common"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	ID           uint    `gorm:"primary_key"`
	Username     string  `gorm:"column:username"`
	Email        string  `gorm:"column:email;unique_index"`
	Bio          string  `gorm:"column:bio;size:1024"`
	Image        *string `gorm:"column:image"`
	PasswordHash string  `gorm:"column:password;not null"`
}

type FollowModel struct {
	gorm.Model
	Following    UserModel
	FollowingID  uint
	FollowedBy   UserModel
	FollowedById uint
}

func AutoMigrate() {
	db := common.GetDB()

	db.AutoMigrate(UserModel{})
	db.AutoMigrate(&FollowModel{})
}

func (u *UserModel) setPassword(password string) error {
	if len(password) < 6 || len(password) > 32 {
		return errors.New("password length should be between 6 and 32")
	}

	bytePassword := []byte(password)
	passwordHass, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHass)
	return nil
}

func (u *UserModel) checkPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func FindOneUser(condition interface{}) (UserModel, error) {
	db := common.GetDB()
	var model UserModel
	err := db.Where(condition).First(&model).Error
	return model, err
}
