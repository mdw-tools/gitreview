// gitreview scans entries found in the `CDPATH` environment variable
// looking for git repositories that are messy or behind and opens
// a git GUI (Sublime Merge by default) for each to facilitate a review.
// It only runs `git status` and `git fetch`, which should be safe.
// After all reviews are complete it prints (to `stdout`) a concatenated
// report of all `git fetch` output for repos that were behind their origin.
package main

func main() {
	config := ReadConfig()
	reviewer := NewGitReviewer(config.GitRoots, config.GitGUI)
	reviewer.FetchAllRepositories()
	reviewer.ReviewAllNotableRepositories()
	reviewer.PrintCodeReviewLogEntry()
}