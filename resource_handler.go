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
	Index(response http.ResponseWriter, request *http.Request, params map[string][]string)
	Create(response http.ResponseWriter, request *http.Request, params map[string][]string)
	Show(response http.ResponseWriter, request *http.Request, params map[string][]string)
	Update(response http.ResponseWriter, request *http.Request, params map[string][]string)
	Destroy(response http.ResponseWriter, request *http.Request, params map[string][]string)
}

type resourceHandler struct {
	ParentChain []string
	Name        string
	controller  ResourceController
	router      *mux.Router
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

func (handler resourceHandler) deriveParams(request *http.Request, match mux.RouteMatch) map[string][]string {
	params := make(map[string][]string)

	for key, value := range match.Vars {
		var values []string
		values = append(values, value)
		params[key] = values
	}

	request.ParseForm()
	for key, value := range request.Form {
		params[key] = value
	}

	return params
}

func (handler resourceHandler) PrintRoutes(writer io.Writer) {
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Get, handler.CollectionRoute(), reflect.TypeOf(handler.controller), "Index")
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Put, handler.CollectionRoute(), reflect.TypeOf(handler.controller), "Create")
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Get, handler.MemberRoute(), reflect.TypeOf(handler.controller), "Show")
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Patch, handler.MemberRoute(), reflect.TypeOf(handler.controller), "Update")
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Post, handler.MemberRoute(), reflect.TypeOf(handler.controller), "Update")
	fmt.Fprintf(writer, "%v\t%v\t\t%v#%v\n", shttp.Delete, handler.MemberRoute(), reflect.TypeOf(handler.controller), "Destroy")
}
