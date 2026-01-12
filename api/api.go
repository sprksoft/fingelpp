package api

import (
	"fingelpp/access"
	"fingelpp/finsyn"
	"fingelpp/parser"
	"fingelpp/utils"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
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

func Routes(r *gin.Engine) {
	accessFile := access.LoadAccessFile("access.txt")

	r.POST("/lessons/reload", func(c *gin.Context) {
		if accessFile.EnforcePermission(c, access.PermissionEditLesson) {
			return
		}
		parser.CurrentBook.Reload()
	})

	r.POST("/lessons/:id/reload", func(c *gin.Context) {
		if accessFile.EnforcePermission(c, access.PermissionEditLesson) {
			return
		}

		id, err := parser.ParseLessonId(c.Param("id"))
		if err != nil {
			utils.ReqError(c, http.StatusBadRequest)
			return
		}

		log.Info("Reloading a lesson via api call")
		parser.CurrentBook.GetLessonById(id).Reload()
	})

	r.GET("/access/key/:key", func(c *gin.Context) {
		key := c.Param("key")
		if !accessFile.KeyExist(key) {
			utils.ReqError(c, http.StatusUnauthorized)
			return
		}
		c.SetCookie("AccessKey", key, 100000000, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/")
	})

	r.GET("/access/permissions", func(c *gin.Context) {
		perms := accessFile.GetPerms(access.GetKey(c))
		var sb strings.Builder
		for _, perm := range perms {
			sb.WriteString(string(perm))
		}
		c.String(http.StatusOK, sb.String())
	})

	r.POST("/lessons/preview", RenderPreview)

	r.PUT("/lessons/:id", SaveLesson)

	r.DELETE("/lessons/:id", DeleteLesson)
}

func SaveLesson(c *gin.Context) {
	id, err := parser.ParseLessonId(c.Param("id"))
	if err != nil {
		ReqError(c, http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Error reading body")
		return
	}
	lesson := parser.CurrentBook.GetLessonById(id)
	lesson.Update(string(bodyBytes))
}

func DeleteLesson(c *gin.Context) {
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
