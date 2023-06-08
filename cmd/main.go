package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-app/internal/dao/mysql"
	"web-app/internal/dao/redis"
	"web-app/internal/logger"
	"web-app/internal/routes"
	"web-app/settings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Go Web 开发通用的脚手架模板

func main() {
	// 1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed err: %v\n", err)
		panic(err)
	}

	// 2.初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failed err: %v\n", err)
		panic(err)
	}
	defer func() {
		zap.L().Sync()
		logger.SugarLogger.Sync()
	}()
	zap.L().Debug("log init success ... ")

	// 3.初始化mysql
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed err: %v\n", err)
		panic(err)
	}

	// 4.初始化redis
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed err: %v\n", err)
		panic(err)
	}
	defer redis.RDB.Close()

	// 5.注册路由
	r := routes.Setup()

	// 6.优雅关机
	Run(r, viper.GetString("app.name"))
}

func Run(r *gin.Engine, srvName string) {
	srv := &http.Server{
		Addr:    viper.GetString("app.port"),
		Handler: r,
	}

	// 保证优雅启停
	go func() {
		logger.SugarLogger.Debugf("%s running in %s \n", srvName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.SugarLogger.Errorf("listen: %s\n", err)
		}
	}()

	// 相当于监听一下 kill -2 和 kill -9
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT (Ctrl + C)
	// kill -9 is syscall. SIGKILL
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.SugarLogger.Debugf("Shutdown %s ...\n", srvName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.SugarLogger.Errorf("%s Shutdown err: %v\n", srvName, err)
	}
	// catching ctx.Done(). timeout of 2 seconds.
	select {
	case <-ctx.Done():
		zap.L().Debug("timeout of 2 seconds.")
	}
	logger.SugarLogger.Debugf("%s exiting", srvName)
}
