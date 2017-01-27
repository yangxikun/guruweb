package guruweb

import (
	"encoding/json"
	"github.com/kisielk/gotool"
	"go/build"
	"go/parser"
	"golang.org/x/tools/go/buildutil"
	"golang.org/x/tools/go/loader"
	"guruweb/output"
)

type config struct {
	loader     loader.Config
	scope      []string
	allPkgs    []string
	loadedPkgs []string
	verbose    bool
}

func (c *config) MarshalJSON() ([]byte, error) {
	jsonObject := make(map[string]interface{})
	jsonObject["inputScopes"] = c.scope
	jsonObject["loadedPkgs"] = c.loadedPkgs
	jsonObject["verbose"] = c.verbose
	return json.Marshal(jsonObject)
}

var defaultConfig = &config{}

func InitConfig(scope []string, verbose bool) {
	if verbose {
		output.Trace("init config")
	}
	defaultConfig.scope = scope
	defaultConfig.allPkgs = gotool.ImportPaths([]string{"all"})
	defaultConfig.loadedPkgs = make([]string, 0)
	defaultConfig.loader.Build = &build.Default
	defaultConfig.loader.Build.CgoEnabled = false
	defaultConfig.loader.AllowErrors = true
	defaultConfig.loader.ParserMode = parser.AllErrors
	defaultConfig.loader.TypeChecker.Error = func(err error) {}
	defaultConfig.loader.ImportPkgs = buildutil.ExpandPatterns(defaultConfig.loader.Build, scope)
	defaultConfig.verbose = verbose
}
