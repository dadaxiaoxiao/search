package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"github.com/dadaxiaoxiao/go-pkg/saramax"
	"github.com/dadaxiaoxiao/search/internal/domain"
	"github.com/dadaxiaoxiao/search/internal/service"
	"time"
)

const topicSyncArticle = "sync_article_event"

type SyncArticleEvent struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Status  int32  `json:"status"`
	Content string `json:"content"`
}

type SyncArticleConsumer struct {
	syncSvc service.SyncService
	client  sarama.Client
	l       accesslog.Logger
}

func NewSyncArticleConsumer(client sarama.Client,
	svc service.SyncService,
	l accesslog.Logger) *SyncArticleConsumer {
	return &SyncArticleConsumer{
		client:  client,
		syncSvc: svc,
		l:       l,
	}
}

func (s *SyncArticleConsumer) Start() error {
	// 创建消费者组
	cg, err := sarama.NewConsumerGroupFromClient("search_sync_user", s.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicSyncArticle},
			saramax.NewHandler[SyncArticleEvent](s.l, s.Consume))
		if err != nil {
			s.l.Error("退出了消费循环异常", accesslog.Error(err))
		}
	}()
	return err
}

func (s *SyncArticleConsumer) Consume(msg *sarama.ConsumerMessage, evt SyncArticleEvent) error {
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return s.syncSvc.InputArticle(ctx, s.toDomain(evt))
}

func (u *SyncArticleConsumer) toDomain(evt SyncArticleEvent) domain.Article {
	return domain.Article{
		Id:      evt.Id,
		Title:   evt.Title,
		Status:  evt.Status,
		Content: evt.Content,
	}
}
