package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) {
	// 1. 判断用户是否存在
	mysql.QueryUserByUsername()
	// 2. 生成 UID
	snowflake.GenID()
	// 3. 密码加密

	// 4. 保存到数据库
	mysql.InsertUser()
}
