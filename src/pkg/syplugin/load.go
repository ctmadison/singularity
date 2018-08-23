// Copyright (c) 2018, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package syplugin

import (
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/singularityware/singularity/src/pkg/buildcfg"
	"github.com/singularityware/singularity/src/pkg/sylog"
)

const MainSymbolName = "SyPlugin"

// registeredPlugins contains a map of name -> GenericInitializer objects
var registeredPlugins map[string]GenericInitializer

var registeredMountPlugins map[string]GenericMount
var registeredBindPlugins map[string]GenericBind

// GenericInitializer is an interface that all plugins MUST satisfy. This is where
// any plugin specific initialization code will be called.
type GenericInitializer interface {
	Init()
}

// GenericMount is an interface that all mount plugins MUST satisfy. This is where
// any plugin specific mount code will be called.
type GenericMount interface {
	Mount()
}

// GenericBind is an interface that all mount plugins MUST satisfy. This is where
// any plugin specific bind code will be called.
type GenericBind interface {
	Bind()
}

// Symbol foo
type Symbol interface{}

var splgExec Symbol
var splgInit Symbol

// Load will open the plugin specified by path relative to the plugin directory. Plugins
// are by default stored at buildcfg.SLIBDIR (LIBEXECDIR + "/singularity/lib")
func Load(abspath string) error {
	// abspath := filepath.Join(buildcfg.LIBEXECDIR, "/singularity/lib/plugins", path)
	sylog.Debugf("Load: plugin abspath: %v\n", abspath)

	return load(abspath)
}

func load(abspath string) error {
	p, err := plugin.Open(abspath)
	if err != nil {
		sylog.Debugf("Failure to open plugin @ abspath: %v\n", abspath)
		return err
	}
	// sylog.Debugf("after plugin %v open\n", abspath)
	sylog.Debugf("Plugin info: %+v\n", p)

	mainSymbol, err := p.Lookup(MainSymbolName)
	if err != nil {
		sylog.Debugf("plugin %v does not export %v symbol", abspath, MainSymbolName)
		return fmt.Errorf("plugin %v does not export %v symbol", abspath, MainSymbolName)
	}
	// sylog.Debugf("after plugin %v lookup: %v\n", abspath, mainSymbol)

	initializer, ok := mainSymbol.(GenericInitializer)
	if !ok {
		sylog.Debugf("symbol %v of plugin %v does not satisfy GenericInitializer", MainSymbolName, abspath)
		return fmt.Errorf("symbol %v of plugin %v does not satisfy GenericInitializer", MainSymbolName, abspath)
	} else {
		sylog.Debugf("symbol %v of plugin %v does satisfy GenericInitializer", MainSymbolName, abspath)
		registeredPlugins[abspath] = initializer
		registeredPlugins[abspath].Init()
	}

	mount, ok := mainSymbol.(GenericMount)
	if !ok {
		sylog.Debugf("symbol %v of plugin %v does not satisfy GenericMount", MainSymbolName, abspath)
	} else {
		sylog.Debugf("symbol %v of plugin %v does satisfy GenericMount", MainSymbolName, abspath)
		registeredMountPlugins[abspath] = mount
		registeredMountPlugins[abspath].Mount() // test
	}

	bind, ok := mainSymbol.(GenericBind)
	if !ok {
		sylog.Debugf("symbol %v of plugin %v does not satisfy GenericBind", MainSymbolName, abspath)
	} else {
		sylog.Debugf("symbol %v of plugin %v does satisfy GenericBind", MainSymbolName, abspath)
		registeredBindPlugins[abspath] = bind
		registeredBindPlugins[abspath].Bind() // test
	}

	return nil
}

// GetByName returns the GenericInitializer of the given path for use
func GetByName(path string) (GenericInitializer, error) {
	abspath := filepath.Join(buildcfg.LIBEXECDIR, "/singularity/lib/plugins", path)

	return getByName(abspath)
}

func getByName(abspath string) (GenericInitializer, error) {
	if p, ok := registeredPlugins[abspath]; ok {
		return p, nil
	}

	return nil, fmt.Errorf("plugin %v not registered", abspath)
}

func init() {
	registeredPlugins = make(map[string]GenericInitializer)
	registeredMountPlugins = make(map[string]GenericMount)
	registeredBindPlugins = make(map[string]GenericBind)
	// sylog.Debugf("syplugin init \n")
}
