package models

import "time"

// 内存对齐概念
// 内存对齐是计算机系统中的一个概念，用于描述数据在内存中的存储方式。在内存中，数据是以字节为单位存储的，而内存对齐规定了数据在内存中存储的起始位置和对齐边界。
// 对齐边界是指数据在内存中的存储位置必须是某个特定字节的倍数。例如，对齐边界为4的情况下，数据的起始位置必须是4的倍数，即内存地址的末尾两位必须是00。
// 内存对齐的目的是为了提高数据的读取和访问效率。当数据按照对齐要求存储时，处理器可以更快地读取和处理数据，而不需要额外的操作来处理不对齐的数据。
// 在结构体中，不同的数据类型会有不同的对齐要求。通常，基本数据类型的对齐要求是其大小的最小值和处理器字长的较小值。结构体的对齐要求是其成员中对齐要求最大的数据类型的对齐要求。
// 了解和使用内存对齐概念可以帮助程序员优化内存使用和提高程序性能。

type Post struct {
	ID          int64     `json:"id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ApiPostDetail 帖子详情接口的结构体
type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"`
	*Post                               // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区信息
}
