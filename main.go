package main

import (
	"gu"
	"net/http"
)

func main() {
	e := gu.New()
	e.GET("/", func(c *gu.Context) {
		c.HTML(http.StatusOK, "<h1>gu-library</h1>")
	})
	e.GET("/hello", func(c *gu.Context) {
		c.String(http.StatusOK, "hello, you're at %s\n", c.Path)
	})

	e.POST("/login", func(c *gu.Context) {
		c.JSON(http.StatusOK, gu.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	e.GET("/hello/:name", func(c *gu.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	e.GET("/assets/*filepath", func(c *gu.Context) {
		c.JSON(http.StatusOK, gu.H{"filepath": c.Param("filepath")})
	})

	e.Run(":9999")
}
