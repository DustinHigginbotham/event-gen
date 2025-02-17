package eventgen

import (
	"strings"
	"text/template"
)

var templateFunctions = template.FuncMap{
	"HasReactors": func(d []DomainSchema) bool {

		for _, domain := range d {
			if len(domain.Reactors) > 0 {
				return true
			}
		}
		return false
	},
	"toLower": func(s string) string {
		return strings.ToLower(s)
	},
	"toType": func(s string) string {
		switch s {
		case "string":
			return "string"
		case "int":
			return "int64"
		default:
			if s == "id" {
				return "ID"
			}
			parts := strings.Split(s, "_")
			for i, part := range parts {
				parts[i] = strings.Title(part)
			}
			return strings.Join(parts, "")
		}
	},
	"toExported": func(s string) string {
		if s == "id" {
			return "ID"
		}
		parts := strings.Split(s, "_")
		for i, part := range parts {
			parts[i] = strings.Title(part)
		}
		return strings.Join(parts, "")
	},
}
