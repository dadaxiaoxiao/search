package repository

import (
	"context"
	"github.com/dadaxiaoxiao/search/internal/domain"
)

// AnyRepository 通用的仓储层
//
// 只是提供模糊的上传功能，不提供查询
type AnyRepository interface {
	Input(ctx context.Context, index string, docId string, data string) error
}

type UserRepository interface {
	InputUser(ctx context.Context, user domain.User) error
	SearchUser(ctx context.Context, keywords []string) ([]domain.User, error)
}

type ArticleRepository interface {
	InputArticle(ctx context.Context, article domain.Article) error
	SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error)
}
