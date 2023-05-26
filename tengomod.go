package tengomod

import (
	"context"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/exec"
	"github.com/analog-substance/tengomod/ffuf"
	"github.com/analog-substance/tengomod/filepath"
	"github.com/analog-substance/tengomod/log"
	"github.com/analog-substance/tengomod/net"
	"github.com/analog-substance/tengomod/nmap"
	"github.com/analog-substance/tengomod/os2"
	"github.com/analog-substance/tengomod/set"
	"github.com/analog-substance/tengomod/slice"
	"github.com/analog-substance/tengomod/url"
	"github.com/analog-substance/tengomod/viper"
)

type moduleFactory func(*ModuleOptions) map[string]tengo.Object

var (
	builtinModules map[string]moduleFactory = map[string]moduleFactory{
		"filepath": func(_ *ModuleOptions) map[string]tengo.Object {
			return filepath.Module()
		},
		"viper": func(_ *ModuleOptions) map[string]tengo.Object {
			return viper.Module()
		},
		"url": func(_ *ModuleOptions) map[string]tengo.Object {
			return url.Module()
		},
		"slice": func(_ *ModuleOptions) map[string]tengo.Object {
			return slice.Module()
		},
		"os2": func(o *ModuleOptions) map[string]tengo.Object {
			return os2.Module(o.getCompiled, o.ctx)
		},
		"set": func(_ *ModuleOptions) map[string]tengo.Object {
			return set.Module()
		},
		"nmap": func(_ *ModuleOptions) map[string]tengo.Object {
			return nmap.Module()
		},
		"exec": func(_ *ModuleOptions) map[string]tengo.Object {
			return exec.Module()
		},
		"log": func(_ *ModuleOptions) map[string]tengo.Object {
			return log.Module()
		},
		"ffuf": func(opt *ModuleOptions) map[string]tengo.Object {
			return ffuf.Module(opt.ctx)
		},
		"net": func(_ *ModuleOptions) map[string]tengo.Object {
			return net.Module()
		},
	}
)

type ModuleOptions struct {
	getCompiled func() *tengo.Compiled
	ctx         context.Context
	modules     []string
}

type ModuleOption func(o *ModuleOptions)

func WithCompiled(compiled *tengo.Compiled) ModuleOption {
	return func(o *ModuleOptions) {
		o.getCompiled = func() *tengo.Compiled {
			return compiled
		}
	}
}

func WithCompiledFunc(fn func() *tengo.Compiled) ModuleOption {
	return func(o *ModuleOptions) {
		o.getCompiled = fn
	}
}

func WithContext(ctx context.Context) ModuleOption {
	return func(o *ModuleOptions) {
		o.ctx = ctx
	}
}

func WithModules(modules ...string) ModuleOption {
	return func(o *ModuleOptions) {
		o.modules = modules
	}
}

func WithoutModules(modules ...string) ModuleOption {
	return func(o *ModuleOptions) {
		for _, m := range AllModuleNames() {
			add := true
			for _, module := range modules {
				if m == module {
					add = false
					break
				}
			}

			if add {
				o.modules = append(o.modules, m)
			}
		}
	}
}

func AllModuleNames() []string {
	var names []string
	for name := range builtinModules {
		names = append(names, name)
	}
	return names
}

func GetModuleMap(opts ...ModuleOption) *tengo.ModuleMap {
	options := &ModuleOptions{}
	for _, opt := range opts {
		opt(options)
	}

	modules := options.modules
	if len(modules) == 0 {
		modules = AllModuleNames()
	}

	moduleMap := tengo.NewModuleMap()

	for _, name := range modules {
		factory, ok := builtinModules[name]
		if ok {
			moduleMap.AddBuiltinModule(name, factory(options))
		}
	}

	return moduleMap
}
