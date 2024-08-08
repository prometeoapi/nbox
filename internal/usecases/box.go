package usecases

import (
	"context"
	"fmt"
	"nbox/internal/domain"
	"strings"
)

type BoxUseCase struct {
	boxOperation    domain.TemplateAdapter
	entryOperations domain.EntryAdapter
}

func NewBox(boxOperation domain.TemplateAdapter, entryOperations domain.EntryAdapter) *BoxUseCase {
	return &BoxUseCase{
		boxOperation:    boxOperation,
		entryOperations: entryOperations,
	}
}

func (b *BoxUseCase) BuildBox(ctx context.Context, service string, stage string, template string) (string, error) {
	box, err := b.boxOperation.RetrieveBox(ctx, service, stage, template)
	if err != nil {
		return "", err
	}

	tmpl := string(box)
	proc := NewProcessor(tmpl)
	prefixes := proc.GetPrefixes()

	tree := map[string]string{}

	for _, k := range prefixes {
		entries, _ := b.entryOperations.List(ctx, k)
		for _, entry := range entries {
			if entry.Value == "" {
				continue
			}
			if k == "" {
				tree[entry.Key] = entry.Value
				continue
			}
			p := strings.NewReplacer("/", ".").Replace(fmt.Sprintf("%s.%s", k, entry.Key))
			tree[p] = entry.Value
		}
	}

	return proc.Replace(tree), nil
}
