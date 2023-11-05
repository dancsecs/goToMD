//nolint:goDot // Ok.
/*
The goToMD utility provides for the maintenance of github README.MD style
pages by permitting go files, go documentation and go test output to be
included by reference into the github README.md file directly from the Go
code permitting program documentation to be maintained in one place (the Go
code.)

It can use a template file (```*.md.gtm") or can maintain a ```*.md``` file
in place.

Usage: goToMD [flags] [path ...]

The flags are:

Usage of goToMD [-c | -r] [-fvl] [-p perm] [-o outDir] file|dir [file|dir...]

    -c
        Reverse operation and remove generated markdown (Cannot be used
        with the -r option). Files with a .md extension are expected with
        a an .md.gtm file being produced.
    -f
        Do not confirm overwrite of destination.
    -l
        Display license before program exits.
    -o string
        Direct all output to the specified directory. (default ".")
    -p int
        Permissions to use when creating new file (can only set RW
        bits). (default 420)
    -r
        Replace the *.MD in place (Cannot be used with the -c flag).
    -v
        Provide more information woth respect to processing.

Directives are placed into the ```*.md.gtm``` file (or directly into the
```*.md``` document if the replace in place option is given.  These directves
are emebbed into HTML style comments.

```html
<!--- goToMD::ACTION::PARAMETERS -->
```

where ACTION can be one of the following:

 const szDocPrefix = "doc::"
        Run the go doc command on the object listed from the directory
        specified.  The PARAMETER is foramtted with the relative directory up
        to the last directory seperater before the end of the string and the
        desired object.  A special object package returns the package
        comments.

 const szDclsPrefix = "dcls::"
        Pull out the declaration for the object and include as a single line
        regardless of how declaed in the go code.  The Parameter is a list of
        go functions, methods and constants (more object coming) to be included
        in a go code block. No comments are included.

 const szDclPrefix = "dcl::"
         Pull the declaration and include exactly as declared in the go
         source.  No Comments are included.

 const szTstPrefix = "tst::"
         Run go test with the tests listed (or package) to run all tests and
         included the output.

 const szFilePrefix = "file::"
         Include the specified file in a code block.


When expanded in the target file the content will be framed by similar
comments prefixed with Bgn and end as:

const szTestBgnPrefix = szTestPrefix + "Bgn::"
const szTestEndPrefix = szTestPrefix + "End::"

A header prefixed with

const szAutoPrefix = szTestPrefix + "Auto::"

and a blank line following will be inserted into the output file.  If the
action is not "replace in place" then an addition ** DO NOT MODIFY **
warning is included.
*/
package main

import (
	"fmt"
	"os"
)

const license = `
Golang To Github Markdown: goToMD.
Copyright (C) 2023  Leslie Dancsecs

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
`

func main() {
	var err error
	var origWd string
	var filesToProcess []string

	// Restore original working directory on exit.
	origWd, err = os.Getwd()
	if err == nil {
		defer func() {
			_ = os.Chdir(origWd)
		}()
	}

	processArgs()

	filesToProcess, err = getFilesToProcess()

	for i, mi := 0, len(filesToProcess); i < mi && err == nil; i++ {
		err = os.Chdir(origWd)
		if err == nil {
			switch {
			case cleanOnly:
				err = cleanMD(filesToProcess[i])
			case replace:
				err = replaceMDInPlace(filesToProcess[i])
			default:
				err = expandMD(filesToProcess[i])
			}
		}
	}

	if showLicense {
		fmt.Print(license)
	}

	if err != nil {
		panic(err)
	}
}
