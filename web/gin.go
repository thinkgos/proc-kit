package web

import (
	"context"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.DisableBindValidation()
}

func Handler[T, R any](handle func(context.Context, *T) (*R, error)) gin.HandlerFunc {
	var t T
	bsize := unsafe.Sizeof(t)
	return func(c *gin.Context) {
		var err error
		var req T
		var reply *R

		carrier := FromCarrier(c.Request.Context())
		if bsize > 0 {
			if err = carrier.ShouldAutoBind(c, &req); err != nil {
				carrier.Error(c, err)
				return
			}
		}
		reply, err = handle(c.Request.Context(), &req)
		if err != nil {
			carrier.Error(c, err)
			return
		}
		carrier.Render(c, reply)
	}
}
