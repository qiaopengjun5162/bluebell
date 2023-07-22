package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
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
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	// 1. 查询帖子信息
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById error ", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	// 2. 根据作者ID查找作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql get user by id failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}
	// 3. 根据社区ID查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}
	// 4. 接口数据拼接
	data = &models.ApiPostDetail{
		AuthorName:      user.UserName,
		Post:            post,
		CommunityDetail: community,
	}
	return
}
