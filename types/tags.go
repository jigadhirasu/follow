package types

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func Ctx2Tags(ctx *gin.Context) Tags {
	tags := Tags{
		"user":       ctx.Query("id"),
		"address":    ctx.ClientIP(),
		"user-agent": ctx.Request.UserAgent(),
	}
	return tags
}

func Metadata2Tags(md metadata.MD, ok bool) Tags {
	tags := Tags{}
	if !ok {
		return tags
	}
	for key := range md {
		tags[key] = md.Get(key)[0]
	}
	return tags
}

type Tags map[string]string

func (tg Tags) String(key string) string {
	lkey := strings.ToLower(key)
	return tg[lkey]
}
func (tg Tags) Bytes(key string) []byte {
	lkey := strings.ToLower(key)
	return []byte(tg.String(lkey))
}
func (tg Tags) Int(key string) int {
	lkey := strings.ToLower(key)
	i, _ := strconv.Atoi(tg.String(lkey))
	return i
}
func (tg Tags) Int64(key string) int64 {
	lkey := strings.ToLower(key)
	i, _ := strconv.ParseInt(tg.String(lkey), 10, 64)
	return i
}
func (tg Tags) Float64(key string) float64 {
	lkey := strings.ToLower(key)
	f, _ := strconv.ParseFloat(tg.String(lkey), 64)
	return f
}
