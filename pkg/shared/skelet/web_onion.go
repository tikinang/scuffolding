package skelet

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RestGuest[Input any, Output any](
	handle func(ctx context.Context, in Input) (Output, error),
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var in Input
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		out, err := handle(c, in)
		if err != nil {
			// FIXME: mpavlicek - translate api errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, out)
	}
}

func HtmlGuestGet[Output any](
	handle func(ctx context.Context) (Output, error),
	tmpl string,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		out, err := handle(c)
		if err != nil {
			// FIXME: mpavlicek - template included here in tikigo
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.HTML(http.StatusOK, tmpl, out)
	}
}
