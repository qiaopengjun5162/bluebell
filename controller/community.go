package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ---- 跟社区相关的 ----

// CommunityHandler CommunityHandler函数用于处理请求，返回所有社区的列表数据。
// 该函数的主要逻辑是调用logic包中的GetCommunityList函数来查询社区数据，并将查询结果以列表的形式返回。
// 如果查询过程中出现错误，会记录日志并返回服务器繁忙的错误信息。如果查询成功，会将查询结果以成功的响应返回给客户端。
func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区 （community_id, community_name）以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	// 1. get community id 获取URL参数 社区id
	communityId := c.Param("community_id")
	// 2. 将字符串类型的 communityId 转换为 int64 类型，并返回转换后的值
	id, err := strconv.ParseInt(communityId, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		zap.L().Error("Community id invalid", zap.Error(err))
		return
	}
	// 3. 根据ID获取社区详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}
	ResponseSuccess(c, data)
}
