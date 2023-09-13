package workers

func waitForRoutingWorkers(pauseChan chan bool, continueChan chan bool, numRoutingGoroutines int) {
	for i := 0; i < numRoutingGoroutines; i++ {
		pauseChan <- true
	}
	for i := 0; i < numRoutingGoroutines; i++ {
		<-continueChan
	}
}
