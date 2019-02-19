// Copyright (c) 2018, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package cli

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sylabs/singularity/docs"
	"github.com/sylabs/singularity/internal/app/singularity"
	"github.com/sylabs/singularity/internal/pkg/sylog"
)

const (
	containerPath      = "/home/mibauer/plugin-compile/compile_plugin.sif"
	containedSourceDir = "/go/src/github.com/sylabs/singularity/plugins/"
)

var (
	out string
)

func init() {
	PluginCompileCmd.Flags().StringVarP(&out, "out", "o", "", "")
}

// PluginCompileCmd allows a user to compile a plugin
//
// singularity plugin compile <path> [-o name]
var PluginCompileCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := filepath.Abs(args[0])
		if err != nil {
			sylog.Fatalf("While sanitizing input path: %s", err)
		}
		sourceDir := filepath.Clean(s)

		destSif := out

		if destSif == "" {
			destSif = sifPath(sourceDir)
		}

		sylog.Debugf("sourceDir: %s; sifPath: %s", sourceDir, destSif)
		return singularity.CompilePlugin(sourceDir, destSif)
	},
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),

	Use:     docs.PluginCompileUse,
	Short:   docs.PluginCompileShort,
	Long:    docs.PluginCompileLong,
	Example: docs.PluginCompileExample,
}
