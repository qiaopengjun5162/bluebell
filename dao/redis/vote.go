package redis

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

// 本项目使用简化版的投票分数
// 投一票就加432分  86400/200  -> 200 张赞成票可以给你的帖子续一天 -> 《Redis实战》

/* 投票的几种情况：
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票      -->  更新分数和投票记录  差值的绝对值：1  +432
	2. 之前投反对票，现在改投赞成票      -->  更新分数和投票记录  差值的绝对值：2  +432*2
direction=0时，有两种情况：
	1. 之前投过反对票，现在要取消投票    -->  更新分数和投票记录  差值的绝对值：1  +432
	2. 之前投过赞成票，现在要取消投票    -->  更新分数和投票记录  差值的绝对值：1  -432
direction=-1时，有两种情况：
	1. 之前没有投过票，现在投反对票      -->  更新分数和投票记录  差值的绝对值：1  -432
	2. 之前投赞成票，现在改投反对票      -->  更新分数和投票记录  差值的绝对值：2  -432*2

投票的限制：
每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1. 到期之后将Redis中保存的赞成票及反对票数存储到MySQL表中
	2. 到期之后删除 KeyPostVotedZSetPrefix
*/

const (
	oneWeekInSeconds = 7 * 24 * 60 * 60
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
)

func CreatePost(postID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// 事务
	// TxPipeline acts like Pipeline, but wraps queued commands with MULTI/EXEC.
	_, err := client.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		// 帖子时间
		// Redis ZADD 命令用于将一个或多个 member 元素及其 score 值加入到有序集 key 当中。
		// https://www.runoob.com/redis/sorted-sets-zadd.html
		// https://redis.com.cn/commands/zadd.html
		pipeliner.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: postID,
		})
		// 帖子分数
		pipeliner.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: postID,
		})
		// Exec is to send all the commands buffered in the pipeline to the redis-server.
		_, err := pipeliner.Exec(ctx)
		return err

	})
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// 1. 判断投票限制
	// Redis中取帖子发布时间
	postTime := client.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2. 更新帖子的分数
	// 2.1 先查当前用户给当前帖子的投票记录
	ov := client.ZScore(ctx, getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	// 3. 记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  value,  // 赞成票还是反对票
			Member: userID, // 用户
		})
	}
	_, err := pipeline.Exec(ctx)
	return err
}
