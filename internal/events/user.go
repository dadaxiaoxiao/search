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

const topicSyncUser = "sync_user_event"

type SyncUserEvent struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
}

type SyncUserConsumer struct {
	syncSvc service.SyncService
	client  sarama.Client
	l       accesslog.Logger
}

func NewSyncUserConsumer(client sarama.Client,
	svc service.SyncService,
	l accesslog.Logger) *SyncUserConsumer {
	return &SyncUserConsumer{
		client:  client,
		syncSvc: svc,
		l:       l,
	}
}

func (s *SyncUserConsumer) Start() error {
	// 创建消费者组
	cg, err := sarama.NewConsumerGroupFromClient("search_sync_user", s.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicSyncUser},
			saramax.NewHandler[SyncUserEvent](s.l, s.Consume))
		if err != nil {
			s.l.Error("退出了消费循环异常", accesslog.Error(err))
		}
	}()
	return err
}

func (s *SyncUserConsumer) Consume(msg *sarama.ConsumerMessage, evt SyncUserEvent) error {
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return s.syncSvc.InputUser(ctx, s.toDomain(evt))
}

func (u *SyncUserConsumer) toDomain(evt SyncUserEvent) domain.User {
	return domain.User{
		Id:       evt.Id,
		Email:    evt.Email,
		Nickname: evt.Nickname,
		Phone:    evt.Phone,
	}
}
