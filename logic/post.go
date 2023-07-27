package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	//	1. 生成 post ID
	p.ID = snowflake.GenID()
	//	2. 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID)
	//	3. 返回
	return
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

// GetPostList 获取帖子列表
func GetPostList(pageNumber, pageSize int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(pageNumber, pageSize)
	if err != nil {
		zap.L().Error("mysql.GetPostList error ", zap.Error(err))
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 2. 根据作者ID查找作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql get user by id failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		// 3. 根据社区ID查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		// 4. 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostList2 根据前端传来的参数动态的获取帖子列表
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 1. 去Redis查询ID列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis get post ids in order return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 2. 根据ID去MySQL数据库查询帖子详细信息
	// 返回的数据按照给定的ID的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostListByIDs", zap.Any("posts", posts))
	// 3. 将帖子的作者及分区信息查询出来填充到帖子中
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 2. 根据作者ID查找作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql get user by id failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		// 3. 根据社区ID查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		// 4. 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}
