package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var wg sync.WaitGroup
	out := make(Bi)

	multiplex := func(c <-chan interface{}) {
		wg.Add(1)
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case out <- i:
			}
		}
	}

	for i := 0; i < 2; i++ {
		go multiplex(runPipeline(in, done, stages))
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func runPipeline(in In, done In, stages []Stage) Out {
	if len(stages) == 0 {
		return in
	}

	outCh := in

	for _, stage := range stages {
		stageIn := make(Bi)
		stageOut := stage(stageIn)

		go func(stageIn Bi, outCh In) {
			defer close(stageIn)
			for {
				select {
				case <-done:
					return
				case v, ok := <-outCh:
					if !ok {
						return
					}
					stageIn <- v
				}
			}
		}(stageIn, outCh)

		outCh = stageOut
	}

	return outCh
}
