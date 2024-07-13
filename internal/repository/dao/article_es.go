package dao

import (
	"context"
	"encoding/json"
	"github.com/ecodeclub/ekit/slice"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

const ArticleIndexName = "article_index"

// Article 文章实体
type Article struct {
	Id      int64    `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Status  int32    `json:"status"`
	Tags    []string `json:"tags"`
}

type ArticleElasticDAO struct {
	client *elastic.Client
}

func NewArticleElasticDAO(client *elastic.Client) ArticleDAO {
	return &ArticleElasticDAO{
		client: client,
	}
}

func (a *ArticleElasticDAO) InputArticle(ctx context.Context, article Article) error {
	_, err := a.client.Index().Index(ArticleIndexName).
		Id(strconv.FormatInt(article.Id, 10)).
		BodyJson(article).
		Do(ctx)
	return err
}

func (a *ArticleElasticDAO) Search(ctx context.Context, tagArtIds []int64, keywords []string) ([]Article, error) {
	queryString := strings.Join(keywords, " ")
	// 查询条件:
	// 文章，标题或者内容任何一个匹配上
	// 并且状态 status 必须是已发表的状态

	// 构建多个查询条件
	// status 使用精确擦查询
	statusOpts := elastic.NewTermQuery("status", 2)
	titleOpts := elastic.NewMatchQuery("title", queryString)
	content := elastic.NewMatchQuery("content", queryString)

	// 标签命中
	tagArtIdAnys := slice.Map(tagArtIds, func(idx int, src int64) any {
		return src
	})

	// or bool query   实现 （title or content）
	orOpts := elastic.NewBoolQuery().Should(titleOpts, content)
	if len(tagArtIds) > 0 {
		// 多个命中
		tag := elastic.NewTermsQuery("id", tagArtIdAnys...).
			Boost(2.0)
		orOpts = orOpts.Should(tag)
	}

	// and bool query  实现 status and (title or content)
	andOpts := elastic.NewBoolQuery().Must(statusOpts, orOpts)

	resp, err := a.client.Search(ArticleIndexName).Query(andOpts).Do(ctx)
	res := make([]Article, 0, resp.Hits.TotalHits.Value)
	for _, hit := range resp.Hits.Hits {
		var art Article
		err = json.Unmarshal(hit.Source, &art)
		if err != nil {
			return nil, err
		}
		res = append(res, art)
	}
	return res, nil
}
