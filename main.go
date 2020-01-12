package main

import (
	"flag"
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
)

type Message string

func (m Message) Message() string {
	return string(m)
}

var message Message

func messageHandler(c *gin.Context) {
	receive := c.Query("message")
	if receive != "" {
		message = Message(receive)

		query := c.Request.URL.Query()
		query.Del("message")
		c.Redirect(302, fmt.Sprintf("./?%s", query.Encode()))
		return
	}

	tpl := template.Must(template.New("Message").Parse(string(MustAsset("template/message.gohtml"))))
	if err := tpl.Execute(c.Writer, message); err != nil {
		fmt.Println(err)
	}
}

func main() {
	port := flag.Int("p", 80, "port")
	baseURI := flag.String("base", "", "base uri")

	r := gin.Default()
	r.Any(*baseURI, messageHandler)

	if err := r.Run(fmt.Sprintf(":%d", *port)); err != nil {
		fmt.Println(err)
	}
}
