package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		return nil
	}
	/*
		The main principe of controlling stages is using sentinel (intermediate) channels.
		It helps to make controllable (via special "done" channel) channels returned by stages.
	*/
	for _, stage := range stages {
		// Create sentinel channel.
		sentinelChannel := make(chan interface{})
		// Start wrapping goroutine.
		go func(ch In) {
			defer close(sentinelChannel)
			for {
				// Try to:
				// 1) receive data from input channel and put it in sentinel channel;
				// 2) or receive signal from "done" channel.
				// Return in case of closed input channel or signal received from "done".
				select {
				case v, ok := <-ch:
					if !ok {
						return
					}

					select {
					case sentinelChannel <- v:
						continue
					case <-done:
						return
					}
				case <-done:
					return
				}
			}
		}(in)
		// Sentinel channel turns into input channel for stage,
		// and output channel of stage becomes new input for next wrapping goroutine.
		in = stage(sentinelChannel)
	}
	return in
}
