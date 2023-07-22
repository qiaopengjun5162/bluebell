package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	//	1. 生成 post ID
	p.ID = snowflake.GenID()
	//	2. 保存到数据库
	return mysql.CreatePost(p)
	//	3. 返回
}

// GetPostById 根据帖子ID查询帖子详情数据
func GetPostById(pid int64) (*models.Post, error) {
	return mysql.GetPostById(pid)
}
