# gitreview

gitreview facilitates visual inspection (code review) of git
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


```
  -fetch
    	When false, suppress all git fetch operations via --dry-run.
    	Repositories with updates will still be included in the review.
    	--> (default true)
  -gui string
    	The external git GUI application to use for visual reviews.
    	--> (default "smerge")
  -roots string
    	The name of the environment variable containing colon-separated
    	path values to scan for any git repositories contained therein.
    	Scanning is NOT recursive.
    	NOTE: this flag will be ignored in the case that non-flag command
    	line arguments representing paths to git repositories are provided.
    	--> (default "GITREVIEWPATH")
```
