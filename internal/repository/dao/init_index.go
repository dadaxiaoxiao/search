package dao

import (
	"context"
	_ "embed"
	"github.com/olivere/elastic/v7"
	"golang.org/x/sync/errgroup"
	"time"
)

var (
	//go:embed user_index.json
	userIndex string
	//go:embed article_index.json
	articleIndex string
	//go:embed tag_index.json
	tagIndex string
)

func InitEs(client *elastic.Client) error {
	// 超时控制
	const timeout = time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var eg errgroup.Group
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, UserIndexName, userIndex)
	})
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, ArticleIndexName, articleIndex)
	})
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, "tag_index", tagIndex)
	})
	return eg.Wait()
}

// tryCreateIndex
//
// try to create es index
func tryCreateIndex(ctx context.Context,
	client *elastic.Client,
	idxName string, idxCfg string) error {

	// 判断索引是否已经创建
	ok, err := client.IndexExists(idxName).Do(ctx)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	// 创建索引
	_, err = client.CreateIndex(idxName).Do(ctx)
	return err
}
