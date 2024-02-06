package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	con "github.com/DmitriyKomarovCoder/short_link/common/config"
	"github.com/DmitriyKomarovCoder/short_link/common/logger"
	"github.com/DmitriyKomarovCoder/short_link/internal/pkg/linkGenerator"
	grpcLink "github.com/DmitriyKomarovCoder/short_link/internal/pkg/shortLink/delivery/grpc"
	pb "github.com/DmitriyKomarovCoder/short_link/internal/pkg/shortLink/delivery/grpc/gen"
	handler "github.com/DmitriyKomarovCoder/short_link/internal/pkg/shortLink/delivery/http"
	"github.com/DmitriyKomarovCoder/short_link/internal/pkg/shortLink/repository"
	"github.com/DmitriyKomarovCoder/short_link/internal/pkg/shortLink/usecase"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var usePostgreSQL bool

func init() {
	flag.BoolVar(&usePostgreSQL, "usePostgresSQL", true, "Set to true if you want to use PostgreSQL, otherwise Redis will be used.")
	flag.Parse()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := con.LoadConfig()
	if err != nil {
		log.Fatalf("Error config loading: %v", err)
	}

	logger, err := logger.NewLogger(config.LogFile.Path)
	if err != nil {
		log.Fatalf("Error logger loading: %v", err)
	}

	var repo repository.Repository
	if usePostgreSQL {
		logger.Info("Use PostgresSQL")
		repo = repository.NewPostgreSQLRepository(ctx, config.Postgres.Host,
			config.Postgres.Port,
			config.Postgres.User,
			config.Postgres.Password,
			config.Postgres.Name,
			*logger)
	} else {
		logger.Info("Use Redis")
		repo = repository.NewRedisRepository(config.Redis.Address, config.Redis.DB, *logger)
	}

	if err := repo.Connect(); err != nil {
		logger.Errorf("Error Initializing connect: %v", err)
		return
	}

	logger.Info("Db Connect successfully")

	defer func(repo repository.Repository) {
		err := repo.Close()
		if err != nil {
			logger.Errorf("Error closing connection: %v", err)
		}
	}(repo)

	linkGen := linkGenerator.NewLinkHash(config.UrlGenerate.Alphabet, config.UrlGenerate.Length)

	useCase := usecase.NewUsecase(repo, *logger, linkGen)
	urlHandler := handler.NewHandler(useCase, *logger)

	router := gin.Default()
	api := router.Group("/api")

	api.POST("/save", urlHandler.CreateLink)
	api.GET("/url/:url", urlHandler.GetLink)

	server := &http.Server{
		Addr:         config.Server.Host + config.Server.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	if usePostgreSQL {
		go func() {
			ticker := time.NewTicker(24 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					err := repo.Clear()
					if err != nil {
						logger.Println("Error clearing records:", err)
					}
				}
			}
		}()
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Error starting server: %v", err)
		}
	}()

	go func() {
		grpcHandler := grpcLink.NewLinkGrpcServer(*useCase)
		grpcServer := grpc.NewServer()
		pb.RegisterShortLinkServer(grpcServer, grpcHandler)

		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Server.Host, config.Server.GrpcPort))
		if err != nil {
			logger.Fatalf("Failed to listen: %v", err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve: %v", err)
		}
		logger.Info("Grpc server start in port: ", config.Server.GrpcPort)
	}()

	<-stop

	logger.Info("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown error: %v", err)
	}

	logger.Info("Server gracefully stopped")
}
