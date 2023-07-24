package logic

import "bluebell/models"

// 推荐阅读
// 基于用户投票的相关算法：http://www.ruanyifeng.com/blog/2012/03/

// 本项目使用简化版的投票分数
// 投一票就加432分  86400/200  -> 200 张赞成票可以给你的帖子续一天 -> 《Redis实战》

/* 投票的几种情况：
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票
	2. 之前投反对票，现在改投赞成票
direction=0时，有两种情况：
	1. 之前投过赞成票，现在要取消投票
	2. 之前投过反对票，现在要取消投票
direction=-1时，有两种情况：
	1. 之前没有投过票，现在投反对票
	2. 之前投赞成票，现在改投反对票

投票的限制：
每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1. 到期之后将Redis中保存的赞成票及反对票数存储到MySQL表中
	2. 到期之后删除 KeyPostVotedZSetPrefix
*/

// VoteForPost 帖子投票的函数
func VoteForPost(userID int64, p *models.ParamVoteData) {
	// 1.
	// 2.
}
