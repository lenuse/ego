package http

import (
	"github.com/gin-gonic/gin"
	"errors"
	"github.com/ebar-go/ego/http/handler"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/constant"
	"github.com/ebar-go/ego/log"
	"path"
	"github.com/sirupsen/logrus"
)


// Server Web服务管理器
type Server struct {
	LogPath string
	AppDebug bool
	SystemName string
	Address string
	Router *gin.Engine
	initialize bool
	NotFoundHandler func(ctx *gin.Context)
	Recover func(ctx *gin.Context)
}

// Init 服务初始化
func (server *Server)Init() error {
	if server.initialize {
		return errors.New("请勿重复初始化Http Server")
	}

	server.Router = gin.Default()

	// 请求日志
	server.Router.Use(middleware.RequestLog)

	if server.NotFoundHandler == nil {
		server.NotFoundHandler = handler.NotFoundHandler
	}

	if server.Recover == nil {
		server.Recover = handler.Recover
	}
	server.Router.Use(server.Recover)

	// 404
	server.Router.NoRoute(server.NotFoundHandler)
	server.Router.NoMethod(server.NotFoundHandler)

	server.initLogger()

	server.initialize = true
	return nil
}


func (server *Server) initLogger() error {
	// 初始化日志目录
	if server.LogPath == "" {
		server.LogPath = constant.DefaultLogPath
	}

	if server.SystemName == "" {
		server.SystemName = constant.DefaultSystemName
	}


	appPath := path.Join(server.LogPath, server.SystemName, constant.AppLogPrefix + server.SystemName + constant.LogSuffix)
	systemPath := path.Join(server.LogPath, server.SystemName, constant.SystemLogPrefix + server.SystemName + constant.LogSuffix)
	requestPath := path.Join(server.LogPath, server.SystemName, constant.RequestLogPrefix + server.SystemName + constant.LogSuffix)

	log.AppLogger = log.NewFileLogger(appPath)

	if server.AppDebug {
		log.AppLogger.SetLevel(logrus.DebugLevel)
	}
	log.SystemLogger = log.NewFileLogger(systemPath)
	log.RequestLogger = log.NewFileLogger(requestPath)

	return nil
}

// Start 启动服务
func (server *Server) Start() error {
	if !server.initialize {
		return errors.New("请先初始化服务")
	}

	return server.Router.Run(server.Address)
}
