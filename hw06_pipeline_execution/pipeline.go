package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	return runPipeline(in, done, stages)
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
					for range outCh {
					}
					return
				case v, ok := <-outCh:
					if !ok {
						return
					}
					select {
					case <-done:
						for range outCh {
						}
						return
					case stageIn <- v:
					}
				}
			}
		}(stageIn, outCh)

		outCh = stageOut
	}

	return outCh
}
