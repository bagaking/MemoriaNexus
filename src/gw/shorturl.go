package gw

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ShortURLLen   = 8
	ShortURLTTL   = 30 * 365 * 24 * time.Hour // 约30年
	ShortURLRoute = "/utils/s/"
)

var shortURLCache sync.Map

type ShortURLEntry struct {
	OriginalURL string
	ExpiresAt   time.Time
}

// SetupShortURLRoutes 设置短网址相关路由
func SetupShortURLRoutes(router *gin.Engine) {
	router.GET(ShortURLRoute+":shortURL", handleShortURL)
	router.POST(ShortURLRoute, createShortURL)
}

// 处理短网址重定向
func handleShortURL(c *gin.Context) {
	shortURL := c.Param("shortURL")
	entry, ok := shortURLCache.Load(shortURL)
	if !ok {
		c.String(http.StatusNotFound, "Short URL not found")
		return
	}

	urlEntry := entry.(ShortURLEntry)
	if time.Now().After(urlEntry.ExpiresAt) {
		shortURLCache.Delete(shortURL)
		c.String(http.StatusNotFound, "Short URL expired")
		return
	}

	c.Redirect(http.StatusMovedPermanently, urlEntry.OriginalURL)
}

// GetFullURLFromGinCtx 从 Gin 上下文中获取完整的 URL, 处理 X-Original-URI, X-Forwarded-Proto 和 X-Forwarded-Host 头
// 网关代理场景下，可能进行 redirect 操作，此时可以通过设置 X-Original-URI 头来指定原始 URI
func GetFullURLFromGinCtx(c *gin.Context) string {
	// 首先尝试从 X-Original-URI 头获取原始 URI
	originalURI := c.GetHeader("X-Original-URI")
	if originalURI == "" {
		originalURI = c.Request.URL.Path
	}

	// 确定 scheme
	scheme := "https"
	if c.Request.TLS == nil {
		scheme = "http"
	}
	// 如果 X-Forwarded-Proto 头存在，优先使用它
	if forwardedProto := c.GetHeader("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = forwardedProto
	}

	// 确定 host
	host := c.Request.Host
	// 如果 X-Forwarded-Host 头存在，优先使用它
	if forwardedHost := c.GetHeader("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	// 构建并返回完整的 URL
	return fmt.Sprintf("%s://%s%s", scheme, host, originalURI)
}

// 创建短网址
// e.g. curl -X POST https://utils.kenv.tech/s/ -H "Content-Type: application/json" -d '{"url": "xxxxxxxxxx"}'
func createShortURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortURL := generateShortURL(req.URL)
	expiresAt := time.Now().Add(ShortURLTTL)

	shortURLCache.Store(shortURL, ShortURLEntry{
		OriginalURL: req.URL,
		ExpiresAt:   expiresAt,
	})

	fullURL := GetFullURLFromGinCtx(c)

	// 解析 URL
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析 URL"})
		return
	}

	parsedURL.RawPath = path.Join(parsedURL.RawPath, shortURL)

	c.JSON(http.StatusOK, gin.H{
		"full_url":   fullURL,
		"short_url":  parsedURL.RawPath,
		"expires_at": expiresAt,
	})
}

// 生成短网址
func generateShortURL(url string) string {
	hash := md5.Sum([]byte(url))
	return base64.URLEncoding.EncodeToString(hash[:])[:ShortURLLen]
}
