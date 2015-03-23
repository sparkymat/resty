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

type ObjectWithId interface {
	GetId() int64
}

type Resource struct {
	wrappedObject ObjectWithId
}

func MakeResource(object ObjectWithId) Resource {
	return Resource{wrappedObject: object}
}

func (resource Resource) AddRoutes(router *mux.Router) {
	router.HandleFunc(resource.IndexRoute(), func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("got index"))
	}).Methods(string(resource.IndexMethod()))

	router.HandleFunc(resource.ShowRoute(), func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("got show"))
	}).Methods(string(resource.ShowMethod()))
}

func (resource Resource) Type() reflect.Type {
	return reflect.TypeOf(resource.wrappedObject)
}

func (resource Resource) MemberName() string {
	return inflect.Singularize(inflect.Underscore(resource.Type().Name()))
}

func (resource Resource) CollectionName() string {
	return inflect.Pluralize(inflect.Underscore(resource.Type().Name()))
}

func (resource Resource) IndexRoute() string {
	return fmt.Sprintf("/%v.json", resource.CollectionName())
}

func (resource Resource) IndexMethod() shttp.Method {
	return shttp.Get
}

func (resource Resource) ShowRoute() string {
	return fmt.Sprintf("/%v/{id:[0-9]+}.json", resource.CollectionName())
}

func (resource Resource) ShowMethod() shttp.Method {
	return shttp.Get
}

func (resource Resource) WriteRoutes(writer io.Writer) {
	io.WriteString(writer, fmt.Sprintf("%v_index\t%v\t%v\n", resource.CollectionName(), resource.IndexMethod(), resource.IndexRoute()))
	io.WriteString(writer, fmt.Sprintf("%v_show\t%v\t%v\n", resource.MemberName(), resource.ShowMethod(), resource.ShowRoute()))
}
