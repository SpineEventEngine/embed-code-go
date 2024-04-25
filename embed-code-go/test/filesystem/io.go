// Copyright 2024, TeamDev. All rights reserved.
//
// Redistribution and use in source and/or binary forms, with or without
// modification, must retain the above copyright notice and the following
// disclaimer.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package filesystem

import (
	"io"
	"os"
	"path/filepath"
)

// Removes directory and all its subdirectories if exists, does nothing if not exists.
//
// dir_path — a path (full or relative) of the directory to be removed.
func CleanupDir(dir_path string) {
	if _, err := os.Stat(dir_path); err == nil {
		err = os.RemoveAll(dir_path)
		if err != nil {
			panic(err)
		}
	}
}

// Copies directory from source path to target path with all subdirs and children.
//
// source_dir_path — a path (full or relative) of the directory to be copied.
//
// target_dir_path — a path (full or relative) of the directory to be copied to.
func CopyDirRecursive(source_dir_path string, target_dir_path string) {
	info, err := os.Stat(source_dir_path)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(target_dir_path, info.Mode())
	if err != nil {
		panic(err)
	}

	entries, err := os.ReadDir(source_dir_path)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(source_dir_path, entry.Name())
		targetPath := filepath.Join(target_dir_path, entry.Name())

		if entry.IsDir() {
			CopyDirRecursive(sourcePath, targetPath)
		} else {
			err = CopyFile(sourcePath, targetPath)
			if err != nil {
				panic(err)
			}
		}
	}
}

// Copies file from source_file_path to target_file_path.
func CopyFile(source_file_path string, target_file_path string) (err error) {
	sourceFile, err := os.Open(source_file_path)
	if err != nil {
		return
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(target_file_path)
	if err != nil {
		return
	}
	defer func() {
		cerr := targetFile.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(targetFile, sourceFile); err != nil {
		return
	}

	err = os.Chmod(target_file_path, os.FileMode(0666))
	return
}
