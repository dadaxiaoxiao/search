package dao

import "context"

// AnyDAO 通用的数据访问对象
//
// 只是提供模糊的上传功能，不提供查询
type AnyDAO interface {
	// Input
	// index 索引
	// docId 文档id
	// data json 结构体
	Input(ctx context.Context, index, docId, data string) error
}

type UserDAO interface {
	InputUser(ctx context.Context, user User) error
	Search(ctx context.Context, keywords []string) ([]User, error)
}

type ArticleDAO interface {
	InputArticle(ctx context.Context, article Article) error
	Search(ctx context.Context, tagArtIds []int64, keywords []string) ([]Article, error)
}

type TagDAO interface {
	Search(ctx context.Context, uid int64, biz string, keywords []string) ([]BizTags, error)
}
