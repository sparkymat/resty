package router

import "strconv"

type Params map[string][]string

func (p Params) HasKey(key string) bool {
	_, ok := p[key]

	return ok
}

func (p Params) GetStrings(key string) []string {
	return p[key]
}

func (p Params) GetString(key string) string {
	if len(p[key]) > 0 {
		return p[key][0]
	} else {
		return ""
	}
}

func (p Params) GetInt(key string) int {
	if len(p[key]) > 0 {
		intString := p[key][0]
		intValue, _ := strconv.Atoi(intString)
		return intValue
	} else {
		return 0
	}
}
