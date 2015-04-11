package resty

import (
	"bitbucket.org/pkg/inflect"
	shttp "github.com/sparkymat/webdsl/http"
)

type Verb string

const Create Verb = "create"
const Show Verb = "show"
const Update Verb = "update"
const Index Verb = "index"
const Destroy Verb = "destroy"

func (verb Verb) Action() string {
	return inflect.Camelize(string(verb))
}

func (verb Verb) Methods() []shttp.Method {
	switch verb {
	case Create:
		return []shttp.Method{shttp.Put}
	case Update:
		return []shttp.Method{shttp.Patch, shttp.Post}
	case Destroy:
		return []shttp.Method{shttp.Delete}
	}

	return []shttp.Method{shttp.Get}
}

func CustomMethod(verb string) Verb {
	return Verb(verb)
}
