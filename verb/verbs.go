package verb

import (
	"fmt"

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
	format     string
}

var Create = Verb{name: "create", methods: []shttp.Method{shttp.Put}, actionType: CollectionAction}
var Show = Verb{name: "show", methods: []shttp.Method{shttp.Get}, actionType: MemberAction}
var Update = Verb{name: "update", methods: []shttp.Method{shttp.Patch, shttp.Post}, actionType: MemberAction}
var Index = Verb{name: "index", methods: []shttp.Method{shttp.Get}, actionType: CollectionAction}
var Destroy = Verb{name: "destroy", methods: []shttp.Method{shttp.Delete}, actionType: MemberAction}

func (v Verb) Action() string {
	return inflect.Camelize(v.name)
}

func (v Verb) Methods() []shttp.Method {
	return v.methods
}

func MemberVerb(methods []shttp.Method, v string) Verb {
	return Verb{name: v, methods: methods, actionType: MemberAction}
}

func CollectionVerb(methods []shttp.Method, v string) Verb {
	return Verb{name: v, methods: methods, actionType: CollectionAction}
}

func (v Verb) routeSuffix() string {
	switch v.name {
	case "create", "index":
		return ""
	case "show", "update", "destroy":
		return "/{id:[0-9]+}"
	}

	if v.actionType == MemberAction {
		return fmt.Sprintf("/{id:[0-9]+}/%v", v.name)
	} else {
		return fmt.Sprintf("/%v", v.name)
	}
}
