package main

import (
	"fmt"
	"log"

	"bbonfim.com/floss/boleto"
	"bbonfim.com/floss/config"
	"bbonfim.com/floss/headers"
)

func main() {

	user := "me"
	srv := config.Init()

	z, err := srv.Users.Messages.List(user).Q("boleto has:attachment").MaxResults(400).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	fmt.Println("Messages: %i", len(z.Messages))
	for _, m := range z.Messages {
		msg, _ := srv.Users.Messages.Get(user, m.Id).Fields().Format("full").Do()
		headers.PrintSubject(msg.Payload.Headers)
		boleto.ExtractCode(msg, srv, user)
	}

}
