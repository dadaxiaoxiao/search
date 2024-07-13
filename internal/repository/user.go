package repository

import (
	"context"
	"github.com/dadaxiaoxiao/search/internal/domain"
	"github.com/dadaxiaoxiao/search/internal/repository/dao"
	"github.com/ecodeclub/ekit/slice"
)

type userRepository struct {
	dao dao.UserDAO
}

func NewUserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{dao: dao}
}

func (u *userRepository) InputUser(ctx context.Context, user domain.User) error {
	return u.dao.InputUser(ctx, u.toDao(user))
}

func (u *userRepository) SearchUser(ctx context.Context, keywords []string) ([]domain.User, error) {
	users, err := u.dao.Search(ctx, keywords)
	if err != nil {
		return nil, err
	}
	return slice.Map(users, func(idx int, src dao.User) domain.User {
		return u.toDomain(src)
	}), nil
}

func (u *userRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Nickname: user.Nickname,
		Phone:    user.Phone,
	}
}

func (u *userRepository) toDao(user domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		Email:    user.Email,
		Nickname: user.Nickname,
		Phone:    user.Phone,
	}
}
