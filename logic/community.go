package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查询数据库中community表中的所有信息并返回
	return mysql.GetCommunityList()
}
