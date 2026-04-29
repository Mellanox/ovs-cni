// Copyright 2026 NVIDIA CORPORATION & AFFILIATES
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	defaultCNIBinDir         = "/host/opt/cni/bin"
	defaultOVSBinFile        = "/ovs"
	ovsMirrorProducerBinFile = "/ovs-mirror-producer"
	ovsMirrorConsumerBinFile = "/ovs-mirror-consumer"
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"This is an entrypoint script for OVS CNI to overlay its\n"+
			"binaries into location in a filesystem. Binary files will\n"+
			"be copied to the corresponding CNI bin directory.\n\n"+
			"./entrypoint\n"+
			"\t-h --help\n"+
			"\t--cni-bin-dir=%s\n"+
			"\t--ovs-bin-file=%s\n"+
			"\t--ovs-mirror-producer-bin-file=%s\n"+
			"\t--ovs-mirror-consumer-bin-file=%s\n"+
			"\t--no-sleep\n",
		defaultCNIBinDir, defaultOVSBinFile, ovsMirrorProducerBinFile, ovsMirrorConsumerBinFile)
}

func run() int {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cniBinDir := fs.String("cni-bin-dir", defaultCNIBinDir, "CNI binary destination directory")
	ovsBinFile := fs.String("ovs-bin-file", defaultOVSBinFile, "Source ovs binary path")
	mirrorProducerBinFile := fs.String("ovs-mirror-producer-bin-file", ovsMirrorProducerBinFile, "Source ovs-mirror-producer binary path")
	mirrorConsumerBinFile := fs.String("ovs-mirror-consumer-bin-file", ovsMirrorConsumerBinFile, "Source ovs-mirror-consumer binary path")
	noSleep := fs.Bool("no-sleep", false, "Exit after copying binary without waiting for signal")
	fs.Usage = usage

	err := fs.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to parse flags: %v\n", err)
		return 1
	}

	cniBinDirClean := filepath.Clean(*cniBinDir)
	if !filepath.IsAbs(cniBinDirClean) {
		fmt.Fprintf(os.Stderr, "cni-bin-dir must be an absolute path, got: %s\n", *cniBinDir)
		return 1
	}

	info, err := os.Stat(cniBinDirClean)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cni-bin-dir %q does not exist: %v\n", cniBinDirClean, err)
		return 1
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "cni-bin-dir %q is not a directory\n", cniBinDirClean)
		return 1
	}

	if _, err := os.Stat(*ovsBinFile); err != nil {
		fmt.Fprintf(os.Stderr, "ovs-bin-file %q does not exist: %v\n", *ovsBinFile, err)
		return 1
	}
	if _, err := os.Stat(*mirrorProducerBinFile); err != nil {
		fmt.Fprintf(os.Stderr, "ovs-mirror-producer-bin-file %q does not exist: %v\n", *mirrorProducerBinFile, err)
		return 1
	}
	if _, err := os.Stat(*mirrorConsumerBinFile); err != nil {
		fmt.Fprintf(os.Stderr, "ovs-mirror-consumer-bin-file %q does not exist: %v\n", *mirrorConsumerBinFile, err)
		return 1
	}

	binToCopy := []*string{ovsBinFile, mirrorProducerBinFile, mirrorConsumerBinFile}
	for _, bin := range binToCopy {
		binBase := filepath.Base(*bin)
		destPath := filepath.Join(cniBinDirClean, binBase)
		tempPattern := fmt.Sprintf("%s.temp", binBase)
		if err := copyFileAtomic(*bin, cniBinDirClean, tempPattern, binBase); err != nil {
			fmt.Fprintf(os.Stderr, "failed to copy %q to %q: %v\n", *bin, destPath, err)
			return 1
		}
	}

	if *noSleep {
		fmt.Println("OVS CNI binary installed.")
		return 0
	}

	fmt.Println("OVS CNI binaries installed, waiting for termination signal.")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(ch)
	<-ch
	return 0
}

// copyFileAtomic does file copy atomically
func copyFileAtomic(srcFilePath, destDir, tempFileName, destFileName string) error {
	tempFilePath := filepath.Join(destDir, tempFileName)
	// check temp filepath and remove old file if exists
	if _, err := os.Stat(tempFilePath); err == nil {
		err = os.Remove(tempFilePath)
		if err != nil {
			return fmt.Errorf("cannot remove old temp file %q: %v", tempFilePath, err)
		}
	}

	// create temp file
	f, err := os.CreateTemp(destDir, tempFileName)
	if err != nil {
		return fmt.Errorf("cannot create temp file %q in %q: %v", tempFileName, destDir, err)
	}
	defer f.Close()

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("cannot open file %q: %v", srcFilePath, err)
	}
	defer srcFile.Close()

	// Copy file to tempfile
	_, err = io.Copy(f, srcFile)
	if err != nil {
		f.Close()
		os.Remove(tempFilePath)
		return fmt.Errorf("cannot write data to temp file %q: %v", tempFilePath, err)
	}
	if err = f.Sync(); err != nil {
		return fmt.Errorf("cannot flush temp file %q: %v", tempFilePath, err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("cannot close temp file %q: %v", tempFilePath, err)
	}

	// change file mode if different
	destFilePath := filepath.Join(destDir, destFileName)
	_, err = os.Stat(destFilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	srcFileStat, err := os.Stat(srcFilePath)
	if err != nil {
		return err
	}

	if err := os.Chmod(f.Name(), srcFileStat.Mode()); err != nil {
		return fmt.Errorf("cannot set stat on temp file %q: %v", f.Name(), err)
	}

	// replace file with tempfile
	if err := os.Rename(f.Name(), destFilePath); err != nil {
		return fmt.Errorf("cannot replace %q with temp file %q: %v", destFilePath, tempFilePath, err)
	}

	return nil
}

func main() {
	os.Exit(run())
}
