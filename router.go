package resty

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type router struct {
	resourceHandlers []resourceHandler
	muxRouter        *mux.Router
}

func NewRouter() *router {
	r := router{}
	r.muxRouter = mux.NewRouter()

	return &r
}

func (router *router) MuxRouter() *mux.Router {
	return router.muxRouter
}

func (router *router) HandleFunc(path string, handlerFunc http.HandlerFunc) {
	router.muxRouter.HandleFunc(path, handlerFunc)
}

func (router *router) Resource(path []string, controller interface{}) *resourceHandler {
	handler := resourceHandler{}

	if len(path) == 0 {
		return &handler
	}

	handler.name = path[len(path)-1]
	handler.parentChain = path[:len(path)-1]
	handler.router = mux.NewRouter()
	handler.controller = controller
	handler.verbs = []Verb{Create, Update, Show, Index, Destroy}

	router.resourceHandlers = append(router.resourceHandlers, handler)

	return &router.resourceHandlers[len(router.resourceHandlers)-1]
}

func (router router) HandleRoot() {
	http.Handle("/", router)
}

func (router router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	for _, resource := range router.resourceHandlers {
		var handled = resource.checkAndHandleRequest(router.muxRouter, response, request)
		if handled {
			return
		}
	}

	log.Printf("Gonna mux\n")
	router.muxRouter.ServeHTTP(response, request)
}

func (router router) PrintRoutes(writer io.Writer) {
	for _, handler := range router.resourceHandlers {
		handler.PrintRoutes(writer)
	}
}
