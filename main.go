package main

import (
	"fingelpp/api"
	"fingelpp/parser"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

var book = parser.LoadBook("./content")
var accessFile = LoadAccessFile("access.txt")

func createHTMLRenderer(rootDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	globals := []string{rootDir + "/base.tmpl"}

	pages, err := filepath.Glob(rootDir + "/pages/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	for _, page := range pages {
		files := make([]string, len(globals)+1)
		copy(files, globals)
		files[len(files)-1] = page
		r.AddFromFiles(filepath.Base(page), files...)
	}
	return r
}

func reqError(c *gin.Context, code int) {
	c.HTML(code, "error.tmpl", code)
}

func main() {
	r := gin.Default()
	r.HTMLRender = createHTMLRenderer("./www/templates")
	r.Static("/static", "./www/static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.tmpl", book.Chapters)
	})

	r.POST("/lessons/reload", func(c *gin.Context) {
		if accessFile.EnforcePermission(c, PermissionEditLesson) {
			return
		}
		book.Reload()
	})

	r.POST("/lesson/:id/reload", func(c *gin.Context) {
		if accessFile.EnforcePermission(c, PermissionEditLesson) {
			return
		}

		id, err := parser.ParseLessonId(c.Param("id"))
		if err != nil {
			reqError(c, http.StatusBadRequest)
			return
		}

		book.GetLessonById(id).Reload()
	})

	r.GET("/lesson/:id", func(c *gin.Context) {
		id, err := parser.ParseLessonId(c.Param("id"))
		if err != nil {
			reqError(c, http.StatusBadRequest)
			return
		}

		lesson := book.GetLessonById(id)
		if lesson == nil {
			reqError(c, http.StatusNotFound)
			return
		}

		editPerms := accessFile.HasPermission(c, PermissionEditLesson)

		chap := book.GetChapterById(lesson.Id.ChapterId())

		c.HTML(http.StatusOK, "lesson.tmpl", gin.H{"Lesson": lesson, "ChapterName": chap.Name, "ChapterId": lesson.Id.ChapterId(), "EditPerms": editPerms})
	})

	r.GET("/access/key/:key", func(c *gin.Context) {
		key := c.Param("key")
		if !accessFile.KeyExist(key) {
			reqError(c, http.StatusUnauthorized)
			return
		}
		c.SetCookie("AccessKey", key, 100000000, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/")
	})

	r.GET("/access/permissions", func(c *gin.Context) {
		perms := accessFile.GetPerms(GetKey(c))
		var sb strings.Builder
		for _, perm := range perms {
			sb.WriteString(string(perm))
		}
		c.String(http.StatusOK, sb.String())
	})

	r.POST("/lessons/preview", api.RenderPreview)

	r.GET("/lessons/:id/src", func(c *gin.Context) {
		api.GetLessonSource(c, book)
	})

	r.Run("localhost:2025")

}
