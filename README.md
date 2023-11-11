<!--- goToMD::Auto:: See github.com/dancsecs/goToMD ** DO NOT MODIFY ** -->

# Package goToMd

<!--- goToMD::Bgn::doc::./package -->
```go
package main
```

The goToMD utility provides for the maintenance of github README.MD style
pages by permitting go files, go documentation and go test output to be
included by reference into the github README.md file directly from the Go
code permitting program documentation to be maintained in one place (the Go
code.)

It can use a template file (```*.md.gtm```) or can maintain a ```*.md``` file
in place.

Usage of goToMD [-c | -r] [-fvl] [-p perm] [-o outDir] [-U file] [-u uint] path [path...]

The flags are:

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
    -U  string
        Collect cpu profile data into named file.
    -u  uint
        Number of iterations to run when collecting cpu profile information.
    -v
        Provide more information when processing.

Directives are placed into the ```*.md.gtm``` file (or directly into the
```*.md``` document if the replace in place option is given.  These directves
are emebbed into HTML style comments.

```html
<!--- goToMD::ACTION::PARAMETERS -->
```

where ACTION can be one of the following:

- goToMD::doc::./relativeDirectory/goObject

        Run the go doc command on the object listed from the directory
        specified.  The PARAMETER is foramtted with the relative directory up
        to the last directory seperater before the end of the string and the
        desired object.  A special object package returns the package
        comments.

- goToMD::dcls::./relativeDirectory/declaredObject ListOfDeclaredGoObjects

        Pull out the declaration for the object and include as a single line
        regardless of how declaed in the go code.  The Parameter is a list of
        go functions, methods and constants (more object coming) to be included
        in a go code block. No comments are included.

- goToMD::dcln::./relativeDirectory/delaredObject ListOfDeclaredGoObjects

         Pull the declaration and include exactly as declared in the go
         source including leading comments.

- goToMD::dcl::./relativeDirectory/delaredObject ListOfDeclaredGoObjects

         Pull the declaration and include exactly as declared in the go
         source.  No Comments are included.

- goToMD::tst::goTest::./relativeDirectory/testName

         Run go test with the tests listed (or package) to run all tests and
         included the output.

- goToMD::file::./relativeDirectory/fName

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
<!--- goToMD::End::doc::./package -->
