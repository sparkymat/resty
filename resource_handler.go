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
	parentChain []string
	name        string
	controller  interface{}
	router      *mux.Router
	verbs       []Verb
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
	verb := Verb{name: name, methods: methods, actionType: MemberAction}
	handler.verbs = append(handler.verbs, verb)
	return handler
}

func (handler *resourceHandler) Collection(name string, methods []shttp.Method) *resourceHandler {
	verb := Verb{name: name, methods: methods, actionType: CollectionAction}
	handler.verbs = append(handler.verbs, verb)
	return handler
}

func (handler *resourceHandler) Except(verbs ...Verb) *resourceHandler {
	filteredVerbs := handler.verbs[:0]
	for _, v := range handler.verbs {
		var included = true
		for _, exceptedVerb := range verbs {
			if v.name == exceptedVerb.name {
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

func (handler *resourceHandler) Only(verbs ...Verb) *resourceHandler {
	filteredVerbs := handler.verbs[:0]
	for _, v := range handler.verbs {
		var included = false
		for _, includedVerb := range verbs {
			if v.name == includedVerb.name {
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
		http.Error(response, fmt.Sprintf("No '%v' method on %v", method, reflect.TypeOf(handler.controller).Name()), http.StatusInternalServerError)
	}
}

func (handler resourceHandler) collectParams(request *http.Request, match mux.RouteMatch) map[string][]string {
	params := make(map[string][]string)

	for key, value := range match.Vars {
		var values []string
		values = append(values, value)
		params[key] = values
	}

	for key, value := range request.Form {
		params[key] = value
	}

	return params
}

func (handler resourceHandler) checkAndHandleRequest(router *mux.Router, response http.ResponseWriter, request *http.Request) bool {
	var match mux.RouteMatch
	var route *mux.Route

	for _, v := range handler.verbs {
		route = router.NewRoute().Path(fmt.Sprintf("%v/%v%v", handler.pathPrefix(), handler.name, v.routeSuffix()))
		if route.Match(request, &match) {
			for _, m := range v.Methods() {
				if string(m) == request.Method {
					handler.callController(v.Action(), response, request, handler.collectParams(request, match))
					return true
				}
			}
		}
	}

	return false
}

func (handler resourceHandler) PrintRoutes(writer io.Writer) {
	for _, verb := range handler.verbs {
		for _, method := range verb.methods {
			fmt.Fprintf(writer, "%v\t%v/%v%v\t%v#%v\n", method, handler.pathPrefix(), handler.name, verb.routeSuffix(), reflect.TypeOf(handler.controller), verb.Action())
		}
	}
}
