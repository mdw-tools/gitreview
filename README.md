# gitreview

WARNING: This README file is built from the source code. To modify its contents:

1. Edit the string values in `config.go`
2. Run `make docs` for the changes to show up in this file

`gitreview` facilitates visual inspection (code review) of git
repositories that meet any of the following criteria:

1. New content was fetched
2. Behind origin/<default-branch>
3. Ahead of origin/<default-branch>
4. Messy (have uncommitted state)
5. Throw errors for the required git operations (listed below)

We use variants of the following commands to ascertain the
status of each repository:

- `git remote`           (shows remote address)
- `git status`           (shows uncommitted files)
- `git fetch`            (finds new commits/tags/branches)
- `git rev-list`         (lists commits behind/ahead-of <default-branch>)
- `git config --get ...` (show config parameters of a repo)

...all of which should be safe enough.

Each repository that meets any criteria above will be
presented for review.

Repositories are gathered recursively from the current working directory.

## Prerequisites:

1. [Git](https://git-scm.com/)
	- The `gitreview` tool will invoke the aforementioned git commands.
1. A Git GUI (recommended: [Sublime Merge](https://www.sublimemerge.com/) or [SourceTree](https://www.sourcetreeapp.com/))
	- The `gitreview` tool invokes your git GUI according to the 'gui' flag (details below).
	- The git GUI used must support invocation from the command line in the form "<git-gui-command> /path/to/repository"
1. A go compiler (v1.16 or higher)
   - download [installer](https://golang.org/dl/)
   - or use brew: "brew install golang"


## Compilation/Installation:

    git clone git@github.com:mdw-tools/gitreview
    cd gitreview
    make install

The above `make install` runs `go install` which installs executables
"in the directory named by the `GOBIN` environment variable, which defaults
to `$GOPATH/bin` or `$HOME/go/bin` if the `GOPATH` environment variable is not set."

See the Makefile for other helpful targets/operations.


## Usage Description

After installation, the `gitreview` executable, when invoked, will commence
running `git fetch` for all repositories it finds based on what is provided
to the `-roots` flag. It uses a fan-out strategy to fetch multiple repositories
concurrently.

As repositories are scanned a brief summary of the status of each is emitted.
After scanning all repositories a report of which repositories require a review.
The `gitreview` program will then halt, waiting for the user to press `<ENTER>`.

Upon receiving `<ENTER>`, `gitreview` will invoke the git GUI program specified
by the `-gui` flag for each of the repositories previously listed. At this point
the `gitreview` tool's execution is concluded.

All that remains is for the user to review each opened repository in the git GUI.
It is the user's responsibility to review each new tag, branch, and commit as well
as to run `git pull` as desired to fully synchronize the local repository with the
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


```
  -default-branch string
    	The default branch to use. Defaults to the value of the
    	GITREVIEWBRANCH environment variable, if declared,
    	otherwise 'main'.
    	--> (default "main")
  -fetch
    	When false, suppress all git fetch operations via --dry-run.
    	Repositories with updates will still be included in the review.
    	--> (default true)
  -gui string
    	The external git GUI application to use for visual reviews.
    	--> (default "smerge")
```
