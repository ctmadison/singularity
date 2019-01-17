// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	osexec "os/exec"
	"runtime"

	"github.com/sylabs/singularity/internal/pkg/buildcfg"
)

func startVm(sifImage, singAction, cliExtra string, isInternal bool) error {

	var stdoutBuf, stderrBuf bytes.Buffer
	hdString := fmt.Sprintf("2:0,ahci-hd,%s", sifImage)

	bzImage := fmt.Sprintf(buildcfg.LIBEXECDIR+"%s"+runtime.GOARCH, "/singularity/vm/syos-kernel-")
	initramfs := fmt.Sprintf(buildcfg.LIBEXECDIR+"%s"+runtime.GOARCH+".gz", "/singularity/vm/initramfs_")
	kexecArgs := fmt.Sprintf("kexec,%s,%s,console=ttyS0 quiet root=/dev/ram0 singularity_action=%s singularity_arguments=\"%s\"", bzImage, initramfs, singAction, cliExtra)

	defArgs := []string{""}
	if cliExtra == "syos" && isInternal {
		//fmt.Println("defArgs - without -hda")
		defArgs = []string{"-A", "-m", "6G", "-c", "2", "-s", "0:0,hostbridge", "-s", "31,lpc", "-l", "com1,stdio", "-f", kexecArgs}
	} else {
		//fmt.Println("defArgs - with -hda")
		defArgs = []string{"-A", "-m", "6G", "-c", "2", "-s", "0:0,hostbridge", "-s", hdString, "-s", "31,lpc", "-l", "com1,stdio", "-f", kexecArgs}
	}

	pgmExec, lookErr := osexec.LookPath("/usr/local/libexec/xhyve/build/xhyve")
	if lookErr != nil {
		panic(lookErr)
	}

	cmd := osexec.Command(pgmExec, defArgs...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	cmdErr := cmd.Run()
	if cmdErr != nil {
		//log.Infof("cmd.Start() failed with '%s'\n", cmdErr)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}

	return nil
}
