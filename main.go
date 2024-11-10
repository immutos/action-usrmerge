/*
 * Copyright 2024 Damian Peckett <damian@pecke.tt>.
 *
 * Licensed under the Immutos Community Edition License, Version 1.0
 * (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 *    http://immutos.com/licenses/LICENSE-1.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

const (
	// rootfsMountDir is the directory where the root filesystem is mounted.
	rootfsMountDir = "/rootfs"
)

// usrMergeDirectories is the complete list of directories that can be merged into /usr.
var usrMergeDirectories = []string{"/bin", "/lib", "/lib32", "/lib64", "/libo32", "/libx32", "/sbin"}

// MergeUsr merges the /usr directory into the root filesystem
// See: https://wiki.debian.org/UsrMerge
func main() {
	for _, dir := range usrMergeDirectories {
		// The architecture does not have this directory.
		if _, err := os.Stat(filepath.Join(rootfsMountDir, dir)); os.IsNotExist(err) {
			continue
		}

		// The directory is already usr merged.
		if info, err := os.Lstat(filepath.Join(rootfsMountDir, dir)); err == nil && info.Mode()&os.ModeSymlink != 0 {
			log.Println("Directory is already usr merged", dir)
			continue
		}

		canonDir := filepath.Join("/usr", dir)

		log.Println("Merging into /usr", dir, canonDir)

		if err := cp.Copy(filepath.Join(rootfsMountDir, dir), canonDir, cp.Options{OnSymlink: func(src string) cp.SymlinkAction {
			return cp.Shallow
		}}); err != nil {
			log.Fatal("Failed to copy directory", dir, canonDir, err)
		}

		if err := os.RemoveAll(filepath.Join(rootfsMountDir, dir)); err != nil {
			log.Fatal("Failed to remove directory", dir, err)
		}

		if err := os.Symlink(canonDir, dir); err != nil {
			log.Fatal("Failed to symlink directory", dir, canonDir, err)
		}
	}

	// Remove any usr-is-merged pre/postinst scripts.
	// These are no longer needed after the usr merge.
	_ = os.RemoveAll("/var/lib/dpkg/info/usr-is-merged.preinst")
	_ = os.RemoveAll("/var/lib/dpkg/info/usr-is-merged.postinst")
}
