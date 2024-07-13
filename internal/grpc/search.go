package grpc

import (
	"context"
	searchv1 "github.com/dadaxiaoxiao/api-repository/api/proto/gen/search/v1"
	"github.com/dadaxiaoxiao/search/internal/domain"
	"github.com/dadaxiaoxiao/search/internal/service"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

type SearchServiceServer struct {
	searchv1.UnimplementedSearchServiceServer
	svc service.SearchService
}

func NewSearchServiceServer(svc service.SearchService) *SearchServiceServer {
	return &SearchServiceServer{
		svc: svc,
	}
}

func (s *SearchServiceServer) Register(server *grpc.Server) {
	searchv1.RegisterSearchServiceServer(server, s)
}

func (s *SearchServiceServer) Search(ctx context.Context, req *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
	resp, err := s.svc.Search(ctx, req.Uid, req.Expression)
	if err != nil {
		return nil, err
	}
	return &searchv1.SearchResponse{
		User: &searchv1.UserResult{
			Users: slice.Map[domain.User, *searchv1.User](resp.Users, func(idx int, src domain.User) *searchv1.User {
				return &searchv1.User{
					Id:       src.Id,
					Nickname: src.Nickname,
					Email:    src.Email,
					Phone:    src.Phone,
				}
			}),
		},
		Article: &searchv1.ArticleResult{
			Articles: slice.Map[domain.Article, *searchv1.Article](resp.Articles, func(idx int, src domain.Article) *searchv1.Article {
				return &searchv1.Article{
					Id:      src.Id,
					Title:   src.Title,
					Status:  src.Status,
					Content: src.Content,
				}
			}),
		},
	}, err

}
