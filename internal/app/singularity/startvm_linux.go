// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	osexec "os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func startVm(sifImage, singAction, cliExtra string, isInternal bool) error {

	var stdoutBuf, stderrBuf bytes.Buffer
	hdString := fmt.Sprintf("-hda", sifImage)

	bzImage := fmt.Sprintf(buildcfg.LIBEXECDIR+"%s"+runtime.GOARCH, "/singularity/vm/syos-kernel-")
	initramfs := fmt.Sprintf(buildcfg.LIBEXECDIR+"%s"+runtime.GOARCH+".gz", "/singularity/vm/initramfs_")	
	appendArgs := fmt.Sprintf("root=/dev/ram0 console=ttyS0 quiet singularity_action=%s singularity_arguments=\"%s\"", singAction, cliExtra)

	defArgs := []string{""}
	if cliExtra == "syos" && isInternal {
		//fmt.Println("defArgs - without -hda")
		defArgs = []string{"-cpu", "host", "-enable-kvm", "-device", "virtio-rng-pci", "-display", "none", "-realtime", "mlock=on", "-serial", "stdio", "-kernel", %s, "-initrd", %s, "-m", "4096", "-append", bzImage, initramfs, appendArgs}
	} else {
		//fmt.Println("defArgs - with -hda")
		defArgs = []string{"-cpu", "host", "-enable-kvm", "-device", "virtio-rng-pci", "-display", "none", "-realtime", "mlock=on", hdString, "-serial", "stdio", "-kernel", %s, "-initrd", %s, "-m", "4096", "-append", bzImage, initramfs, appendArgs}
	}

	pgmexec, lookErr := osexec.LookPath("/usr/libexec/qemu-kvm")
	if lookErr != nil {
		panic(lookErr)
	}

	cmd := osexec.Command(pgmexec, defArgs...)
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

