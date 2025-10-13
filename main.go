package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

type les struct {
	Name string
	Id   int
}

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
		lessen := []les{
			les{Id: 0, Name: "test les 1"},
			les{Id: 1, Name: "test les 2"},
			les{Id: 2, Name: "test les 3"},
		}
		c.HTML(http.StatusOK, "home.tmpl", lessen)

	})

	r.GET("/lesson/:id", func(c *gin.Context) {
		id := c.Param("id")

		c.HTML(http.StatusOK, "lesson.tmpl", id)
	})

	r.Run("localhost:2025")

}
