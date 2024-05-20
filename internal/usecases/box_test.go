package usecases

import (
	"fmt"
	"nbox/internal/domain/models"
	"strings"
	"testing"
)

func TestBoxUseCase_BuildBox(t *testing.T) {
	tmpl := `{
  "ENV_1": "{widget-x.development.key}",
  "ENV_2": "{widget-x.development.debug}",
  "GLOBAL_SERVICE": "{widget-x.sentry}",
  "domain": "{private-domain}"
}`
	//values := map[string]string{
	//	"widget-x.development.key":   "test 1",
	//	"widget-x.development.debug": "test 2",
	//	"widget-x.sentry":            "test 3",
	//	"private-domain":             "test 4",
	//}

	tree := map[string]string{}

	entries := []models.Entry{
		{
			Path:  "widget-x/development",
			Key:   "key",
			Value: []byte("test 1"),
		},
		{
			Path:  "widget-x/development",
			Key:   "debug",
			Value: []byte("test 1"),
		},
		{
			Path:  "widget-x",
			Key:   "sentry",
			Value: []byte("test 1"),
		},
		{
			Path:  "",
			Key:   "private-domain",
			Value: []byte("test 1"),
		},
	}

	for _, entry := range entries {
		if entry.Value == nil {
			continue
		}
		if entry.Path == "" {
			tree[entry.Key] = string(entry.Value)
			continue
		}
		p := strings.NewReplacer("/", ".").Replace(fmt.Sprintf("%s.%s", entry.Path, entry.Key))
		tree[p] = string(entry.Value)
	}

	proc := NewProcessor(tmpl)

	var oldnew []string
	for _, v := range proc.GetVars() {
		oldnew = append(oldnew, fmt.Sprintf(`{%s}`, v), tree[v])
	}
	fmt.Println(strings.NewReplacer(oldnew...).Replace(tmpl))
}
