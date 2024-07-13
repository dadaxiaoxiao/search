package repository

import (
	"context"
	"github.com/dadaxiaoxiao/search/internal/repository/dao"
)

type anyRepository struct {
	dao dao.AnyDAO
}

func NewAnyRepository(dao dao.AnyDAO) AnyRepository {
	return &anyRepository{dao: dao}
}

func (a *anyRepository) Input(ctx context.Context, index string, docId string, data string) error {
	return a.dao.Input(ctx, index, docId, data)
}
