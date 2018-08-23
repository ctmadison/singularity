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

// GenericInitializer is an interface that all plugins MUST satisfy. This is where
// any plugin specific initialization code will be called.
type GenericInitializer interface {
	Init()
}

// Symbol foo
type Symbol interface{}

var splgExec Symbol
var splgInit Symbol

// Load will open the plugin specified by path relative to the plugin directory. Plugins
// are by default stored at buildcfg.SLIBDIR (LIBEXECDIR + "/singularity/lib")
func Load(path string) error {
	abspath := filepath.Join(buildcfg.LIBEXECDIR, "/singularity/lib/plugins", path)
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
	}
	// sylog.Debugf("symbol %v of plugin %v does satisfy GenericInitializer", MainSymbolName, abspath)

	// test
	// splgInit, err := p.Lookup("Init")
	// if err != nil {
	// 	sylog.Debugf("plugin %v lookup: %v error\n", abspath, mainSymbol)
	// }
	// sylog.Debugf("after plugin %v lookup: %v\n", abspath, mainSymbol)

	// plgInit, ok := splgInit.(GenericInitializer)
	// if !ok {
	// 	sylog.Debugf("Plugin has no 'Init' function")
	// }
	// sylog.Debugf("after plugin %v Init() check\n", abspath)

	// plgInit.Init()
	// if err != nil {
	// 	sylog.Debugf("plugin %v Init() error: %v\n", abspath, err)
	// }
	// sylog.Debugf("after plugin %v Init(): %v\n", abspath, plgInit)

	//	mainSymbol.Init()

	registeredPlugins[abspath] = initializer
	registeredPlugins[abspath].Init()
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
	// sylog.Debugf("syplugin init \n")
}
