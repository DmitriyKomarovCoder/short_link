package grpc

import (
	"context"
	"errors"
	"github.com/DmitriyKomarovCoder/short_link/internal/models"
	pb "github.com/DmitriyKomarovCoder/short_link/internal/pkg/shortLink/delivery/grpc/gen"
	"github.com/DmitriyKomarovCoder/short_link/internal/pkg/shortLink/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
)

type LinkGrpcServer struct {
	pb.UnimplementedShortLinkServer
	usecase usecase.Usecase
}

func NewLinkGrpcServer(usecase usecase.Usecase) *LinkGrpcServer {
	return &LinkGrpcServer{
		usecase: usecase,
	}
}

func (g *LinkGrpcServer) CreateLink(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	longUrl := req.Url

	parsedURL, err := url.ParseRequestURI(longUrl)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, status.Error(codes.InvalidArgument, "bad url")
	}

	link, err := g.usecase.CreateLink(parsedURL.Host)
	if err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	return &pb.Response{Url: link}, nil
}

func (g *LinkGrpcServer) GetLink(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	shortURL := req.Url

	link, err := g.usecase.GetUrl(shortURL)

	var noSuchLinkErr *models.NoSuchLink
	if err != nil {
		if errors.As(err, &noSuchLinkErr) {
			return nil, status.Error(codes.NotFound, "such link does not exist or it has expired")
		} else {
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &pb.Response{Url: link}, nil
}
