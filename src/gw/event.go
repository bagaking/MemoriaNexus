package gw

import (
	"net/http"
	"time"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/internal/utils/cache"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type (
	Header struct {
		EventID   string `json:"event_id"`
		EventType string `json:"event_type"`
	}

	Callback[T any] struct {
		// challenge validate
		Type      string `json:"type,omitempty"`
		Challenge string `json:"challenge,omitempty"`

		Schema string `json:"schema"`
		Header Header `json:"header"`
		Event  T      `json:"event"`
	}

	ChallengeResp struct {
		Code      int    `json:"code,omitempty"`
		Msg       string `json:"msg,omitempty"`
		Challenge string `json:"challenge,omitempty"`
	}
)

func LarkEventHandler(c *gin.Context) {
	log := wlog.ByCtx(c, "LarkEventHandler").WithField("schema", "-")
	callback := new(Callback[struct{}])

	if err := c.BindJSON(callback); err != nil {
		log.WithError(err).Error("Failed to bind JSON")
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Failed to bind JSON")
		return
	}

	log = log.WithFields(logrus.Fields{
		"event_type": callback.Header.EventType,
		"event_id":   callback.Header.EventID,
		"schema":     callback.Schema,
	})

	log.Info("Event received")

	if callback.Type == "url_verification" {
		log.Info("Processing URL verification")
		c.JSON(http.StatusOK, ChallengeResp{Challenge: callback.Challenge})
		return
	}

	if callback.Schema != "2.0" {
		log.Warn("Unsupported API schema version")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported API schema version"})
		return
	}

	cacheKey := "larkoapi_event_" + callback.Header.EventID
	if err := cache.Client().Set(c, cacheKey, callback.Header.EventType, time.Hour).Err(); err != nil {
		log.WithError(err).Error("Failed to cache event")
	} else {
		log.WithField("cache_key", cacheKey).Info("Event cached")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event processed"})
	log.Info("Event processing completed")
}
