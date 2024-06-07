package utils

import (
	"net/http"
	"strconv"

	"github.com/khgame/ranger_iam/pkg/authcli"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/khicago/irr"
)

func GinMWParseID() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := wlog.ByCtx(c, "gin_parse_id")
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			GinHandleError(c, log, http.StatusBadRequest, irr.Wrap(err, "id= %s", idStr), "Invalid ID")
			c.Abort()
			return
		}
		c.Set("__parsed_id", UInt64(id))
		c.Next()
	}
}

// GinMustGetID should be used with GinMWParseID
func GinMustGetID(c *gin.Context) (id UInt64) {
	return c.MustGet("__parsed_id").(UInt64)
}

func GinMustGetUserID(c *gin.Context) (id UInt64) {
	idu, exist := authcli.GetUIDFromGinCtx(c)
	if !exist {
		panic("UserID does not exist")
	}
	return UInt64(idu)
}
