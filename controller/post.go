package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子的处理函数
// 这段代码是用于创建新帖子的处理函数。它遵循以下步骤：
// 1. 从请求的 JSON 数据中获取参数并进行校验，使用了 JSON 校验器和绑定标签。
// 2. 从请求上下文中获取当前用户的 ID。如果用户未登录，则返回需要登录的错误响应。
// 3. 将帖子的作者 ID 设置为当前用户的 ID。
// 4. 调用逻辑层的 CreatePost 函数在数据库中创建帖子。
// 5. 如果在帖子创建过程中出现错误，将记录错误日志并返回服务器繁忙的错误响应。
// 6. 如果帖子成功创建，将返回成功的响应。
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数的校验
	// c.ShouldBindJSON() validator --> binding tag
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("create post with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从 c 取到当前发请求的用户ID
	user_id, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = user_id
	// 2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// CreatePostDetailHandler 获取帖子详情的处理函数
func CreatePostDetailHandler(c *gin.Context) {
	// 1. 获取参数（从URL中获取帖子的ID）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 根据ID取出帖子数据（查数据库）
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic get post by pid failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}
