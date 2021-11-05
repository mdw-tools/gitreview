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

1. Edit the string values in config.go
2. Run make docs for the changes to show up in this file

gitreview facilitates visual inspection (code review) of git
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

Installation:

    go get -u github.com/mdwhatcott/gitreview


Skipping Repositories:

If you have repositories in your list that you would rather not review,
you can mark them to be skipped by adding a config variable to the
repository. The following command will produce this result:

    git config --add review.skip true


Omitting Repositories:

If you have repositories in your list that you would still like to audit
but aren't responsible to sign off (it's code from another team), you can
mark them to be omitted from the final report by adding a config variable
to the repository. The following command will produce this result:

    git config --add review.omit true


Specifying the 'default' branch:

This tool assumes that the default branch of all repositories is 'master'.
If a repository uses a non-standard default branch (ie. 'master', 'trunk')
and you want this tool to focus  reviews on commits pushed to that branch
instead, run the following command:

	git config --add review.branch <branch-name>


CLI Flags:
`

var doc = strings.ReplaceAll(rawDoc, "|", "`")
