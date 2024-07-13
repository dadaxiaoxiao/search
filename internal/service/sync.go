package service

import (
	"context"
	"github.com/dadaxiaoxiao/search/internal/domain"
	"github.com/dadaxiaoxiao/search/internal/repository"
)

// SyncService 同步服务
//
// 只是在Service 进行CQRS
type SyncService interface {
	InputAny(ctx context.Context, index, docId, data string) error
	InputUser(ctx context.Context, user domain.User) error
	InputArticle(ctx context.Context, article domain.Article) error
}

type syncService struct {
	anyRepo     repository.AnyRepository
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
}

func NewSyncService(anyRepo repository.AnyRepository,
	userRepo repository.UserRepository,
	articleRepo repository.ArticleRepository) SyncService {
	return &syncService{
		anyRepo:     anyRepo,
		userRepo:    userRepo,
		articleRepo: articleRepo,
	}
}

func (s *syncService) InputAny(ctx context.Context, index, docId, data string) error {
	return s.anyRepo.Input(ctx, index, docId, data)
}

func (s *syncService) InputUser(ctx context.Context, user domain.User) error {
	return s.userRepo.InputUser(ctx, user)
}

func (s *syncService) InputArticle(ctx context.Context, article domain.Article) error {
	return s.articleRepo.InputArticle(ctx, article)
}
