package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"github.com/dadaxiaoxiao/go-pkg/saramax"
	"github.com/dadaxiaoxiao/search/internal/service"
	"time"
)

const topicSyncData = "sync_data_event"

// SyncAnyDataEvent 通用的同步数据
//
// 假如说用于同步 user
// IndexName = user_index
// DocID = "123"
// Data = {"id": 123, "email":xx, nickname: ""}
type SyncDataEvent struct {
	IndexName string `json:"indexName"`
	DocID     string `json:"docID"`
	Data      string `json:"data"`
}

type SyncDataEventConsumer struct {
	client sarama.Client
	svc    service.SyncService
	l      accesslog.Logger
}

func NewSyncAnyDataEventConsumer(client sarama.Client,
	svc service.SyncService,
	l accesslog.Logger) *SyncDataEventConsumer {
	return &SyncDataEventConsumer{
		client: client,
		svc:    svc,
		l:      l,
	}
}

func (s *SyncDataEventConsumer) Start() error {
	// 创建消费者组
	cg, err := sarama.NewConsumerGroupFromClient("search_sync_data", s.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicSyncData},
			saramax.NewHandler[SyncDataEvent](s.l, s.Consume))
		if err != nil {
			s.l.Error("退出了消费循环异常", accesslog.Error(err))
		}
	}()
	return err
}

func (s *SyncDataEventConsumer) Consume(msg *sarama.ConsumerMessage, evt SyncDataEvent) error {
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return s.svc.InputAny(ctx, evt.IndexName, evt.DocID, evt.Data)
}
