package redis

import (
	"bluebell/models"
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	start := (page - 1) * size
	end := start + size - 1
	// 3. ZREVRANGE 按分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(ctx, key, start, end).Result()
}

// GetPostIDsInOrder 查询ID列表
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从Redis获取ID
	// http://redis.cn/commands/zrange.html
	/*
		➜ redis-cli
		127.0.0.1:6379> keys bluebell*
		1) "bluebell:post:voted:12807636613337088"
		2) "bluebell:post:score"
		3) "bluebell:post:time"
		127.0.0.1:6379> zrange bluebell:post:time 0 2
		1) "12807636613337088"
		127.0.0.1:6379> zrange bluebell:post:time 0 2 withscores
		1) "12807636613337088"
		2) "1690445578"
		127.0.0.1:6379> zrevrange bluebell:post:time 0 2 withscores
		1) "12807636613337088"
		2) "1690445578"
		127.0.0.1:6379>
	*/
	//ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	//defer cancel()
	// 1. 根据用户请求中携带的order参数确定要查询的 Redis Key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2. 确定查询的索引起始点
	//start := (p.Page - 1) * p.Size
	//end := start + p.Size - 1
	//// 3. ZREVRANGE 按分数从大到小的顺序查询指定数量的元素
	//return client.ZRevRange(ctx, key, start, end).Result()
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPrefix + id)
	//	// 查找Key中分数是1的元素的数量 -> 统计每篇帖子的赞成票的数量
	//	v := client.ZCount(ctx, key, "1", "1").Val()
	//	data = append(data, v)
	//}
	// 使用pipeline一次发送多条命令，减少RTT
	pipeline := client.TxPipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 使用 ZINTERSTORE 把分区的帖子SET与帖子分数的ZSET生成一个新的ZSET
	// 针对新的ZSET 按之前的逻辑提取数据
	// http://redis.cn/commands/zinterstore.html
	// 社区的Key
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))
	// 利用缓存Key减少 ZINTERSTORE 执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(ctx, key).Val() < 1 {
		//	不存在，需要计算
		pipeline := client.Pipeline()
		// Define the options for the ZInterStore operation
		options := redis.ZStore{
			Keys:      []string{cKey, orderKey},
			Aggregate: "MAX",
		}
		pipeline.ZInterStore(ctx, key, &options)
		pipeline.Expire(ctx, key, 60*time.Second)
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	// 存在则直接根据Key查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}
