package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sparkymat/resty/verb"
	shttp "github.com/sparkymat/webdsl/http"
)

type router struct {
	resourceHandlers []resourceHandler
	handlerFuncs     []handlerConfig
	muxRouter        *mux.Router
	debugMode        bool
}

func New() *router {
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

func (router *router) Head(path string, handlerFunc HandlerFunc) {
	config := handlerConfig{path: path, handler: handlerFunc, method: shttp.Head}

	router.handlerFuncs = append(router.handlerFuncs, config)
}

func (router *router) Get(path string, handlerFunc HandlerFunc) {
	config := handlerConfig{path: path, handler: handlerFunc, method: shttp.Get}

	router.handlerFuncs = append(router.handlerFuncs, config)
}

func (router *router) Post(path string, handlerFunc HandlerFunc) {
	config := handlerConfig{path: path, handler: handlerFunc, method: shttp.Post}

	router.handlerFuncs = append(router.handlerFuncs, config)
}

func (router *router) Put(path string, handlerFunc HandlerFunc) {
	config := handlerConfig{path: path, handler: handlerFunc, method: shttp.Put}

	router.handlerFuncs = append(router.handlerFuncs, config)
}

func (router *router) Patch(path string, handlerFunc HandlerFunc) {
	config := handlerConfig{path: path, handler: handlerFunc, method: shttp.Patch}

	router.handlerFuncs = append(router.handlerFuncs, config)
}

func (router *router) Delete(path string, handlerFunc HandlerFunc) {
	config := handlerConfig{path: path, handler: handlerFunc, method: shttp.Delete}

	router.handlerFuncs = append(router.handlerFuncs, config)
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
	handler.verbs = []verb.Verb{verb.Create, verb.Update, verb.Show, verb.Index, verb.Destroy}

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

		log.Print("[LOG] ---")
		log.Printf("[LOG] Incoming request: %v on %v with params: {%v}", request.Method, request.RequestURI, strings.Join(params, " , "))
	}

	for _, resource := range router.resourceHandlers {
		startTime := time.Now()
		var handled = resource.checkAndHandleRequest(router.muxRouter, response, request)
		if handled {
			endTime := time.Now()
			duration := endTime.Sub(startTime)
			log.Printf("[LOG] Request handled in %.3f seconds", duration.Seconds())
			return
		}
	}

	for _, handlerConfig := range router.handlerFuncs {
		startTime := time.Now()
		var handled = handlerConfig.checkAndHandleRequest(router.muxRouter, response, request)
		if handled {
			endTime := time.Now()
			duration := endTime.Sub(startTime)
			log.Printf("[LOG] Request handled in %.3f seconds", duration.Seconds())
			return
		}
	}

	log.Printf("[LOG] Warning! No matching resource or route. Passing request to internal Gorilla Mux router.")

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

	for _, handlerConfig := range router.handlerFuncs {
		handlerConfig.PrintRoute(writer)
	}
}

func collectParams(request *http.Request, match mux.RouteMatch) map[string][]string {
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
