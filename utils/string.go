package utils

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/go-openapi/inflect"
)

var (
	rules    = ruleset()
	acronyms = make(map[string]struct{})
)

// CollectionName returns the name of the collection from struct name
func CollectionName(model any) string {
	name := StructName(model)

	// Schema package is lower-cased (see Type.Package).
	return SnakeString(rules.Pluralize(name))
}

// StructName returns the name of the struct
func StructName(model any) string {
	name := reflect.TypeOf(model).Name()
	if name != "" {
		return name
	}

	return reflect.TypeOf(model).Elem().Name()
}

// SnakeString converts the given struct or field name into a snake_case.
//
//	Username => username
//	FullName => full_name
//	HTTPCode => http_code
func SnakeString(s string) string {
	var (
		j int
		b strings.Builder
	)
	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		// Put '_' if it is not a start or end of a word, current letter is uppercase,
		// and previous is lowercase (cases like: "UserInfo"), or next letter is also
		// a lowercase and previous letter is not "_".
		if i > 0 && i < len(s)-1 && unicode.IsUpper(r) {
			if unicode.IsLower(rune(s[i-1])) ||
				j != i-1 && unicode.IsLower(rune(s[i+1])) && unicode.IsLetter(rune(s[i-1])) {
				j = i
				b.WriteString("_")
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// StringCompareOrPassValue returns string compare or pass value
func StringCompareOrPassValue(a, b string) string {
	if strings.Compare(a, b) != 0 {
		return b
	}

	return a
}

// ruleset returns the default ruleset.
func ruleset() *inflect.Ruleset {
	rule := inflect.NewDefaultRuleset()
	// Add common initialisms from golint and more.
	for _, w := range []string{
		"ACL", "API", "ASCII", "AWS", "CPU", "CSS", "DNS", "EOF", "GB", "GUID",
		"HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "KB", "LHS", "MAC", "MB",
		"QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "SSO", "TCP",
		"TLS", "TTL", "UDP", "UI", "UID", "URI", "URL", "UTF8", "UUID", "VM",
		"XML", "XMPP", "XSRF", "XSS",
	} {
		acronyms[w] = struct{}{}
		rule.AddAcronym(w)
	}
	return rule
}
