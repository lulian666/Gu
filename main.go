package main

import (
	"gu"
	"net/http"
)

func main() {
	e := gu.New()
	e.Use(gu.Logger())
	e.GET("/", func(c *gu.Context) {
		c.HTML(http.StatusOK, "<h1>gu-library</h1>")
	})

	v1 := e.Group("/v1")
	{
		v1.GET("/", func(c *gu.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gu.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := e.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gu.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gu.Context) {
			c.JSON(http.StatusOK, gu.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	e.Run(":9999")
}
