package lessons

import (
	"fingelpp/access"
	"fingelpp/finsyn"
	"fingelpp/utils"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	accessFile := access.CurrentAccessFile

	r.GET("/lessons/:id", func(c *gin.Context) {
		id, err := ParseLessonId(c.Param("id"))
		if err != nil {
			utils.ReqError(c, http.StatusBadRequest)
			return
		}

		lesson := CurrentBook.GetLessonById(id)
		if lesson == nil {
			utils.ReqError(c, http.StatusNotFound)
			return
		}

		chap := CurrentBook.GetChapterById(lesson.Id.ChapterId())

		c.HTML(http.StatusOK, "lesson.tmpl", gin.H{"Lesson": lesson, "ChapterName": chap.Name, "ChapterId": lesson.Id.ChapterId(), "Edit": false})
	})

	r.GET("/lessons/:id/edit", func(c *gin.Context) {
		if accessFile.EnforcePermission(c, access.PermissionEditLesson) {
			return
		}
		id, err := ParseLessonId(c.Param("id"))
		if err != nil {
			utils.ReqError(c, http.StatusBadRequest)
			return
		}

		lesson := CurrentBook.GetLessonById(id)
		if lesson == nil {
			utils.ReqError(c, http.StatusNotFound)
			return
		}

		chap := CurrentBook.GetChapterById(lesson.Id.ChapterId())

		c.HTML(http.StatusOK, "lesson.tmpl", gin.H{"Lesson": lesson, "ChapterName": chap.Name, "ChapterId": lesson.Id.ChapterId(), "Edit": true})

	})

	r.POST("/lessons/reload", func(c *gin.Context) {
		if accessFile.EnforcePermission(c, access.PermissionEditLesson) {
			return
		}
		CurrentBook.Reload()
	})

	r.POST("/lessons/:id/reload", func(c *gin.Context) {
		if accessFile.EnforcePermission(c, access.PermissionEditLesson) {
			return
		}

		id, err := ParseLessonId(c.Param("id"))
		if err != nil {
			utils.ReqError(c, http.StatusBadRequest)
			return
		}

		log.Info("Reloading a lesson via api call")
		CurrentBook.GetLessonById(id).Reload()
	})

	r.POST("/lessons/preview", func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, "Error reading body")
			return
		}
		parsedBody := string(finsyn.ParseFinSyn(string(bodyBytes)))
		c.String(http.StatusOK, parsedBody)
	})

	r.PUT("/lessons/:id", func(c *gin.Context) {
		id, err := ParseLessonId(c.Param("id"))
		if err != nil {
			utils.ReqError(c, http.StatusBadRequest)
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, "Error reading body")
			return
		}
		lesson := CurrentBook.GetLessonById(id)
		lesson.Update(string(bodyBytes))
	})

	r.DELETE("/lessons/:id", func(c *gin.Context) {

	})
}
