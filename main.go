package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"sync"
	"text/template"

	"github.com/gin-gonic/gin"
)

type Message struct {
	message string
	keys    []string
}

func (m *Message) Message() string {
	return m.message
}

func (m *Message) Keys() []string {
	return m.keys
}

var messages = make(map[string]string)
var messageLock = &sync.RWMutex{}

func extractKeys(m map[string]string, columnSize int) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func messageHandler(c *gin.Context) {
	path := c.Request.URL.Path
	if path[len(path)-1] != '/' {
		c.Redirect(307, fmt.Sprintf("%s/?%s", path, c.Request.URL.RawQuery))
		return
	}

	if c.Request.Method == "POST" {
		messagePost(c, path)
	} else {
		messageGet(c, path)
	}
}

func messagePost(c *gin.Context, path string) {
	messageLock.Lock()
	defer messageLock.Unlock()

	receive := strings.TrimSpace(c.PostForm("message"))
	if receive != "" {
		if len(messages) > 1000 {
			messages = make(map[string]string)
		}
		messages[path] = receive
	} else {
		delete(messages, path)
	}

	c.Redirect(303, "./")
}

func messageGet(c *gin.Context, path string) {
	messageLock.RLock()
	defer messageLock.RUnlock()

	tpl := template.Must(template.New("Message").Parse(string(MustAsset("template/message.gohtml"))))
	message := &Message{
		message: messages[path],
		keys:    extractKeys(messages, 5),
	}
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
