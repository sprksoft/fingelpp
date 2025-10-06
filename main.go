package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type les struct {
	Name string
	Id   int
}

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("./www/templates/base.html", "./www/templates/pages/home.html")
	r.LoadHTMLGlob("./www/templates/*")

	r.Static("/static", "./www/static")
	r.GET("/", func(c *gin.Context) {
		lessen := []les{
			les{Id: 0, Name: "test les 1"},
			les{Id: 1, Name: "test les 2"},
			les{Id: 2, Name: "test les 3"},
		}
		c.HTML(http.StatusOK, "pages/home.html", lessen)

	})

	r.Run("localhost:2025")

}
