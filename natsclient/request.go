package natsclient

func NatsRequest(subject string) {

	nc := NewNatsClient()
	defer nc.Close()

}
