package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

const secret = "qiaopengjun.com"

// 把每一步数据库操作封装成函数
// 待 Logic 层根据业务需求调用

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `SELECT count(user_id) FROM user WHERE username = ?`
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		// 用户已存在的错误
		return errors.New("user already")
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 执行SQL 语句入库
	sqlStr := `INSERT INTO user (user_id, username, password) VALUES (?, ?, ?)`
	_, err = db.Exec(sqlStr, user.UserID, user.UserName, user.Password)
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
