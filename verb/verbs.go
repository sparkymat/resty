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
	Name       string
	Methods    []shttp.Method
	ActionType ActionType
	Format     string
}

var Create = Verb{Name: "create", Methods: []shttp.Method{shttp.Put}, ActionType: CollectionAction}
var Show = Verb{Name: "show", Methods: []shttp.Method{shttp.Get}, ActionType: MemberAction}
var Update = Verb{Name: "update", Methods: []shttp.Method{shttp.Patch, shttp.Post}, ActionType: MemberAction}
var Index = Verb{Name: "index", Methods: []shttp.Method{shttp.Get}, ActionType: CollectionAction}
var Destroy = Verb{Name: "destroy", Methods: []shttp.Method{shttp.Delete}, ActionType: MemberAction}

func (v Verb) Action() string {
	return inflect.Camelize(v.Name)
}

func MemberVerb(methods []shttp.Method, v string) Verb {
	return Verb{Name: v, Methods: methods, ActionType: MemberAction}
}

func CollectionVerb(methods []shttp.Method, v string) Verb {
	return Verb{Name: v, Methods: methods, ActionType: CollectionAction}
}

func (v Verb) RouteSuffix() string {
	switch v.Name {
	case "create", "index":
		return ""
	case "show", "update", "destroy":
		return "/{id:[0-9]+}"
	}

	if v.ActionType == MemberAction {
		return fmt.Sprintf("/{id:[0-9]+}/%v", v.Name)
	} else {
		return fmt.Sprintf("/%v", v.Name)
	}
}
