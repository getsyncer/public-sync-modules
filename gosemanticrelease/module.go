package gosemanticrelease

import (
	"context"
	_ "embed"
	"sort"

	"github.com/getsyncer/public-sync-modules/buildaction"

	"github.com/getsyncer/public-sync-modules/buildgo"

	"github.com/getsyncer/syncer-core/drift/templatefiles"
	"github.com/getsyncer/syncer-core/syncer"
)

func init() {
	syncer.FxRegister(Module)
}

//go:embed bump_tag_step.yaml.template
var templateStr string

const Name = syncer.Name("gosemanticrelease")

var Module = templatefiles.NewModule(templatefiles.NewModuleConfig[Config]{
	Name: Name,
	Files: map[string]string{
		// Note: Empty string filename is removed by PostGenProcessor
		"": templateStr,
	},
	Priority: buildgo.RunPriority + 1, // Force it to run before buildgo so our mutation is rendered.
	PostGenProcessor: templatefiles.PostGenProcessorList{
		&templatefiles.PostGenConfigMutator[buildgo.Config]{
			ToMutate: buildgo.Name,
			PostGenMutatorFunc: func(_ context.Context, renderedTemplate string, cfg buildgo.Config) (buildgo.Config, error) {
				cfg.Jobs = append(cfg.Jobs, renderedTemplate)
				return cfg, nil
			},
		},
		&templatefiles.PostGenConfigMutator[buildaction.Config]{
			ToMutate: buildaction.Name,
			PostGenMutatorFunc: func(_ context.Context, renderedTemplate string, cfg buildaction.Config) (buildaction.Config, error) {
				cfg.Jobs = append(cfg.Jobs, renderedTemplate)
				return cfg, nil
			},
		},
	},
})

type Config struct {
	RequiredSteps []string `yaml:"required_steps"`
}

func (c Config) AllRequiredSteps() []string {
	ret := make([]string, 0, len(c.RequiredSteps))
	ret = append(ret, "build", "test")
	ret = append(ret, c.RequiredSteps...)
	ret = removeDuplicate(ret)
	sort.Strings(ret)
	return ret
}

func removeDuplicate[T comparable](items []T) []T {
	ret := make([]T, 0, len(items))
	seen := make(map[T]struct{})
	for _, item := range items {
		if _, ok := seen[item]; !ok {
			ret = append(ret, item)
			seen[item] = struct{}{}
		}
	}
	return ret
}
