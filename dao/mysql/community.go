package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

// GetCommunityList 获取社区列表
func GetCommunityList() (communityList []*models.Community, err error) {
	// 该函数用于获取社区列表。它通过查询数据库中的community表，获取社区的ID和名称，并将结果存储在communityList变量中。
	// 如果查询成功，将返回communityList和nil；如果查询结果为空，则会记录一条警告日志并返回nil；如果查询出现错误，则将返回错误信息。
	sqlStr := "SELECT community_id, community_name FROM community"
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in database")
			err = nil
		}
	}
	return
}

// GetCommunityDetailByID 根据ID查询社区详情
func GetCommunityDetailByID(id int64) (community_detail *models.CommunityDetail, err error) {
	// 根据给定的id获取社区详情。函数接受一个int64类型的参数id，返回一个名为community_detail的指针类型变量和一个错误变量err。
	// 函数首先创建一个models.CommunityDetail类型的指针变量community_detail。
	// 然后定义一个SQL查询语句，查询社区表中的community_id、community_name、introduction和create_time字段，条件是community_id等于给定的id。
	// 接着使用db.Get方法执行查询，并将结果赋值给community_detail变量。
	// 如果查询出错，会根据错误类型进行处理，如果是没有找到对应的行，则会记录警告日志，并将错误变量err设置为ErrorInvalidID。
	// 最后，函数返回community_detail和err变量。
	community_detail = new(models.CommunityDetail)
	sqlStr := "SELECT community_id, community_name, introduction, create_time FROM community WHERE community_id = ?"
	if err = db.Get(community_detail, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("get community detail from database failed", zap.Error(err))
			err = ErrorInvalidID
		}
	}
	return community_detail, err
}
