package main

import (
	"fmt"
	"gu"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	e := gu.New()
	e.Use(gu.Logger(), gu.Recovery())

	e.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})

	e.LoadHTMLGlob("templates/*")
	e.Static("/assets", "./static")

	stu1 := &student{Name: "Gu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}

	e.GET("/", func(c *gu.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	e.GET("/students", func(c *gu.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gu.H{
			"title":  "gu",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	e.GET("/date", func(c *gu.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gu.H{
			"title": "gu",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	e.GET("/panic", func(c *gu.Context) {
		names := []string{"gu"}
		c.String(http.StatusOK, names[100])
	})

	e.Run(":9999")
}
