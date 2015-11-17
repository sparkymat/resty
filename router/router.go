package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type router struct {
	resourceHandlers []resourceHandler
	muxRouter        *mux.Router
	debugMode        bool
}

func NewRouter() *router {
	r := router{}
	r.muxRouter = mux.NewRouter()

	return &r
}

func (router *router) EnableDebug() {
	router.debugMode = true
}

func (router *router) DisableDebug() {
	router.debugMode = false
}

func (router *router) HandleFunc(path string, handlerFunc http.HandlerFunc) {
	router.muxRouter.HandleFunc(path, handlerFunc)
}

func (router *router) PathPrefix(tpl string) *mux.Route {
	return router.muxRouter.PathPrefix(tpl)
}

func (router *router) Resource(path []string, controller interface{}) *resourceHandler {
	handler := resourceHandler{}

	if len(path) == 0 {
		return &handler
	}

	handler.name = path[len(path)-1]
	handler.parentChain = path[:len(path)-1]
	handler.router = mux.NewRouter()
	handler.controller = controller
	handler.verbs = []Verb{Create, Update, Show, Index, Destroy}

	router.resourceHandlers = append(router.resourceHandlers, handler)

	return &router.resourceHandlers[len(router.resourceHandlers)-1]
}

func (router router) HandleRoot() {
	http.Handle("/", router)
}

func (router router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	if router.debugMode {
		var params []string
		for key, values := range request.Form {
			if len(values) == 0 {
				params = append(params, fmt.Sprintf("%v: %v", key, "nil"))
			} else if len(values) == 1 {
				params = append(params, fmt.Sprintf("%v: \"%v\"", key, values[0]))
			} else {
				var quotedValues []string
				for _, value := range values {
					quotedValues = append(quotedValues, fmt.Sprintf("\"%v\"", value))
				}
				params = append(params, fmt.Sprintf("%v: [%v]", key, strings.Join(quotedValues, " , ")))
			}
		}

		log.Printf("[LOG] Incoming request: %v on %v with params: {%v}", request.Method, request.RequestURI, strings.Join(params, " , "))
	}

	for _, resource := range router.resourceHandlers {
		var handled = resource.checkAndHandleRequest(router.muxRouter, response, request)
		if handled {
			return
		}
	}

	router.muxRouter.ServeHTTP(response, request)
}

func WriteJSONResponse(response http.ResponseWriter, status int, body map[string]interface{}) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)

	encoder := json.NewEncoder(response)
	encoder.Encode(body)
}

func (router router) PrintRoutes(writer io.Writer) {
	for _, handler := range router.resourceHandlers {
		handler.PrintRoutes(writer)
	}
}
