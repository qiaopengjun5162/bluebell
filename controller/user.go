package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid parameters", zap.Error(err))
		// 判断 err 是否是 validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			c.JSON(http.StatusOK, gin.H{
				"message": err.Error(),
			})
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		c.JSON(http.StatusOK, gin.H{
			"message": removeTopStruct(errs.Translate(trans)),
		})
		fmt.Printf("paramSignUp error %v\n", err)
		return
	}

	fmt.Printf("signUp params: %v\n", p)
	// 2. 业务处理
	// 结构体是值类型，字段很多的时候，会有性能影响，故最好传指针
	if err := logic.SignUp(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "registration failed",
		})
		return
	}
	// 3. 返回响应
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func LoginHandler(c *gin.Context) {
	// 1. 获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid parameters", zap.Error(err))
		// 判断 err 是否是 validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": removeTopStruct(errs.Translate(trans)), // 翻译错误
		})
		return
	}
	// 2. 业务逻辑处理
	if err := logic.Login(p); err != nil {
		zap.L().Error("logic Login failed", zap.String("username", p.Username), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Login failed, username or password incorrect",
		})
		return
	}
	// 3. 返回响应
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}
