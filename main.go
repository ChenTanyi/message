package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/gin-gonic/gin"
)

type Message string

func (m Message) Message() string {
	return string(m)
}

var messages = make(map[string]Message)
var messageLock = &sync.RWMutex{}

func messageHandler(c *gin.Context) {
	receive := strings.TrimSpace(c.Query("message"))
	path := c.Request.URL.Path
	if path[len(path)-1] != '/' {
		c.Redirect(307, fmt.Sprintf("%s/?%s", path, c.Request.URL.RawQuery))
		return
	}

	if receive != "" {
		messageLock.Lock()
		defer messageLock.Unlock()
		if len(messages) > 1000 {
			messages = make(map[string]Message)
		}
		messages[path] = Message(receive)

		query := c.Request.URL.Query()
		query.Del("message")
		c.Redirect(307, fmt.Sprintf("./?%s", query.Encode()))
		return
	}
	messageLock.RLock()
	defer messageLock.RUnlock()

	tpl := template.Must(template.New("Message").Parse(string(MustAsset("template/message.gohtml"))))
	message, _ := messages[path]
	if err := tpl.Execute(c.Writer, message); err != nil {
		fmt.Println(err)
	}
}

func main() {
	port := flag.Int("p", 80, "port")
	baseURI := flag.String("base", "", "base uri")

	r := gin.Default()
	group := r.Group(*baseURI)
	group.Any("/*message", messageHandler)

	if err := r.Run(fmt.Sprintf(":%d", *port)); err != nil {
		fmt.Println(err)
	}
}
