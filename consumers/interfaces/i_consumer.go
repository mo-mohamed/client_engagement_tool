package interfaces

/*
Consumer interface, should be implemented by all consumers.

Run() method is responsible for starting the consumer, it expectes the consumer struct to be already configured.

Stop() gracfully shuts down the consumer.
*/
type IConsumer interface {
	// Starts the consumer
	Run()
	// Stops the consumer
	Stop()
}
