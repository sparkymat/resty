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
	handler.verbs = []Verb{Create, Update, Show, Index, Destroy}

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
			if (request.Method == string(shttp.Get)) && resource.handlesVerb(Show) {
				resource.callController(Show.Action(), response, request, resource.deriveParams(request, match))
				return
			} else if (request.Method == string(shttp.Patch) || request.Method == string(shttp.Post)) && resource.handlesVerb(Update) {
				resource.callController(Update.Action(), response, request, resource.deriveParams(request, match))
				return
			} else if (request.Method == string(shttp.Delete)) && resource.handlesVerb(Destroy) {
				resource.callController(Destroy.Action(), response, request, resource.deriveParams(request, match))
				return
			}
		}

		route = router.router.NewRoute().Path(resource.CollectionRoute())
		if route.Match(request, &match) {
			if (request.Method == string(shttp.Get)) && resource.handlesVerb(Index) {
				resource.callController(Index.Action(), response, request, resource.deriveParams(request, match))
				return
			} else if (request.Method == string(shttp.Put)) && resource.handlesVerb(Create) {
				resource.callController(Create.Action(), response, request, resource.deriveParams(request, match))
				return
			}
		}

	}

	http.Error(response, "Page not found", http.StatusNotFound)
}
