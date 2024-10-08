// File: cmd/memorial_nexus.go

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"

	"github.com/bagaking/memorianexus/src/gw"

	"github.com/adjust/redismq"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/khgame/ranger_iam/pkg/authcli"

	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/khicago/irr"
	"github.com/sirupsen/logrus"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/internal/utils/cache"
)

const APIGroup = "/api/v1"

// dsn for dev
// todo: using env config
func dsn() string {
	// 使用docker-compose环境变量来设置数据库DSN
	username := "user"
	password := "password"

	// todo: using config
	host := "localhost"
	switch utils.Env() {
	case utils.RuntimeENVDev:
		host = "mysql" // run in docker compose
	case utils.RuntimeENVStaging:
		host = "0.0.0.0" // run at remote service (staging cluster)
	case utils.RuntimeENVProd:
		host = "0.0.0.0" // run at remote service (product cluster)
	case utils.RuntimeENVLocal:
		fallthrough
	default:
	}

	port := "3306"
	dbname := "memorianexus"
	charset := "utf8mb4"
	loc := "Local"
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s", username, password, host, port, dbname, charset, loc)
}

func redisHost() string {
	host := "localhost"
	switch utils.Env() {
	case utils.RuntimeENVDev:
		host = "redis"
	case utils.RuntimeENVStaging:
		host = "0.0.0.0"
	case utils.RuntimeENVProd:
		host = "0.0.0.0"
	case utils.RuntimeENVLocal:
		fallthrough
	default:
	}
	return host
}

func redisDSN() string {
	port := "6379"
	return redisHost() + ":" + port
}

// 初始化 RedisMQ
func mustInitRedisMQ(redisAddr string) *redismq.Queue {
	queue := redismq.CreateQueue(redisAddr, "6379", "", 0, "memnexus")
	if queue == nil {
		wlog.Common("memorial_nexus", "mustInitRedisMQ").Fatal("failed to create RedisMQ queue")
	}
	return queue
}

func main() {
	// 配置一个lumberjack.Logger
	logRoller := &lumberjack.Logger{
		Filename:   "./logs/memnexus.log", // 日志文件的位置
		MaxSize:    10,                    // 日志文件的最大大小（MB）
		MaxBackups: 31,                    // 保存的旧日志文件最大个数
		MaxAge:     31,                    // 保存的旧日志文件的最大天数
		Compress:   true,                  // 是否压缩归档的日志文件
	}

	defer func() {
		if err := logRoller.Close(); err != nil {
			fmt.Println("Failed to close log", err)
		}
	}()

	mustInitLogger(logRoller)
	startLogger := wlog.Common("memnexus")

	// 初始化数据库连接
	db := mustInitDB()

	// 初始化缓存
	cache.Init(redisDSN())
	redisMQInst := mustInitRedisMQ(redisHost())

	// 初始化HTTP路由
	router := gin.Default()
	router.Use(
		ginRecoveryWithLog(),
		corsMiddleware(),
	)

	// todo: 这些值应该从配置中安全获取，现在 MVP 一下
	iamCli := authcli.New("my_secret_key", "http://0.0.0.0:8090/")

	model.MustInit(context.TODO(), db, redisMQInst)

	// todo: 挪到单独的服务里 ?
	gw.RegRouter(router, db, iamCli, APIGroup, "/var/static/memnexus")

	startLogger.Trace("memnexus initialed")
	startLogger.Debug("memnexus initialed")
	startLogger.Info("memnexus initialed")
	// startLogger.Error("memnexus initialed")

	// 开启HTTP服务
	if err := router.Run(":8080"); err != nil {
		startLogger.WithError(err).Infof("gin exit")
	}
}

func mustInitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn()), &gorm.Config{})
	if err != nil {
		wlog.Common("memorial_nexus", "mustInitDB").Fatal("failed to connect database:", err)
	}
	return db
}

func mustInitLogger(fileLogger io.Writer) {
	// 创建error日志文件
	errorLogRoller := &lumberjack.Logger{
		Filename:   "./logs/memnexus.error.log",
		MaxSize:    10,
		MaxBackups: 31,
		MaxAge:     31,
		Compress:   true,
	}

	if utils.Env() == utils.RuntimeENVProd {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: true,
		})
		logrus.SetLevel(logrus.TraceLevel)
	}

	// 创建一个自定义的Writer,将error及以上级别的日志写入error日志文件
	errorWriter := &logrusErrorWriter{errorLogRoller}

	// 使用MultiWriter将日志同时输出到标准输出、普通日志文件和error日志文件
	multiLogger := io.MultiWriter(os.Stdout, fileLogger, errorWriter)
	logrus.SetOutput(multiLogger)

	gin.DefaultWriter = logrus.StandardLogger().Out

	wlog.SetEntryGetter(func(ctx context.Context) *logrus.Entry {
		return logrus.WithContext(ctx)
	})
}

// logrusErrorWriter 是一个自定义的io.Writer,只写入error及以上级别的日志
type logrusErrorWriter struct {
	io.Writer
}

func (w *logrusErrorWriter) Write(p []byte) (n int, err error) {
	entry := &logrus.Entry{}
	err = jsoniter.Unmarshal(p, entry)
	if err != nil {
		return 0, err
	}

	if entry.Level <= logrus.ErrorLevel {
		return w.Writer.Write(p)
	}
	return len(p), nil
}

// ginRecoveryWithLog 返回一个中间件，当程序发生 panic 时记录错误日志，并返回 HTTP 500 错误。
func ginRecoveryWithLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var er irr.IRR
				if e, ok := err.(error); ok {
					er = irr.TrackSkip(1, e, "recover from panic!!")
				} else {
					er = irr.TrackSkip(1, irr.Error("%v", err), "recover from panic!!")
				}
				log := wlog.ByCtx(c, c.HandlerName())

				utils.GinHandleError(c, log, http.StatusInternalServerError, er, "Unexpected internal Server Error")
			}
		}()

		// 处理请求
		c.Next()
	}
}

// corsMiddleware 处理跨域请求
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
