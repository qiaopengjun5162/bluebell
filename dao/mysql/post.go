package mysql

import "bluebell/models"

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `INSERT INTO post (post_id, title, content, author_id, community_id) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostById
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `SELECT post_id, title, content, author_id, community_id, create_time FROM post WHERE post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return
}
