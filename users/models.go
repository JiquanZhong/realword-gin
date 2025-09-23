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
	if len(password) < 8 || len(password) > 32 {
		return errors.New("password length should be between 8 and 32")
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

func SaveOne(data interface{}) error {
	db := common.GetDB()
	return db.Save(data).Error
}

func (u *UserModel) Update(data interface{}) error {
	db := common.GetDB()
	err := db.Model(u).Updates(data).Error
	return err
}

// u关注v
func (u *UserModel) Following(v UserModel) error {
	db := common.GetDB()
	var following FollowModel
	err := db.FirstOrCreate(&following, &FollowModel{
		FollowingID:  v.ID,
		FollowedById: u.ID,
	}).Error
	return err
}

// u是否关注了v
func (u *UserModel) isFollowing(v UserModel) bool {
	db := common.GetDB()
	var follow FollowModel
	db.Where(FollowModel{
		FollowingID:  v.ID,
		FollowedById: u.ID,
	}).First(&follow)
	return follow.ID != 0
}

// u取消关注v
func (u *UserModel) unFollowing(v UserModel) error {
	db := common.GetDB()
	err := db.Where(FollowModel{
		FollowingID:  v.ID,
		FollowedById: u.ID,
	}).Delete(FollowModel{}).Error
	return err
}

// 获取u的关注列表
func (u *UserModel) GetFollowings() []UserModel {
	db := common.GetDB()
	tx := db.Begin()
	var followings []UserModel
	var follows []FollowModel
	tx.Where(FollowModel{
		FollowedById: u.ID,
	}).Find(&follows)
	for _, follow := range follows {
		// 先获取关注关系模型，拿到每个关注关系的FollowingID
		// 再根据FollowingID去用户表里查找对应的用户
		// 这样做的原因是gorm不支持多级嵌套预加载
		var userModel UserModel
		tx.Model(&follow).Related(&userModel, "Following")
		followings = append(followings, userModel)
	}
	tx.Commit()
	return followings
}
