package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(ctx *Context)

var (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Router %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *RouterGroup) GET(pattern string, handler HandlerFunc) {
	engine.addRoute(GET, pattern, handler)
}

func (engine *RouterGroup) POST(pattern string, handler HandlerFunc) {
	engine.addRoute(POST, pattern, handler)
}

func (engine *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	engine.addRoute(PUT, pattern, handler)
}

func (engine *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	engine.addRoute(DELETE, pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
