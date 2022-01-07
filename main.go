package main

import (
	"Gara/gara"
	"net/http"
)

func main() {
	engine := gara.NewEngine()
	engine.GET("/", func(c *gara.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gara</h1>")
	})

	engine.GET("/hello", func(c *gara.Context) {
		// expect /hello?name=gara
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	engine.GET("/hello/:name", func(c *gara.Context) {
		// expect /hello/gara
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	engine.GET("/assets/*filepath", func(c *gara.Context) {
		c.JSON(http.StatusOK, gara.H{"filepath": c.Param("filepath")})
	})

	engine.Run(":9999")
}
