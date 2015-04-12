package resty

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type ResourceRouter struct {
	resources []resourceHandler
	router    *mux.Router
}

func (router *ResourceRouter) Resource(path []string, controller interface{}) *resourceHandler {
	handler := resourceHandler{}

	if len(path) == 0 {
		return &handler
	}

	handler.Name = path[len(path)-1]
	handler.ParentChain = path[:len(path)-1]
	handler.router = mux.NewRouter()
	handler.controller = controller
	handler.verbs = []Verb{Create, Update, Show, Index, Destroy}

	router.resources = append(router.resources, handler)

	return &router.resources[len(router.resources)-1]
}

func (router ResourceRouter) HandleRoot() {
	http.Handle("/", router)
}

func (router ResourceRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if router.router == nil {
		router.router = mux.NewRouter()
	}

	for _, resource := range router.resources {
		var handled = resource.checkAndHandleRequest(router.router, response, request)
		if handled {
			return
		}
	}

	http.Error(response, "Page not found", http.StatusNotFound)
}

func (router ResourceRouter) PrintRoutes(writer io.Writer) {
	for _, handler := range router.resources {
		handler.PrintRoutes(writer)
	}
}
