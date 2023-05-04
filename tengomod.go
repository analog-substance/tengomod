package tengomod

import (
	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/filepath"
	"github.com/analog-substance/tengomod/url"
	"github.com/analog-substance/tengomod/viper"
)

var (
	builtinModules map[string]map[string]tengo.Object = map[string]map[string]tengo.Object{
		"filepath": filepath.Module(),
		"viper":    viper.Module(),
		"url":      url.Module(),
	}
)

func AllModuleNames() []string {
	var names []string
	for name := range builtinModules {
		names = append(names, name)
	}
	return names
}

func GetModuleMap(modules ...string) *tengo.ModuleMap {
	moduleMap := tengo.NewModuleMap()

	for _, name := range modules {
		module, ok := builtinModules[name]
		if ok {
			moduleMap.AddBuiltinModule(name, module)
		}
	}

	return moduleMap
}
