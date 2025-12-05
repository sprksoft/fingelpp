package api

import (
	"fingelpp/finsyn"
	"fingelpp/parser"
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

func ReqError(c *gin.Context, code int) {
	c.HTML(code, "error.tmpl", code)
}

func RenderPreview(c *gin.Context) { //POST   /lessons/preview   # Render markdown to HTML
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Error reading body")
		return
	}
	parsedBody := string(finsyn.ParseFinSyn(string(bodyBytes)))
	c.String(http.StatusOK, parsedBody)
}

func GetLessonSource(c *gin.Context, b *parser.Book) {
	id, err := parser.ParseLessonId(c.Param("id"))
	if err != nil {
		ReqError(c, http.StatusBadRequest)
		return
	}
	println(b)
	les := b.GetLessonById(id)

	if les == nil {
		ReqError(c, http.StatusNotFound)
		return
	}
	c.String(http.StatusOK, string(les.Src))
}
