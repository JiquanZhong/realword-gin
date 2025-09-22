package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionDatabase(t *testing.T) {
	asserts := assert.New(t)
	db := Init()
	_, err := os.Stat("./../gorm.db")
	asserts.NoError(err, "Db should be created")
	asserts.NoError(db.DB().Ping(), "Db should be able to connect")

	connection := GetDB()
	asserts.NoError(connection.DB().Ping(), "Db should be able to connect")
	db.Close()

	// 确保完全关闭数据库连接
	db.DB().Close()

	// 删除数据库文件
	err = os.Remove("./../gorm.db")
	asserts.NoError(err)

	// 验证文件已被删除
	_, err = os.Stat("./../gorm.db")
	asserts.Error(err, "Database file should be deleted")

	// SQLite会自动创建新的数据库文件，所以连接应该成功
	db = Init()
	asserts.NoError(db.DB().Ping(), "New db should be able to connect")
	db.Close()
}
