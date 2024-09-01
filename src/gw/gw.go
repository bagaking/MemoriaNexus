package gw

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/bagaking/goulp/wlog"
	"github.com/sirupsen/logrus"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/khgame/ranger_iam/pkg/authcli"
	"gorm.io/gorm"

	"github.com/bagaking/memorianexus/doc"
)

const (
	IndexFile = "index.html"
)

func RegRouter(router *gin.Engine, db *gorm.DB, iamCli *authcli.Cli, APIGroup string, staticFilePath string) {
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(PerformanceMonitor())

	if staticFilePath != "" {
		SetupStaticFileServer(router, staticFilePath, APIGroup)
	}

	doc.SwaggerInfo.BasePath = APIGroup
	group := router.Group(APIGroup)
	RegisterCallbacks(group)
	RegisterRoutes(group, db, iamCli)
}

// SetupStaticFileServer 配置静态文件服务和前端路由处理
func SetupStaticFileServer(router gin.IRouter, staticDir string, apiGroup string) {
	fileSystem := http.Dir(staticDir)

	router.Use(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果是 API 路由，跳过静态文件处理
		if strings.HasPrefix(path, apiGroup) {
			c.Next()
			return
		}

		// 尝试提供静态文件
		if path != "/" {
			if file, err := fileSystem.Open(path); err == nil {
				defer file.Close()
				if info, err := file.Stat(); err == nil && !info.IsDir() {
					// 设置缓存控制头
					setCacheControlHeaders(c, path)
					http.ServeContent(c.Writer, c.Request, info.Name(), info.ModTime(), file)
					c.Abort()
					return
				}
			}
		}

		// 如果文件不存在或是根路径，返回 index.html
		wlog.ByCtx(c, "static").Infof("Serving index.html for path: %s", path)
		setCacheControlHeaders(c, IndexFile)
		c.File(filepath.Join(staticDir, IndexFile))
		c.Abort()
	})
}

// setCacheControlHeaders 设置缓存控制头
func setCacheControlHeaders(c *gin.Context, path string) {
	// 对于 HTML 文件，不使用缓存
	if strings.HasSuffix(path, ".html") || path == IndexFile {
		c.Header("Cache-Control", "no-store, must-revalidate")
		return
	}

	// 对于根目录下的其他文件，缓存 3 分钟
	if !strings.Contains(path, "/") {
		c.Header("Cache-Control", "public, max-age=180")
		return
	}

	// 对于其他静态资源，使用较长的缓存时间
	ext := filepath.Ext(path)
	switch ext {
	case ".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".svg":
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	default:
		c.Header("Cache-Control", "public, max-age=86400")
	}
}

func PerformanceMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(start)

		// 获取响应状态
		status := c.Writer.Status()

		// 获取客户端 IP
		clientIP := c.ClientIP()

		log := wlog.ByCtx(c, "static").WithFields(logrus.Fields{
			"path":      path,
			"method":    method,
			"status":    status,
			"client_ip": clientIP,
			"duration":  duration,
			"body_size": c.Writer.Size(),
		})
		// Log detailed performance information using wlog
		log.Info("Request Performance")

		// Log a warning if the processing time exceeds a threshold
		if duration > time.Second*2 { // e.g., 2 seconds
			log.Warn("Slow Request")
		}
	}
}
