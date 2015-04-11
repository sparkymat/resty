package resty

import (
	"bitbucket.org/pkg/inflect"
	shttp "github.com/sparkymat/webdsl/http"
)

type ActionType string

const MemberAction ActionType = "member"
const CollectionAction ActionType = "collection"

type Verb struct {
	name       string
	methods    []shttp.Method
	actionType ActionType
}

var Create = Verb{name: "create", methods: []shttp.Method{shttp.Put}, actionType: CollectionAction}
var Show = Verb{name: "show", methods: []shttp.Method{shttp.Get}, actionType: MemberAction}
var Update = Verb{name: "update", methods: []shttp.Method{shttp.Patch, shttp.Post}, actionType: MemberAction}
var Index = Verb{name: "index", methods: []shttp.Method{shttp.Get}, actionType: CollectionAction}
var Destroy = Verb{name: "destroy", methods: []shttp.Method{shttp.Delete}, actionType: MemberAction}

func (verb Verb) Action() string {
	return inflect.Camelize(verb.name)
}

func (verb Verb) Methods() []shttp.Method {
	return verb.Methods()
}

func MemberVerb(methods []shttp.Method, verb string) Verb {
	return Verb{name: verb, methods: methods, actionType: MemberAction}
}

func CollectionVerb(methods []shttp.Method, verb string) Verb {
	return Verb{name: verb, methods: methods, actionType: CollectionAction}
}
