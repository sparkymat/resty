package resty

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

	"bitbucket.org/pkg/inflect"
	"github.com/gorilla/mux"
	shttp "github.com/sparkymat/webdsl/http"
)

type ResourceController interface {
	Index(response http.ResponseWriter, request *http.Request, params map[string]string)
	Create(response http.ResponseWriter, request *http.Request, params map[string]string)
	Show(response http.ResponseWriter, request *http.Request, params map[string]string)
	Update(response http.ResponseWriter, request *http.Request, params map[string]string)
	Destroy(response http.ResponseWriter, request *http.Request, params map[string]string)
}

type resourceHandler struct {
	ParentChain []string
	Name        string
	controller  ResourceController
	router      *mux.Router
	NextHandler http.Handler
}

func Resource(path ...string) resourceHandler {
	handler := resourceHandler{}

	if len(path) == 0 {
		return handler
	}

	handler.Name = path[len(path)-1]
	handler.ParentChain = path[:len(path)-1]
	handler.router = mux.NewRouter()

	return handler
}

func (handler resourceHandler) Controller(controller ResourceController) resourceHandler {
	handler.controller = controller
	return handler
}

func (handler resourceHandler) HandleRoot() {
	http.Handle("/", handler)
}

func (handler resourceHandler) pathPrefix() string {
	var prefix string
	for _, parentPath := range handler.ParentChain {
		singularParentPath := inflect.Singularize(parentPath)
		prefix = fmt.Sprintf("%v/%v/{%v_id:[0-9]+}", prefix, parentPath, singularParentPath)
	}
	return prefix
}

func (handler resourceHandler) CollectionRoute() string {
	return fmt.Sprintf("%v/%v.json", handler.pathPrefix(), handler.Name)
}

func (handler resourceHandler) MemberRoute() string {
	return fmt.Sprintf("%v/%v/{id:[0-9]+}.json", handler.pathPrefix(), handler.Name)
}

func (handler resourceHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if handler.router == nil {
		handler.router = mux.NewRouter()
	}

	var match mux.RouteMatch
	var route *mux.Route

	route = handler.router.NewRoute().Path(handler.MemberRoute())
	if route.Match(request, &match) {
		if request.Method == string(shttp.Get) {
			handler.controller.Show(response, request, match.Vars)
			return
		} else if request.Method == string(shttp.Patch) || request.Method == string(shttp.Post) {
			handler.controller.Update(response, request, match.Vars)
			return
		} else if request.Method == string(shttp.Delete) {
			handler.controller.Destroy(response, request, match.Vars)
			return
		}
	}

	route = handler.router.NewRoute().Path(handler.CollectionRoute())
	if route.Match(request, &match) {
		if request.Method == string(shttp.Get) {
			handler.controller.Index(response, request, match.Vars)
			return
		} else if request.Method == string(shttp.Put) {
			handler.controller.Create(response, request, match.Vars)
			return
		}
	}

	if handler.NextHandler != nil {
		handler.NextHandler.ServeHTTP(response, request)
	} else {
		http.Error(response, "Page not found", http.StatusNotFound)
	}
}

func (handler resourceHandler) PrintRoutes(writer io.Writer) {
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Get, handler.CollectionRoute(), reflect.TypeOf(handler.controller), "Index")
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Put, handler.CollectionRoute(), reflect.TypeOf(handler.controller), "Create")
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Get, handler.MemberRoute(), reflect.TypeOf(handler.controller), "Show")
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Patch, handler.MemberRoute(), reflect.TypeOf(handler.controller), "Update")
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Post, handler.MemberRoute(), reflect.TypeOf(handler.controller), "Update")
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Delete, handler.MemberRoute(), reflect.TypeOf(handler.controller), "Destroy")
}
