package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	GitDefaultBranch   string
	GitFetch           bool
	GitRepositoryPaths []string
	GitRepositoryRoots []string
	GitGUILauncher     string
}

func ReadConfig() *Config {
	log.SetFlags(log.Ltime | log.Lshortfile)

	config := new(Config)

	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, doc)
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "```")
		flag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr, "```")
	}

	flag.StringVar(&config.GitDefaultBranch,
		"default-branch", or(os.Getenv("GITREVIEWBRANCH"), "main"), ""+
			"The default branch to use. Defaults to the value of the\n"+
			"GITREVIEWBRANCH environment variable, if declared,\n"+
			"otherwise 'main'.\n"+
			"-->",
	)

	flag.StringVar(&config.GitGUILauncher,
		"gui", "smerge", ""+
			"The external git GUI application to use for visual reviews."+"\n"+
			"-->",
	)

	flag.BoolVar(&config.GitFetch,
		"fetch", true, ""+
			"When false, suppress all git fetch operations via --dry-run."+"\n"+
			"Repositories with updates will still be included in the review."+"\n"+
			"-->",
	)

	gitRoots := flag.String(
		"roots", "GITREVIEWPATH", ""+
			"The name of the environment variable containing colon-separated"+"\n"+
			"path values to scan for any git repositories contained therein."+"\n"+
			"Scanning is NOT recursive."+"\n"+
			"NOTE: this flag will be ignored in the case that non-flag command"+"\n"+
			"line arguments representing paths to git repositories are provided."+"\n"+
			"-->",
	)

	flag.Parse()

	config.GitRepositoryPaths = flag.Args()
	if len(config.GitRepositoryPaths) == 0 {
		config.GitRepositoryRoots = strings.Split(os.Getenv(*gitRoots), ":")
	}
	if !config.GitFetch {
		log.Println("Running git fetch with --dry-run (updated repositories will not be reviewed).")
		gitFetchCommand += " --dry-run"
	}
	return config
}

func or(a string, b string) string {
	if len(a) > 0 {
		return a
	}
	return b
}

const rawDoc = `# gitreview

WARNING: This README file is built from the source code. To modify its contents:

1. Edit the string values in |config.go|
2. Run |make docs| for the changes to show up in this file

|gitreview| facilitates visual inspection (code review) of git
repositories that meet any of the following criteria:

1. New content was fetched
2. Behind origin/<default-branch>
3. Ahead of origin/<default-branch>
4. Messy (have uncommitted state)
5. Throw errors for the required git operations (listed below)

We use variants of the following commands to ascertain the
status of each repository:

- |git remote|           (shows remote address)
- |git status|           (shows uncommitted files)
- |git fetch|            (finds new commits/tags/branches)
- |git rev-list|         (lists commits behind/ahead-of <default-branch>)
- |git config --get ...| (show config parameters of a repo)

...all of which should be safe enough.

Each repository that meets any criteria above will be
presented for review.

Repositories are identified for consideration from path values
supplied as non-flag command line arguments or via the roots
flag (see details below).


## Prerequisites:

1. [Git](https://git-scm.com/)
	- The |gitreview| tool will invoke the aforementioned git commands.
1. A Git GUI (recommended: [Sublime Merge](https://www.sublimemerge.com/) or [SourceTree](https://www.sourcetreeapp.com/))
	- The |gitreview| tool invokes your git GUI according to the 'gui' flag (details below).
	- The git GUI used must support invocation from the command line in the form "<git-gui-command> /path/to/repository"
1. A go compiler (v1.16 or higher)
   - download [installer](https://golang.org/dl/)
   - or use brew: "brew install golang"


## Compilation/Installation:

    git clone git@github.com:mdwhatcott/gitreview
    cd gitreview
    make install

The above |make install| runs |go install| which installs executables
"in the directory named by the |GOBIN| environment variable, which defaults
to |$GOPATH/bin| or |$HOME/go/bin| if the |GOPATH| environment variable is not set."

See the Makefile for other helpful targets/operations.


## Usage Description

After installation, the |gitreview| executable, when invoked, will commence
running |git fetch| for all repositories it finds based on what is provided
to the |-roots| flag. It uses a fan-out strategy to fetch multiple repositories
concurrently.

As repositories are scanned a brief summary of the status of each is emitted.
After scanning all repositories a report of which repositories require a review.
The |gitreview| program will then halt, waiting for the user to press |<ENTER>|.

Upon receiving |<ENTER>|, |gitreview| will invoke the git GUI program specified
by the |-gui| flag for each of the repositories previously listed. At this point
the |gitreview| tool's execution is concluded.

All that remains is for the user to review each opened repository in the git GUI.
It is the user's responsibility to review each new tag, branch, and commit as well
as to run |git pull| as desired to fully synchronize the local repository with the
remote.


## Pre-repository Settings

Skipping Repositories:

If you have repositories in your list that you would rather not review,
you can mark them to be skipped by adding a config variable to the
repository. The following command will produce this result:

    git config --add review.skip true


Specifying the 'default' branch:

This tool assumes that the default branch of all repositories is 'main'.
If a repository uses a non-standard default branch (ie. 'master', 'trunk')
and you want this tool to focus  reviews on commits pushed to that branch
instead, run the following command:

	git config --add review.branch <branch-name>


CLI Flags:
`

var doc = strings.ReplaceAll(rawDoc, "|", "`")
