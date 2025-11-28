package main

import (
	"fingelpp/api"
	"fingelpp/parser"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

var lessonManager = parser.LoadLessons("./content")

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

func ReqError(c *gin.Context, code int) {
	c.HTML(code, "error.tmpl", code)
}

func main() {
	r := gin.Default()
	r.HTMLRender = createHTMLRenderer("./www/templates")
	r.Static("/static", "./www/static")

	r.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "home.tmpl", lessonManager.Chapters)
	})

	r.GET(("/reload"), func(c *gin.Context) {
		lessonManager.Reload()
	})

	r.GET("/lesson/:id", func(c *gin.Context) {
		id, err := parser.ParseLessonId(c.Param("id"))
		if err != nil {
			ReqError(c, http.StatusBadRequest)
			return
		}

		lesson := lessonManager.GetLessonById(id)
		if lesson == nil {
			ReqError(c, http.StatusNotFound)
			return
		}

		chap := lessonManager.GetChapterById(lesson.Id.ChapterId())

		c.HTML(http.StatusOK, "lesson.tmpl", gin.H{"Lesson": lesson, "ChapterName": chap.Name, "ChapterId": lesson.Id.ChapterId()})
	})

	r.POST("/lessons/preview", api.RenderPreview)

	r.Run("localhost:2025")

}
