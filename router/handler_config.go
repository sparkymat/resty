package router

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
	shttp "github.com/sparkymat/webdsl/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, Params)

type handlerConfig struct {
	path    string
	handler HandlerFunc
	method  shttp.Method
}

func (config handlerConfig) checkAndHandleRequest(router *mux.Router, response http.ResponseWriter, request *http.Request) bool {
	var match mux.RouteMatch
	var route *mux.Route

	route = router.NewRoute().Path(config.path)
	if route.Match(request, &match) && string(config.method) == request.Method {
		config.handler(response, request, collectParams(request, match))
		return true
	}

	return false
}

func (config handlerConfig) PrintRoute(writer io.Writer) {
	functionName := runtime.FuncForPC(reflect.ValueOf(config.handler).Pointer()).Name()
	fWords := strings.Split(functionName, "/")
	functionShortName := fWords[len(fWords)-1]
	fmt.Fprintf(writer, "%6s%64s\t%s()\n", config.method, config.path, functionShortName)
}
