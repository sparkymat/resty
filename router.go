package resty

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	resourceHandlers []resourceHandler
	muxRouter        *mux.Router
}

func (router *Router) Resource(path []string, controller interface{}) *resourceHandler {
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

func (router Router) HandleRoot() {
	http.Handle("/", router)
}

func (router Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if router.muxRouter == nil {
		router.muxRouter = mux.NewRouter()
	}

	for _, resource := range router.resourceHandlers {
		var handled = resource.checkAndHandleRequest(router.muxRouter, response, request)
		if handled {
			return
		}
	}

	http.Error(response, "Page not found", http.StatusNotFound)
}

func (router Router) PrintRoutes(writer io.Writer) {
	for _, handler := range router.resourceHandlers {
		handler.PrintRoutes(writer)
	}
}
