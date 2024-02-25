package broadcasting_to_multiple_goroutine

//Instead of fan-out, we can use a broadcast pattern—one that replicates messages to
//a set of output channels.

//To implement this broadcast utility, we just need to create a list of output channels
//and then use a goroutine that writes every received message to each channel

func CreatAll[K any](n int) []chan K {
	res := make([]chan K, n)
	/*
		在Go中，channel必须通过make函数初始化之后才能使用。
		未初始化的channel默认为nil，尝试向nil channel发送或从中接收数据会导致goroutine永远阻塞
	*/
	for i, _ := range res {
		res[i] = make(chan K)
	}
	return res
}

func CloseAll[K any](channels ...chan K) {
	for _, ch := range channels {
		close(ch)
	}
}

func BroadCast[K any](quit <-chan struct{}, input <-chan K, n int) []chan K {
	outputs := CreatAll[K](n)
	go func() {
		defer CloseAll(outputs...)
		var msg K
		moreData := true
		for moreData {
			select {
			case msg, moreData = <-input:
				if moreData {
					for _, output := range outputs {
						output <- msg
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return outputs
}
