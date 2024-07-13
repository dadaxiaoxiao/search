package grpc

import (
	"context"
	searchv1 "github.com/dadaxiaoxiao/api-repository/api/proto/gen/search/v1"
	"github.com/dadaxiaoxiao/search/internal/domain"
	"github.com/dadaxiaoxiao/search/internal/service"
	"google.golang.org/grpc"
)

type SyncServiceServer struct {
	searchv1.UnimplementedSyncServiceServer
	svc service.SyncService
}

func NewSyncServiceServer(svc service.SyncService) *SyncServiceServer {
	return &SyncServiceServer{
		svc: svc,
	}
}

func (s *SyncServiceServer) Register(server *grpc.Server) {
	searchv1.RegisterSyncServiceServer(server, s)
}

func (s *SyncServiceServer) InputAny(ctx context.Context, req *searchv1.InputAnyRequest) (*searchv1.InputAnyResponse, error) {
	err := s.svc.InputAny(ctx, req.GetIndexName(), req.GetDocId(), req.GetData())
	return &searchv1.InputAnyResponse{}, err
}

func (s *SyncServiceServer) InputUser(ctx context.Context, req *searchv1.InputUserRequest) (*searchv1.InputUserResponse, error) {
	err := s.svc.InputUser(ctx, s.toDomainUser(req.GetUser()))
	return &searchv1.InputUserResponse{}, err
}
func (s *SyncServiceServer) InputArticle(ctx context.Context, req *searchv1.InputArticleRequest) (*searchv1.InputArticleResponse, error) {
	err := s.svc.InputArticle(ctx, s.toDomainArticle(req.GetArticle()))
	return &searchv1.InputArticleResponse{}, err
}

func (s *SyncServiceServer) toDomainUser(vuser *searchv1.User) domain.User {
	return domain.User{
		Id:       vuser.GetId(),
		Email:    vuser.GetEmail(),
		Nickname: vuser.GetNickname(),
		Phone:    vuser.GetPhone(),
	}
}

func (s *SyncServiceServer) toDomainArticle(varticle *searchv1.Article) domain.Article {
	return domain.Article{
		Id:     varticle.GetId(),
		Title:  varticle.GetTitle(),
		Status: varticle.GetStatus(),
	}
}
