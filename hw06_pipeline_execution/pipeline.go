package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func wrapValue(val interface{}) In {
	ch := make(Bi)
	go func() {
		defer close(ch)
		ch <- val
	}()
	return ch
}

func stageWithDone(in In, done In, stage Stage) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				stageOut := stage(wrapValue(val))
				for v := range stageOut {
					out <- v
				}
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
