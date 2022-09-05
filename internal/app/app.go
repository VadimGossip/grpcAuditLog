package app

import (
	"context"
	"github.com/VadimGossip/grpcAuditLog/internal/config"
	"github.com/VadimGossip/grpcAuditLog/internal/repository"
	"github.com/VadimGossip/grpcAuditLog/internal/server"
	"github.com/VadimGossip/grpcAuditLog/internal/service"
	"github.com/VadimGossip/grpcAuditLog/pkg/database"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func Run(configDir string) {
	cfg, err := config.Init(configDir)
	if err != nil {
		logrus.Fatalf("Config initialization error %s", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, err := database.NewMongoConnection(ctx, cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.URI)
	if err != nil {
		logrus.Fatalf("Mongo connection error %s", err)
	}
	db := dbClient.Database(cfg.Mongo.Database)

	auditRepo := repository.NewAudit(db)
	auditService := service.NewAudit(auditRepo)
	auditSrv := server.NewAuditServer(auditService)
	srv := server.New(auditSrv)

	logrus.Info("Audit Server for fin manager service started")
	if err := srv.ListenAndServe(cfg.Server.Port); err != nil {
		logrus.Fatalf("error occured while running audit server for fin manager: %s", err.Error())
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("Audit Server for fin manager service stopped")

	//db, err := database.NewPostgresConnection(cfg.Postgres)
	//if err != nil {
	//	logrus.Fatalf("Postgres connection error %s", err)
	//}
	//
	//usersRepo := psql.NewUsers(db)
	//tokensRepo := psql.NewTokens(db)
	//hasher := hash.NewSHA1Hasher(cfg.Auth.Salt)
	//usersService := service.NewUsers(usersRepo, tokensRepo, hasher, []byte(cfg.Auth.Secret), cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL)
	//
	//docsRepo := psql.NewDocs(db)
	//cache := simpleCache.NewCache()
	//docsService := service.NewBooks(docsRepo, cache)
	//
	//handler := rest.NewHandler(usersService, docsService)
	//server := http.NewServer()
	//
	//go func() {
	//	if err := server.Run(cfg.Server, handler.InitRoutes()); err != nil {
	//		logrus.Fatalf("error occured while running rest server: %s", err.Error())
	//	}
	//}()
	//
	//logrus.Info("Http Server for fin manager service started")
	//
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	//<-quit
	//
	//logrus.Info("Http Server for fin manager service stopped")
	//
	//if err := db.Close(); err != nil {
	//	logrus.Errorf("Error occured on postgres connection close: %s", err.Error())
	//}
	//
	//if err := server.Shutdown(context.Background()); err != nil {
	//	logrus.Errorf("Error occured on http server for fin manager service shutting down: %s", err.Error())
	//}
}
