package utils

import "github.com/gin-gonic/gin"

func ReqError(c *gin.Context, code int) {
	c.HTML(code, "error.tmpl", code)
}
