package controller

import "bluebell/models"

// 接口文档使用的model
// 接口文档返回的数据格式是一致的，具体的data类型不一致

type _ResponsePostList struct {
	Code    ResCode                 `json:"code"`    // code 业务响应状态码
	Message string                  `json:"message"` // message 提示信息
	Data    []*models.ApiPostDetail `json:"data"`    // data 数据
}
