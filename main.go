package main

import (
	"fingelpp/access"
	"fingelpp/api"
	"fingelpp/parser"
	"fingelpp/utils"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func createHTMLRenderer(rootDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	funcMap := template.FuncMap{
		"nl2br": func(text string) template.HTML {
			escaped := template.HTMLEscapeString(text)
			withBreaks := strings.ReplaceAll(escaped, "\n", "<br>")
			return template.HTML(withBreaks)
		},
	}

	globals := []string{rootDir + "/base.tmpl"}

	pages, err := filepath.Glob(rootDir + "/pages/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	for _, page := range pages {
		files := make([]string, len(globals)+1)
		copy(files, globals)
		files[len(files)-1] = page
		r.AddFromFilesFuncsWithOptions(filepath.Base(page), funcMap, multitemplate.TemplateOptions{}, files...)
	}
	return r
}

func main() {
	r := gin.Default()
	r.HTMLRender = createHTMLRenderer("./www/templates")
	r.Static("/static", "./www/static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.tmpl", parser.CurrentBook.Chapters)
	})

	r.GET("/lessons/:id", func(c *gin.Context) {
		id, err := parser.ParseLessonId(c.Param("id"))
		if err != nil {
			utils.ReqError(c, http.StatusBadRequest)
			return
		}

		lesson := parser.CurrentBook.GetLessonById(id)
		if lesson == nil {
			utils.ReqError(c, http.StatusNotFound)
			return
		}

		editPerms := access.CurrentAccessFile.HasPermission(c, access.PermissionEditLesson)

		chap := parser.CurrentBook.GetChapterById(lesson.Id.ChapterId())

		c.HTML(http.StatusOK, "lesson.tmpl", gin.H{"Lesson": lesson, "ChapterName": chap.Name, "ChapterId": lesson.Id.ChapterId(), "EditPerms": editPerms})
	})

	api.Routes(r)

	r.Run("localhost:2025")
}
