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
// @Summary 创建帖子的接口
// @Tags 帖子相关接口(api分组展示使用的)
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Security ApiKeyAuth
// @Router /post [post]
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
	userId, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userId
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
// @Summary 获取帖子详情的接口
// @Tags 帖子相关接口(api分组展示使用的)
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Security ApiKeyAuth
// @Router /post/:id [get]
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

// GetPostListHandler 获取帖子列表的处理函数
// @Summary 获取帖子列表的接口
// @Tags 帖子相关接口(api分组展示使用的)
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Security ApiKeyAuth
// @Router /posts [get]
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	pageNumber, pageSize := getPageInfo(c)
	// 1. 获取数据
	data, err := logic.GetPostList(pageNumber, pageSize)
	if err != nil {
		zap.L().Error("get post list failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表接口
// 根据前端传来的参数动态的获取帖子列表
// 按创建时间排序 或者 按分数排序
// 1. 获取请求的 query String 参数
// 2. 去Redis查询ID列表
// 3. 根据ID去数据库查询帖子详细信息
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口(api分组展示使用的)
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	// GET请求参数(query string)：/api/v1/posts2?page=1&size=10&order=time
	// 获取分页参数
	// 初始化结构体时指定初始参数
	// 定义结构体参数变量时可以指定默认值
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, // magic string
	}
	// c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
	// c.ShouldBindJSON() 如果请求中携带的是JSON格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("get post list handler2 failed, with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 1. 获取数据
	data, err := logic.GetPostListNew(p)
	if err != nil {
		zap.L().Error("get post list failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2. 返回响应
	ResponseSuccess(c, data)
}

//// GetCommunityPostListHandler 根据社区去查询帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	p := &models.ParamCommunityPostList{
//		ParamPostList: &models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.OrderTime, // magic string
//		},
//	}
//	// c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
//	// c.ShouldBindJSON() 如果请求中携带的是JSON格式的数据，才能用这个方法获取到数据
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("ParamCommunityPostList failed, with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//	// 1. 获取数据
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("get post list failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	// 2. 返回响应
//	ResponseSuccess(c, data)
//}
