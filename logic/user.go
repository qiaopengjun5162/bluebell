package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"fmt"
)

// 存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 2. 生成 UID
	userID := snowflake.GenID()
	fmt.Printf("generation started with userID: %v\n", userID)
	// 3. 构造一个 User 实例
	user := &models.User{
		UserID:   userID,
		UserName: p.Username,
		Password: p.Password,
	}
	// 4. 保存到数据库
	return mysql.InsertUser(user)
}
