package main

import (
	"sort"
	"sync"
)

type Analyzer struct {
	workerCount int
	workerInput chan string
}

func NewAnalyzer(workerCount int, workerInput chan string) *Analyzer {
	return &Analyzer{
		workerCount: workerCount,
		workerInput: workerInput,
	}
}

func (this *Analyzer) AnalyzeAll() (fetches []*GitReport) {
	outputs := this.startWorkers()
	for fetch := range merge(outputs...) {
		fetches = append(fetches, fetch)
	}
	sort.Slice(fetches, func(i, j int) bool {
		return fetches[i].RepoPath < fetches[j].RepoPath
	})
	return fetches
}

func (this *Analyzer) startWorkers() (outputs []chan *GitReport) {
	for x := 0; x < this.workerCount; x++ {
		output := make(chan *GitReport)
		outputs = append(outputs, output)
		go NewWorker(x, this.workerInput, output).Start()
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
