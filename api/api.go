package api

import (
	"fingelpp/finsyn"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Lesson API Endpoints
	GET    /lessons/:id       # Get specific lesson

	PUT    /lessons/:id       # Update lesson
	DELETE /lessons/:id       # Delete lesson
*/

func RenderPreview(c *gin.Context) { //POST   /lessons/preview   # Render markdown to HTML
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Error reading body")
		return
	}
	parsedBody := string(finsyn.ParseFinSyn(string(bodyBytes)))
	c.String(http.StatusOK, parsedBody)
}
