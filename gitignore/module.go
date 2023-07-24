package gitignore

import (
	_ "embed"
	"sort"
	"strings"

	"github.com/getsyncer/syncer/sharedapi/files/existingfileparser"

	"github.com/getsyncer/syncer/sharedapi/drift/templatefiles"
	"github.com/getsyncer/syncer/sharedapi/syncer"
)

func init() {
	syncer.FxRegister(Module)
}

type Config struct {
	Ignores        []string `yaml:"line"`
	postAutogenMsg string
	preAutogenMsg  string
}

func (c Config) SectionStart() string {
	return existingfileparser.RecommendedSectionStart
}

func (c Config) SectionEnd() string {
	return existingfileparser.RecommendedSectionEnd
}

func (c Config) PostAutogenMsg() string {
	return c.postAutogenMsg
}
func (c Config) PreAutogenMsg() string {
	return c.preAutogenMsg
}

func (c Config) UniqueLines() []string {
	seen := map[string]struct{}{}
	var ret []string
	for _, ignoreLine := range c.Ignores {
		ignoreLine = strings.TrimSpace(ignoreLine)
		if ignoreLine == "" {
			continue
		}
		if _, ok := seen[ignoreLine]; ok {
			continue
		}
		seen[ignoreLine] = struct{}{}
		ret = append(ret, ignoreLine)
	}
	sort.Strings(ret)
	return ret
}

func (c Config) ApplyParse(parse *existingfileparser.ParseResult) (Config, error) {
	c.preAutogenMsg = parse.PreAutogenMsg
	c.postAutogenMsg = parse.PostAutogenMsg
	return c, nil
}

//go:embed .gitignore.template
var templateStr string

var Module = templatefiles.NewModule(templatefiles.NewModuleConfig[Config]{
	Name: "gitignore",
	Files: map[string]string{
		".gitignore": templateStr,
	},
	Priority: syncer.PriorityNormal,
	Decoder:  templatefiles.DefaultDecoder[Config](),
	Setup: &syncer.SetupMutator[Config]{
		Name: "gitignore",
		Mutator: &syncer.ParserMutator[Config]{
			Path:  ".gitignore",
			Conf:  existingfileparser.RecommendedNewlineSeparatedConfig(),
			Apply: syncer.ConfigApply[Config](),
		},
	},
})
