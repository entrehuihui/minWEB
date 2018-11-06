package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//MiddlewareFunc 中间件定义
type MiddlewareFunc func(http.ResponseWriter, *http.Request, func())

// MiddlewareServe 中间件集合
type MiddlewareServe struct {
	Handler    http.Handler
	Middleware []MiddlewareFunc
}

//Add 加入中间件
func (middle *MiddlewareServe) Add(funcs ...MiddlewareFunc) {
	for _, middleFunc := range funcs {
		middle.Middleware = append(middle.Middleware, middleFunc)
	}
}

// ServeHTTP http.Handler的ServeHTTP 方法  ==关键加载部分
func (middle *MiddlewareServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i := 0
	var next func()
	next = func() {
		if i < len(middle.Middleware) {
			i++
			middle.Middleware[i-1](w, r, next)
		} else if middle.Handler != nil {
			middle.Handler.ServeHTTP(w, r)
		}
	}
	next()
}

func main() {

	middle := new(MiddlewareServe)
	route := http.NewServeMux()
	route.Handle("/", http.HandlerFunc(login))
	middle.Handler = route
	middle.Add(loginLog)
	fmt.Println(http.ListenAndServe(":3000", middle))
}

func login(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("404"))
	fmt.Println("load here")
	fmt.Fprintf(w, "4041111111111111")
}

func loginLog(w http.ResponseWriter, r *http.Request, next func()) {
	file, err := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	newWrite := io.MultiWriter(file, os.Stdout)
	newLog := log.New(newWrite, "[loginTime]", 5)
	newLog.Println(time.Now())
	next()
}
