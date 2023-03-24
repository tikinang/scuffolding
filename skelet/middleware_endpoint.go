package skelet

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func RestGuest[I any, O any](
	handle func(ctx context.Context, in I) (O, error),
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var in I
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, NewBadRequestException(err.Error()))
			return
		}

		out, err := handle(c, in)
		if err != nil {
			var apiException ApiException
			if !errors.As(err, &apiException) {
				apiException = NewApiExceptionFromError(err)
			}
			c.JSON(apiException.StatusCode, apiException)
			return
		}

		c.JSON(http.StatusOK, out)
	}
}

func HtmlGuestGet[O any](
	handle func(ctx context.Context) (O, error),
	tmpl string,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		out, err := handle(c)
		if err != nil {
			// TODO(mpavlicek): template included here in tikigo
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.HTML(http.StatusOK, tmpl, out)
	}
}
