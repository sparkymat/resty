package resty

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

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

type ResourceHandler struct {
	Name        string
	Controller  ResourceController
	router      *mux.Router
	NextHandler http.Handler
}

func (handler ResourceHandler) HandleRoot() {
	http.Handle("/", handler)
}

func (handler ResourceHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if handler.router == nil {
		handler.router = mux.NewRouter()
	}

	var match mux.RouteMatch
	var path string
	var route *mux.Route

	path = fmt.Sprintf("/%v/{id:[0-9]+}.json", handler.Name)
	route = handler.router.NewRoute().Path(path)
	if route.Match(request, &match) {
		if request.Method == string(shttp.Get) {
			handler.Controller.Show(response, request, match.Vars)
			return
		} else if request.Method == string(shttp.Patch) || request.Method == string(shttp.Post) {
			handler.Controller.Update(response, request, match.Vars)
			return
		} else if request.Method == string(shttp.Delete) {
			handler.Controller.Destroy(response, request, match.Vars)
			return
		}
	}

	path = fmt.Sprintf("/%v.json", handler.Name)
	route = handler.router.NewRoute().Path(path)
	if route.Match(request, &match) {
		if request.Method == string(shttp.Get) {
			handler.Controller.Index(response, request, match.Vars)
			return
		} else if request.Method == string(shttp.Put) {
			handler.Controller.Create(response, request, match.Vars)
			return
		}
	}

	if handler.NextHandler != nil {
		handler.NextHandler.ServeHTTP(response, request)
	} else {
		http.Error(response, "Page not found", http.StatusNotFound)
	}
}

func (handler ResourceHandler) PrintRoutes(writer io.Writer) {
	fmt.Fprintf(writer, "%v\t/%v.json\t\t%v#%v\n", shttp.Get, handler.Name, reflect.TypeOf(handler.Controller), "Index")
	fmt.Fprintf(writer, "%v\t/%v.json\t\t%v#%v\n", shttp.Put, handler.Name, reflect.TypeOf(handler.Controller), "Create")
	fmt.Fprintf(writer, "%v\t/%v/{id:[0-9]+}.json\t\t%v#%v\n", shttp.Get, handler.Name, reflect.TypeOf(handler.Controller), "Show")
	fmt.Fprintf(writer, "%v\t/%v/{id:[0-9]+}.json\t\t%v#%v\n", shttp.Patch, handler.Name, reflect.TypeOf(handler.Controller), "Update")
	fmt.Fprintf(writer, "%v\t/%v/{id:[0-9]+}.json\t\t%v#%v\n", shttp.Post, handler.Name, reflect.TypeOf(handler.Controller), "Update")
	fmt.Fprintf(writer, "%v\t/%v/{id:[0-9]+}.json\t\t%v#%v\n", shttp.Delete, handler.Name, reflect.TypeOf(handler.Controller), "Destroy")
}
