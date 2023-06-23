package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// 把每一步数据库操作封装成函数
// 待 Logic 层根据业务需求调用

const secret = "qiaopengjun.com"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
)

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `SELECT count(user_id) FROM user WHERE username = ?`
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		// 用户已存在的错误
		return ErrorUserExist
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

// Login
func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户登录的密码
	sqlStr := `SELECT user_id, username, password FROM user WHERE username = ?`
	if err = db.Get(user, sqlStr, user.UserName); err != nil {
		return err // 查询数据库失败
	}
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	//	判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}
