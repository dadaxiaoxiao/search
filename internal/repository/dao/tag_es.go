package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
)

type TagESDAO struct {
	client *elastic.Client
}

func NewTagESDAO(client *elastic.Client) TagDAO {
	return &TagESDAO{client: client}
}

func (t *TagESDAO) Search(ctx context.Context, uid int64, biz string, keywords []string) ([]BizTags, error) {
	query := elastic.NewBoolQuery().Must(
		// 条件1 精确查询uid
		elastic.NewTermQuery("uid", uid),
		// 条件2 精确查询biz
		elastic.NewTermQuery("biz", biz),
		// 条件3，关键字命中了标签
		elastic.NewTermsQueryFromStrings("tags", keywords...),
	)
	resp, err := t.client.Search().Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]BizTags, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var ele BizTags
		// 因为 文档当时是 json 字节流
		err = json.Unmarshal(hit.Source, &ele)
		if err != nil {
			return nil, err
		}
		// 把 bizTag 拿出来了
		res = append(res, ele)
	}
	return res, err
}

type BizTags struct {
	Uid   int64    `json:"uid"`
	Biz   string   `json:"biz"`
	BizId int64    `json:"biz_id"`
	Tags  []string `json:"tags"`
}
