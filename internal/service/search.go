package service

import (
	"context"
	"github.com/dadaxiaoxiao/search/internal/domain"
	"github.com/dadaxiaoxiao/search/internal/repository"
	"golang.org/x/sync/errgroup"
	"strings"
)

// SearchService 查询服务
//
// 只是在Service 进行CQRS
type SearchService interface {
	Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error)
}

type searchService struct {
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
}

func NewSearchService(
	userRepo repository.UserRepository,
	articleRepo repository.ArticleRepository) SearchService {
	return &searchService{}
}

func (s *searchService) Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error) {
	// 对 expression 进行预处理
	// 统一使用的空格符来分割的
	keywords := strings.Split(expression, " ")

	// 因为设计是全站模糊搜索，它可以搜出任何我们支持的内容
	var eg errgroup.Group
	var res domain.SearchResult
	eg.Go(func() error {
		users, err := s.userRepo.SearchUser(ctx, keywords)
		res.Users = users
		return err
	})
	eg.Go(func() error {
		arts, err := s.articleRepo.SearchArticle(ctx, uid, keywords)
		res.Articles = arts
		return err
	})
	// 后续接入更多的数据，这里可以继续开eg.Go
	return res, eg.Wait()
}
