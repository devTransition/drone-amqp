package main

import (
	"bytes"
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
	"github.com/drone/drone-go/template"
	"github.com/streadway/amqp"
	"log"
	"os"
)

var (
	buildCommit string
)

func main() {
	fmt.Printf("Drone RabbitMQ Plugin built from %s\n", buildCommit)

	system := drone.System{}
	repo := drone.Repo{}
	build := drone.Build{}
	vargs := Params{}

	plugin.Param("system", &system)
	plugin.Param("repo", &repo)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	// Creates the payload, by default the payload
	// is the build details in json format, but a custom
	// template may also be used.

	var buf bytes.Buffer
	err := template.Write(&buf, vargs.Template, &drone.Payload{
		Build:  &build,
		Repo:   &repo,
		System: &system,
	})

	if err != nil {
		log.Printf("Error: Failed to execute the content template. %s\n", err)
		os.Exit(1)
	}

	log.Printf("Data: %q, Template: %s, Payload: %s", vargs, vargs.Template, buf.String())
	msg := amqp.Publishing{
		Headers: vargs.Publishing.Headers,
		// Properties
		ContentType:     vargs.Publishing.ContentType,     // MIME content type
		ContentEncoding: vargs.Publishing.ContentEncoding, // MIME content encoding
		DeliveryMode:    vargs.Publishing.DeliveryMode,    // Transient (0 or 1) or Persistent (2)
		Priority:        vargs.Publishing.Priority,        // 0 to 9
		CorrelationId:   vargs.Publishing.CorrelationId,   // correlation identifier
		ReplyTo:         vargs.Publishing.ReplyTo,         // address to to reply to (ex: RPC)
		Expiration:      vargs.Publishing.Expiration,      // message expiration spec
		MessageId:       vargs.Publishing.MessageId,       // message identifier
		Type:            vargs.Publishing.Type,            // message type name
		UserId:          vargs.Publishing.UserId,          // creating user id - ex: "guest"
		AppId:           vargs.Publishing.AppId,           // creating application id

		// The application specific payload of the message
		Body: buf.Bytes(),
	}

	err = PublishMesssage(
		vargs.Connection.Host,
		vargs.Connection.Username,
		vargs.Connection.Password,
		vargs.Exchange,
		vargs.Key,
		vargs.Mandatory,
		vargs.WaitConfirm,
		&msg,
	)

	if err != nil {
		log.Printf("Error: PublishMesssage %s\n", err)
		os.Exit(1)
	}

}
