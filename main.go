package main

import (
	"fingelpp/access"
	"fingelpp/lessons"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/IBM/fp-go/v2/result"
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

	globals := []string{rootDir + "/base.tmpl", rootDir + "/svg/edit.tmpl"}

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

func toResult[T any](value T, err error) result.Result[T] {
	if err != nil {
		return result.Left[T](err)
	}
	return result.Right(value)
}

func main() {
	r := gin.Default()
	r.HTMLRender = createHTMLRenderer("./www/templates")
	r.Static("/static", "./www/static")

	r.GET("/", func(c *gin.Context) {
		// lastOpened := result.Chain(func(value string) result.Result[lessons.LessonId] {
		// 	lessonId, err := lessons.ParseLessonId(value)
		// 	if err != nil {
		// 		return result.Left[lessons.LessonId](err)
		// 	}
		// 	return result.Right(lessonId)
		// })(c.Cookie("lastOpenedLesson"))

		c.HTML(http.StatusOK, "home.tmpl", gin.H{"Chapters": lessons.CurrentBook.Chapters, "EditPerms": access.CurrentAccessFile.HasPermission(c, access.PermissionEditLesson), "LastOpened": ""})
	})

	lessons.Routes(r)
	access.Routes(r)

	r.Run("localhost:2025")
}
