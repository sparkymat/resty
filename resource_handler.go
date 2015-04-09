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

type resourceHandler struct {
	ParentChain []string
	Name        string
	controller  interface{}
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

func (handler resourceHandler) callController(method string, response http.ResponseWriter, request *http.Request, params map[string][]string) {
	var args []reflect.Value
	args = append(args, reflect.ValueOf(response))
	args = append(args, reflect.ValueOf(request))
	args = append(args, reflect.ValueOf(params))

	methodReflection := reflect.ValueOf(handler.controller).MethodByName(method)
	if methodReflection.IsValid() {
		methodReflection.Call(args)
	} else {
		http.Error(response, fmt.Sprintf("No '%v' method on %v", method, reflect.TypeOf(handler.controller).Name()), http.StatusInternalServerError)
	}
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
