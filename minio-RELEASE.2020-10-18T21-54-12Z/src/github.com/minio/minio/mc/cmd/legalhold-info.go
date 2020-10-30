/*
 * MinIO Client (C) 2020 MinIO, Inc.
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

package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/minio/cli"
	json "github.com/minio/mc/pkg/colorjson"
	"github.com/minio/mc/pkg/probe"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio/pkg/console"
)

var (
	lhInfoFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "recursive, r",
			Usage: "show legal hold status recursively",
		},
		cli.StringFlag{
			Name:  "version-id, vid",
			Usage: "show legal hold status of a specific object version",
		},
		cli.StringFlag{
			Name:  "rewind",
			Usage: "show legal hold status of an object version at specified time",
		},
		cli.BoolFlag{
			Name:  "versions",
			Usage: "show legal hold status of multiple versions of object(s)",
		},
	}
)
var legalHoldInfoCmd = cli.Command{
	Name:   "info",
	Usage:  "show legal hold info for object(s)",
	Action: mainLegalHoldInfo,
	Before: setGlobalsFromContext,
	Flags:  append(lhInfoFlags, globalFlags...),
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} [FLAGS] TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}

EXAMPLES:
   1. Show legal hold on a specific object
      $ {{.HelpName}} myminio/mybucket/prefix/obj.csv

   2. Show legal hold on a specific object version
      $ {{.HelpName}} myminio/mybucket/prefix/obj.csv --version-id "HiMFUTOowG6ylfNi4LKxD3ieHbgfgrvC"

   3. Show object legal hold recursively for all objects at a prefix
      $ {{.HelpName}} myminio/mybucket/prefix --recursive

   4. Show object legal hold recursively for all objects versions older than one year
      $ {{.HelpName}} myminio/mybucket/prefix --recursive --rewind 365d --versions
`,
}

// Structured message depending on the type of console.
type legalHoldInfoMessage struct {
	LegalHold minio.LegalHoldStatus `json:"legalhold"`
	URLPath   string                `json:"urlpath"`
	Key       string                `json:"key"`
	VersionID string                `json:"versionID"`
	Status    string                `json:"status"`
	Err       error                 `json:"error,omitempty"`
}

// Colorized message for console printing.
func (l legalHoldInfoMessage) String() string {
	if l.Err != nil {
		return console.Colorize("LegalHoldMessageFailure", "Unable to get object legal hold status `"+l.Key+"`. "+l.Err.Error())
	}
	var msg string

	var legalhold string
	switch l.LegalHold {
	case "":
		legalhold = console.Colorize("LegalHoldNotSet", "Not set")
	case minio.LegalHoldEnabled:
		legalhold = console.Colorize("LegalHoldOn", l.LegalHold)
	case minio.LegalHoldDisabled:
		legalhold = console.Colorize("LegalHoldOff", l.LegalHold)
	}

	msg += "[ " + centerText(legalhold, 8) + " ] "

	if l.VersionID != "" {
		msg += " " + console.Colorize("LegalHoldVersion", l.VersionID) + " "
	}

	msg += " "
	msg += l.Key
	return msg
}

// JSON'ified message for scripting.
func (l legalHoldInfoMessage) JSON() string {
	msgBytes, e := json.MarshalIndent(l, "", " ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")
	return string(msgBytes)
}

// showLegalHoldInfo - show legalhold for one or many objects within a given prefix, with or without versioning
func showLegalHoldInfo(ctx context.Context, urlStr, versionID string, timeRef time.Time, withOlderVersions, recursive bool) error {
	clnt, err := newClient(urlStr)
	if err != nil {
		fatalIf(err.Trace(), "Unable to parse the provided url.")
	}

	prefixPath := clnt.GetURL().Path
	prefixPath = filepath.ToSlash(prefixPath)
	if !strings.HasSuffix(prefixPath, "/") {
		prefixPath = prefixPath[:strings.LastIndex(prefixPath, "/")+1]
	}
	prefixPath = strings.TrimPrefix(prefixPath, "./")

	if !recursive && !withOlderVersions {
		lhold, err := clnt.GetObjectLegalHold(ctx, versionID)
		if err != nil {
			fatalIf(err.Trace(urlStr), "Failed to show legal hold information of `"+urlStr+"`.")
		} else {
			contentURL := filepath.ToSlash(clnt.GetURL().Path)
			key := strings.TrimPrefix(contentURL, prefixPath)

			printMsg(legalHoldInfoMessage{
				LegalHold: lhold,
				Status:    "success",
				URLPath:   clnt.GetURL().String(),
				Key:       key,
				VersionID: versionID,
			})
		}
		return nil
	}

	alias, _, _ := mustExpandAlias(urlStr)
	var cErr error
	errorsFound := false
	objectsFound := false
	lstOptions := ListOptions{IsRecursive: recursive, ShowDir: DirNone}
	if !timeRef.IsZero() {
		lstOptions.WithOlderVersions = withOlderVersions
		lstOptions.TimeRef = timeRef
	}
	for content := range clnt.List(ctx, lstOptions) {
		if content.Err != nil {
			errorIf(content.Err.Trace(clnt.GetURL().String()), "Unable to list folder.")
			cErr = exitStatus(globalErrorExitStatus) // Set the exit status.
			continue
		}

		if !recursive && alias+getKey(content) != getStandardizedURL(urlStr) {
			break
		}

		objectsFound = true
		newClnt, perr := newClientFromAlias(alias, content.URL.String())
		if perr != nil {
			errorIf(content.Err.Trace(clnt.GetURL().String()), "Invalid URL")
			continue
		}
		lhold, probeErr := newClnt.GetObjectLegalHold(ctx, content.VersionID)
		if probeErr != nil {
			errorsFound = true
			errorIf(probeErr.Trace(content.URL.Path), "Failed to get legal hold information on `"+content.URL.Path+"`")
		} else {
			if !globalJSON {

				contentURL := filepath.ToSlash(content.URL.Path)
				key := strings.TrimPrefix(contentURL, prefixPath)

				printMsg(legalHoldInfoMessage{
					LegalHold: lhold,
					Status:    "success",
					URLPath:   content.URL.String(),
					Key:       key,
					VersionID: content.VersionID,
				})
			}
		}
	}

	if cErr == nil && !globalJSON {
		switch {
		case errorsFound:
			console.Print(console.Colorize("LegalHoldPartialFailure", fmt.Sprintf("Errors found while getting legal hold status on objects with prefix `%s`. \n", urlStr)))
		case !objectsFound:
			console.Print(console.Colorize("LegalHoldMessageFailure", fmt.Sprintf("No objects/versions found while getting legal hold status with prefix `%s`. \n", urlStr)))
		}
	}
	return cErr
}

// main for legalhold info command.
func mainLegalHoldInfo(cliCtx *cli.Context) error {
	console.SetColor("LegalHoldSuccess", color.New(color.FgGreen, color.Bold))
	console.SetColor("LegalHoldNotSet", color.New(color.FgYellow))
	console.SetColor("LegalHoldOn", color.New(color.FgGreen, color.Bold))
	console.SetColor("LegalHoldOff", color.New(color.FgRed, color.Bold))
	console.SetColor("LegalHoldVersion", color.New(color.FgGreen))
	console.SetColor("LegalHoldPartialFailure", color.New(color.FgRed, color.Bold))
	console.SetColor("LegalHoldMessageFailure", color.New(color.FgYellow))

	targetURL, versionID, timeRef, recursive, withVersions := parseLegalHoldArgs(cliCtx)
	if timeRef.IsZero() && withVersions {
		timeRef = time.Now().UTC()
	}

	ctx, cancelLegalHold := context.WithCancel(globalContext)
	defer cancelLegalHold()

	checkBucketLockSupport(ctx, targetURL)

	return showLegalHoldInfo(ctx, targetURL, versionID, timeRef, withVersions, recursive)
}
