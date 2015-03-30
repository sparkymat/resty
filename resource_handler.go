package resty

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ResourceController interface {
	Index(response http.ResponseWriter, request *http.Request, params map[string]string)
	Show(response http.ResponseWriter, request *http.Request, params map[string]string)
}

type ResourceHandler struct {
	Name       string
	Controller ResourceController
	router     *mux.Router
}

func (handler *ResourceHandler) RegisterRoutes() {
	http.Handle("/users/", handler)
}

func (handler *ResourceHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if handler.router == nil {
		handler.router = mux.NewRouter()
	}

	// Check index
	var match mux.RouteMatch
	var path string
	var route *mux.Route

	path = fmt.Sprintf("/%v.json", handler.Name)
	route = handler.router.NewRoute().Path(path)
	if route.Match(request, &match) {
		handler.Controller.Index(response, request, match.Vars)
		return
	}

	path = fmt.Sprintf("/%v/{id:0-9+}.json", handler.Name)
	route = handler.router.NewRoute().Path(path)
	if route.Match(request, &match) {
		handler.Controller.Index(response, request, match.Vars)
		return
	}

	http.Error(response, "Page not found", http.StatusNotFound)
}
