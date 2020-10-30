// +build openbsd solaris

/*
 * MinIO Cloud Storage, (C) 2019-2020 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package disk

import (
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

// GetFileSystemAttrs return the file system attribute as string; containing mode,
// uid, gid, uname, Gname, atime, mtime, ctime and md5
func GetFileSystemAttrs(file string) (string, error) {

	st := syscall.Stat_t{}
	err := syscall.Stat(file, &st)
	if err != nil {
		return "", err
	}

	var fileAttr strings.Builder
	fileAttr.WriteString("atime:")
	fileAttr.WriteString(strconv.FormatInt(st.Atim.Sec, 10) + "#" + strconv.FormatInt(st.Atim.Nsec, 10))
	fileAttr.WriteString("/gid:")
	fileAttr.WriteString(strconv.Itoa(int(st.Gid)))

	g, err := user.LookupGroupId(strconv.FormatUint(uint64(st.Gid), 10))
	if err == nil {
		fileAttr.WriteString("/gname:")
		fileAttr.WriteString(g.Name)
	}

	fileAttr.WriteString("/mode:")
	fileAttr.WriteString(strconv.Itoa(int(st.Mode)))
	fileAttr.WriteString("/mtime:")
	fileAttr.WriteString(strconv.FormatInt(st.Mtim.Sec, 10) + "#" + strconv.FormatInt(st.Mtim.Nsec, 10))
	fileAttr.WriteString("/uid:")
	fileAttr.WriteString(strconv.Itoa(int(st.Uid)))

	u, err := user.LookupId(strconv.FormatUint(uint64(st.Uid), 10))
	if err == nil {
		fileAttr.WriteString("/uname:")
		fileAttr.WriteString(u.Username)
	}

	return fileAttr.String(), nil
}
