package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

const UserIndexName = "User_index"

// User 用户实体
type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

type UserElasticDAO struct {
	client *elastic.Client
}

func NewUserElasticDAO(client *elastic.Client) UserDAO {
	return &UserElasticDAO{client: client}
}

func (u *UserElasticDAO) InputUser(ctx context.Context, user User) error {
	_, err := u.client.Index().
		Index(UserIndexName).
		// 这里使用update set 语义
		Id(strconv.FormatInt(user.Id, 10)).
		BodyJson(user).Do(ctx)
	return err
}

func (u *UserElasticDAO) Search(ctx context.Context, keywords []string) ([]User, error) {
	// 前面已经预处理了输入
	queryString := strings.Join(keywords, " ")
	// 昵称命中就可以的
	resp, err := u.client.Search(UserIndexName).
		Query(elastic.NewMatchQuery("nickname", queryString)).Do(ctx)

	if err != nil {
		return nil, err
	}
	res := make([]User, 0, resp.Hits.TotalHits.Value)
	for _, hits := range resp.Hits.Hits {
		var u User
		err = json.Unmarshal(hits.Source, &u)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}
