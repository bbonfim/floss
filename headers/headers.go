package headers

import (
	"fmt"
	"strings"

	"google.golang.org/api/gmail/v1"
)

// PrintSubject Prints the subject line of a message
func PrintSubject(headers []*gmail.MessagePartHeader) {

	for _, header := range headers {
		if strings.Contains(header.Name, "Subject") {
			fmt.Printf("Subject: %s\n", header)
		}
	}
}
