package repository

import (
	"context"
	"github.com/dadaxiaoxiao/search/internal/domain"
	"github.com/dadaxiaoxiao/search/internal/repository/dao"
	"github.com/ecodeclub/ekit/slice"
)

type articleRepository struct {
	dao  dao.ArticleDAO
	tags dao.TagDAO
}

func NewArticleRepository(dao dao.ArticleDAO, tags dao.TagDAO) ArticleRepository {
	return &articleRepository{
		dao:  dao,
		tags: tags,
	}
}

func (a *articleRepository) InputArticle(ctx context.Context, article domain.Article) error {
	return a.dao.InputArticle(ctx, a.toDao(article))
}

func (a *articleRepository) SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error) {

	// 查询命中标签的tag数据
	tags, err := a.tags.Search(ctx, uid, "article", keywords)
	// 找到对应的文章id (bizId)
	bizIds := slice.Map(tags, func(idx int, src dao.BizTags) int64 {
		return src.BizId
	})

	arts, err := a.dao.Search(ctx, bizIds, keywords)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return a.toDomain(src)
	}), nil
}

func (a *articleRepository) toDao(art domain.Article) dao.Article {
	return dao.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  art.Status,
		Tags:    art.Tags,
	}
}

func (a *articleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  art.Status,
		Tags:    art.Tags,
	}
}
