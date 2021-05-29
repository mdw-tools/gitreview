package main

import (
	"sort"
	"sync"
)

type Analyzer struct {
	workerCount   int
	workerInput   chan string
	defaultBranch string
}

func NewAnalyzer(workerCount int, branch string) *Analyzer {
	return &Analyzer{
		workerCount:   workerCount,
		workerInput:   make(chan string),
		defaultBranch: branch,
	}
}

func (this *Analyzer) AnalyzeAll(paths []string) (fetches []*GitReport) {
	go this.loadInputs(paths)
	outputs := this.startWorkers()
	for fetch := range merge(outputs...) {
		fetches = append(fetches, fetch)
	}
	sort.Slice(fetches, func(i, j int) bool {
		return fetches[i].RepoPath < fetches[j].RepoPath
	})
	return fetches
}

func (this *Analyzer) loadInputs(paths []string) {
	for _, path := range paths {
		this.workerInput <- path
	}
	close(this.workerInput)
}

func (this *Analyzer) startWorkers() (outputs []chan *GitReport) {
	for x := 0; x < this.workerCount; x++ {
		output := make(chan *GitReport)
		outputs = append(outputs, output)
		go NewWorker(x, this.defaultBranch, this.workerInput, output).Start()
	}
	return outputs
}

func merge(fannedOut ...chan *GitReport) chan *GitReport {
	var waiter sync.WaitGroup
	waiter.Add(len(fannedOut))

	fannedIn := make(chan *GitReport)

	output := func(c <-chan *GitReport) {
		for n := range c {
			fannedIn <- n
		}
		waiter.Done()
	}

	for _, c := range fannedOut {
		go output(c)
	}

	go func() {
		waiter.Wait()
		close(fannedIn)
	}()

	return fannedIn
}
