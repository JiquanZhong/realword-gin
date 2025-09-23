package users

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/jiquanzhong/realword-gin/common"
	"github.com/stretchr/testify/assert"
)

var image_url = "https://golang.org/doc/gopher/frontpage.png"
var test_db *gorm.DB

func newUserModel() UserModel {
	return UserModel{
		ID:           1,
		Username:     "asd123!@#ASD",
		Email:        "wzt@g.cn",
		Bio:          "heheda",
		Image:        &image_url,
		PasswordHash: "",
	}
}

func userModelMocker(n int) []UserModel {
	var users []UserModel
	for i := 0; i < n; i++ {
		image := fmt.Sprintf("http://image/%v.jpg", i)
		user := UserModel{
			Username: "user" + common.RandString(10),
			Email:    common.RandString(10) + "@test.com",
			Bio:      "This is user" + common.RandString(10),
			Image:    &image,
		}
		user.setPassword("123456789")
		test_db.Create(&user)
		users = append(users, user)
	}
	return users
}

func TestUserModel(t *testing.T) {
	resetDBWithMock()
	asserts := assert.New(t)
	var userModel UserModel
	var err error

	userModel = newUserModel()
	err = userModel.checkPassword("")
	asserts.Error(err, "empty password should not match")

	userModel = newUserModel()
	err = userModel.setPassword("")
	asserts.Error(err, "empty password should not set")

	userModel = newUserModel()
	err = userModel.setPassword("asd123!@#ASD")
	asserts.NoError(err, "password should be set successfully")
	asserts.Len(userModel.PasswordHash, 60, "password hash should be 60 characters long")

	err = userModel.checkPassword("asd123!@#ASD!")
	asserts.Error(err, "wrong password should not match")

	err = userModel.checkPassword("asd123!@#ASD")
	asserts.NoError(err, "correct password should match")

	users := userModelMocker(3)
	a := users[0]
	b := users[1]
	c := users[2]
	asserts.Equal(0, len(a.GetFollowings()), "a should not follow anyone")
	asserts.Equal(false, a.isFollowing(b), "a should not follow b")
	a.Following(b)
	asserts.Equal(1, len(a.GetFollowings()), "a should follow 1 person now")
	asserts.Equal(true, a.isFollowing(b), "a should follow b now")
	a.Following(c)
	asserts.Equal(2, len(a.GetFollowings()), "a should follow 2 person now")
	asserts.Equal(true, a.isFollowing(c), "a should follow c now")
	asserts.Equal(false, b.isFollowing(c), "b should not follow c")
	asserts.Equal(false, b.isFollowing(a), "b should not follow a")
	asserts.Equal(false, c.isFollowing(b), "c should not follow b")
	asserts.Equal(false, c.isFollowing(a), "c should not follow a")

}

func resetDBWithMock() {
	common.TestDBFree(test_db)
	test_db = common.TestDBInit()
	AutoMigrate()
}

func HeaderTokenMock(req *http.Request, u uint) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", common.GenToken(u)))
}
