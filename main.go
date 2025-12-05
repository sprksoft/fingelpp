package main

import (
	"fingelpp/access"
	"fingelpp/api"
	"fingelpp/parser"
	"fingelpp/utils"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

var book = parser.LoadBook("./content")

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

func main() {
	r := gin.Default()
	r.HTMLRender = createHTMLRenderer("./www/templates")
	r.Static("/static", "./www/static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.tmpl", book.Chapters)
	})

	r.GET("/lessons/:id", func(c *gin.Context) {
		id, err := parser.ParseLessonId(c.Param("id"))
		if err != nil {
			utils.ReqError(c, http.StatusBadRequest)
			return
		}

		lesson := book.GetLessonById(id)
		if lesson == nil {
			utils.ReqError(c, http.StatusNotFound)
			return
		}

		editPerms := access.CurrentAccessFile.HasPermission(c, access.PermissionEditLesson)

		chap := book.GetChapterById(lesson.Id.ChapterId())

		c.HTML(http.StatusOK, "lesson.tmpl", gin.H{"Lesson": lesson, "ChapterName": chap.Name, "ChapterId": lesson.Id.ChapterId(), "EditPerms": editPerms})
	})

	api.Routes(r)

	r.Run("localhost:2025")

}
