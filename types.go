package main

import "github.com/streadway/amqp"

// Params represents the valid parameter options for the rabbitmq plugin.
type Params struct {
	Connection Connection
	Exchange   string
	Key        string
	Mandatory  bool
	//Immediate bool
	WaitConfirm bool
	Publishing  Publishing
	Template    string
}

// Connection represents amqp connection options.
type Connection struct {
	Host     string
	Username string
	Password string
}

type Publishing struct {
	// Application or exchange specific fields,
	// the headers exchange will inspect this field.
	Headers amqp.Table

	// Properties
	ContentType     string // MIME content type
	ContentEncoding string // MIME content encoding
	DeliveryMode    uint8  // Transient (0 or 1) or Persistent (2)
	Priority        uint8  // 0 to 9
	CorrelationId   string // correlation identifier
	ReplyTo         string // address to to reply to (ex: RPC)
	Expiration      string // message expiration spec
	MessageId       string // message identifier
	Type            string // message type name
	UserId          string // creating user id - ex: "guest"
	AppId           string // creating application id
}
