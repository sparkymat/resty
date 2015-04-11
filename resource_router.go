package resty

import (
	"net/http"

	"github.com/gorilla/mux"
	shttp "github.com/sparkymat/webdsl/http"
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
	handler.methods = []string{"create", "update", "show", "index", "destroy"}

	router.resources = append(router.resources, handler)

	return &router.resources[len(router.resources)-1]
}

func (router ResourceRouter) HandleRoot() {
	http.Handle("/", router)
}

func (router ResourceRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	var match mux.RouteMatch
	var route *mux.Route

	if router.router == nil {
		router.router = mux.NewRouter()
	}

	for _, resource := range router.resources {
		route = router.router.NewRoute().Path(resource.MemberRoute())
		if route.Match(request, &match) {
			if request.Method == string(shttp.Get) {
				resource.callController("Show", response, request, resource.deriveParams(request, match))
				return
			} else if request.Method == string(shttp.Patch) || request.Method == string(shttp.Post) {
				resource.callController("Update", response, request, resource.deriveParams(request, match))
				return
			} else if request.Method == string(shttp.Delete) {
				resource.callController("Destroy", response, request, resource.deriveParams(request, match))
				return
			}
		}

		route = router.router.NewRoute().Path(resource.CollectionRoute())
		if route.Match(request, &match) {
			if request.Method == string(shttp.Get) {
				resource.callController("Index", response, request, resource.deriveParams(request, match))
				return
			} else if request.Method == string(shttp.Put) {
				resource.callController("Create", response, request, resource.deriveParams(request, match))
				return
			}
		}

	}

	http.Error(response, "Page not found", http.StatusNotFound)
}
