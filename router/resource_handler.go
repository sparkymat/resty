package router

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"bitbucket.org/pkg/inflect"
	"github.com/gorilla/mux"
	"github.com/sparkymat/resty/verb"
	shttp "github.com/sparkymat/webdsl/http"
)

type resourceHandler struct {
	parentChain []string
	name        string
	controller  interface{}
	router      *mux.Router
	verbs       []verb.Verb
}

func (handler resourceHandler) Controller(controller interface{}) resourceHandler {
	handler.controller = controller
	return handler
}

func (handler resourceHandler) pathPrefix() string {
	var prefix string
	for _, parentPath := range handler.parentChain {
		singularParentPath := inflect.Singularize(parentPath)
		prefix = fmt.Sprintf("%v/%v/{%v_id:[0-9]+}", prefix, parentPath, singularParentPath)
	}
	return prefix
}

func (handler *resourceHandler) Member(name string, methods []shttp.Method) *resourceHandler {
	v := verb.Verb{Name: name, Methods: methods, ActionType: verb.MemberAction}
	handler.verbs = append(handler.verbs, v)
	return handler
}

func (handler *resourceHandler) Collection(name string, methods []shttp.Method) *resourceHandler {
	verb := verb.Verb{Name: name, Methods: methods, ActionType: verb.CollectionAction}
	handler.verbs = append(handler.verbs, verb)
	return handler
}

func (handler *resourceHandler) Except(verbs ...verb.Verb) *resourceHandler {
	filteredVerbs := handler.verbs[:0]
	for _, v := range handler.verbs {
		var included = true
		for _, exceptedVerb := range verbs {
			if v.Name == exceptedVerb.Name {
				included = false
			}
		}
		if included {
			filteredVerbs = append(filteredVerbs, v)
		}
	}
	handler.verbs = filteredVerbs
	return handler
}

func (handler *resourceHandler) Only(verbs ...verb.Verb) *resourceHandler {
	filteredVerbs := handler.verbs[:0]
	for _, v := range handler.verbs {
		var included = false
		for _, includedVerb := range verbs {
			if v.Name == includedVerb.Name {
				included = true
			}
		}
		if included {
			filteredVerbs = append(filteredVerbs, v)
		}
	}
	handler.verbs = filteredVerbs
	return handler
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
		log.Printf("[LOG] Error! No '%v' method on %v", method, reflect.TypeOf(handler.controller).Name())
		http.Error(response, fmt.Sprintf("No '%v' method on %v", method, reflect.TypeOf(handler.controller).Name()), http.StatusInternalServerError)
	}
}

func (handler resourceHandler) checkAndHandleRequest(router *mux.Router, response http.ResponseWriter, request *http.Request) bool {
	var match mux.RouteMatch
	var route *mux.Route

	for _, v := range handler.verbs {
		route = router.NewRoute().Path(fmt.Sprintf("%v/%v%v", handler.pathPrefix(), handler.name, v.RouteSuffix()))
		if route.Match(request, &match) {
			for _, m := range v.Methods {
				if string(m) == request.Method {
					handler.callController(v.Action(), response, request, collectParams(request, match))
					return true
				}
			}
		}
	}

	return false
}

func (handler resourceHandler) PrintRoutes(writer io.Writer) {
	for _, verb := range handler.verbs {
		for _, method := range verb.Methods {
			path := fmt.Sprintf("%v/%v%v", handler.pathPrefix(), handler.name, verb.RouteSuffix())
			fmt.Fprintf(writer, "%6s%64s\t%s#%s()\n", method, path, reflect.TypeOf(handler.controller), verb.Action())
		}
	}
}
