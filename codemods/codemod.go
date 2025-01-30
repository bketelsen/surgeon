package codemods

type CodeMod interface {
	Apply(source string, target string, match string, args ...string) error
	Validate(source string, target string, match string, args ...string) error
}

var Mods = map[string]CodeMod{}
