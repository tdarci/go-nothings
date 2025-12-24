package channelsplitter

func Split[T any](inChan chan T, outChanCount int, buffer int) []chan T {
	var out []chan T = make([]chan T, outChanCount)
	for idx := range outChanCount {
		out[idx] = make(chan T, buffer)
	}

	go func() {
		for z := range inChan {
			for chanIdx := range outChanCount {
				out[chanIdx] <- z
			}
		}
		for chanIdx := range outChanCount {
			close(out[chanIdx])
		}
	}()

	return out
}
