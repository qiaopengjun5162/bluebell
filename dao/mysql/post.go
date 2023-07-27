package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `INSERT INTO post (post_id, title, content, author_id, community_id) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostById 根据帖子ID查询单个帖子信息
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `SELECT post_id, title, content, author_id, community_id, create_time FROM post WHERE post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(pageNumber, pageSize int64) (posts []*models.Post, err error) {
	sqlStr := `SELECT post_id, title, content, author_id, community_id, create_time FROM post 
                                                                     ORDER BY create_time DESC LIMIT ?,?`
	posts = make([]*models.Post, 0, pageSize)
	err = db.Select(&posts, sqlStr, (pageNumber-1)*pageSize, pageSize)
	return
}

// GetPostListByIDs 根据给定的ID列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	postList = make([]*models.Post, 0, len(ids))
	sqlStr := `SELECT post_id, title, content, author_id, community_id, create_time FROM post 
                                                                     WHERE post_id in (?) ORDER BY FIND_IN_SET(post_id, ?)`
	// https://www.liwenzhou.com/posts/Go/sqlx/
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
