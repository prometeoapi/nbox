package usecases

import (
	"context"
	"nbox/internal/domain"
	"strings"
)

type BoxUseCase struct {
	templateAdapter domain.TemplateAdapter
	entryAdapter    domain.EntryAdapter
	pathUseCase     *PathUseCase
}

func NewBox(boxOperation domain.TemplateAdapter, entryOperations domain.EntryAdapter, pathUseCase *PathUseCase) *BoxUseCase {
	return &BoxUseCase{
		templateAdapter: boxOperation,
		entryAdapter:    entryOperations,
		pathUseCase:     pathUseCase,
	}
}

func (b *BoxUseCase) BuildBox(ctx context.Context, service string, stage string, template string) (string, error) {
	box, err := b.templateAdapter.RetrieveBox(ctx, service, stage, template)
	if err != nil {
		return "", err
	}

	tmpl := strings.NewReplacer(":stage", stage, ":service", service).Replace(string(box))
	proc := NewProcessor(tmpl)
	prefixes := proc.GetPrefixes()

	tree := map[string]string{}

	for _, k := range prefixes {
		entries, _ := b.entryAdapter.List(ctx, k)
		for _, entry := range entries {
			if k == strings.TrimSpace(entry.Path) {
				p := b.pathUseCase.Concat(k, entry.Key)
				tree[p] = entry.Value
			}
		}
	}

	return proc.Replace(tree), nil
}
