package usecases

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

type Processor struct {
	tmpl       string
	subPattern string
	vars       []string
}

const (
	ExpressionSingle = `{([^{}]*)}`
	ExpressionDouble = `{{(.*?)}}` // ExpressionDouble double curly braces
)

func NewProcessor(tmpl string) *Processor {
	processor := &Processor{
		tmpl:       tmpl,
		subPattern: ExpressionDouble,
	}
	processor.vars = processor.populateVars()
	return processor
}

func (p *Processor) populateVars() []string {
	var vars []string
	r := regexp.MustCompile(ExpressionDouble)
	matches := r.FindAllStringSubmatch(p.tmpl, -1)
	for _, s := range matches {
		vars = append(vars, s[1])
	}
	return vars
}

func (p *Processor) GetVars() []string {
	return p.vars
}

func (p *Processor) GetPrefixes() []string {
	prefixes := map[string]bool{}
	var k []string
	for _, v := range p.vars {
		cleaned := strings.TrimSpace(v)
		prefix := path.Dir(cleaned)
		if prefix == "." {
			prefix = ""
		}
		if _, ok := prefixes[prefix]; !ok {
			k = append(k, prefix)
			prefixes[prefix] = true
		}
	}
	return k
}

func (p *Processor) Replace(values map[string]string) string {
	var oldnew []string
	for _, v := range p.vars {
		cleaned := strings.TrimSpace(v)
		oldnew = append(oldnew, fmt.Sprintf(`{{%s}}`, v), values[cleaned])
	}
	return strings.NewReplacer(oldnew...).Replace(p.tmpl)
}
