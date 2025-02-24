package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func stageWithDone(in In, done In, stage Stage) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		stageOut := stage(in)
		for {
			select {
			case <-done:
				return
			case val, ok := <-stageOut:
				if !ok {
					return
				}
				out <- val
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stageWithDone(in, done, stage)
	}
	return in
}
