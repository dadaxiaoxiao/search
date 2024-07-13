package ioc

import (
	"fmt"
	"github.com/dadaxiaoxiao/search/internal/repository/dao"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"time"
)

func InitESClient() *elastic.Client {
	type Config struct {
		Url   string `yaml:"url"`
		Sniff bool   `yaml:"sniff"`
	}
	var cfg Config
	err := viper.UnmarshalKey("elasticsearch", &cfg)
	if err != nil {
		panic(fmt.Errorf("读取 Elasticsearch 配置失败 %w", err))
	}
	const timeout = 100 * time.Second
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(cfg.Url),
		elastic.SetSniff(cfg.Sniff),
		// 健康检查超时启动
		elastic.SetHealthcheckTimeoutStartup(timeout),
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	err = dao.InitEs(client)
	if err != nil {
		panic(err)
	}
	return client
}
