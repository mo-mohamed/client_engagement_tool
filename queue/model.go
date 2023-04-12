package queue

type SendRequest struct {

	// String value of the Queue url intended to use.
	QueueUrl string

	// String value of the message body.
	Body string

	// List of attributes associated with the message.
	// A single attributes holds 3 values:
	//   - Key: describes the value
	//	 - Value: the actual value
	//	 - Type: the value's data type
	Attributes []Attribute
}

type Attribute struct {
	// String value of the Queue url intended to use.
	Key string

	// String value of the Queue url intended to use.
	Value string

	// String value of the Queue url intended to use.
	Type string
}

type Message struct {
	// Unique message ID. Populated with the brooker's message id.
	ID string
	// Receipt Handler. Populated with the brooker's message id.
	ReceiptHandler string
	// Message body.
	Body string
	// List of attributes associated with the queue message.
	Attributes map[string]Attribute
}
