package boleto

import (
	"encoding/base64"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

//Prints the "Boleto" code that can be used to pay out the invoice
func ExtractCode(msg *gmail.Message, srv *gmail.Service, user string) {
	var boletoFound = false
	for _, part := range msg.Payload.Parts {
		if strings.HasSuffix(part.Filename, ".pdf") {
			err := writePdfToDisk(msg, srv, user, part)

			entries, err := convertingPdfToTextFiles(part, err)
			if err != nil {
				continue
			}
			boletoFound = extractBoletoCodeFromText(entries, boletoFound)

			os.RemoveAll("/tmp/pepper")
		}
		if !boletoFound {
			fmt.Println("Boleto not found")
		}
	}
}

func extractBoletoCodeFromText(entries []os.FileInfo, boletoFound bool) bool {
	for _, entry := range entries {
		content, err := ioutil.ReadFile("/tmp/noSurprise/out/" + entry.Name())
		if err != nil {
			log.Fatal(err)
		}

		re := regexp.MustCompile(`\d{5}\.\d{5} \d{5}\.\d{6} \d{5}\.\d{6} \d \d{14}`)
		text := string(content)
		code := re.FindString(text)
		if len(code) > 0 {
			fmt.Printf("found code: %s\n", code)
			boletoFound = true
		}
	}
	return boletoFound
}

func convertingPdfToTextFiles(part *gmail.MessagePart, err error) ([]os.FileInfo, error) {
	err = api.ExtractContentFile("/tmp/noSurprise/"+part.PartId+".pdf", "/tmp/noSurprise/out", nil, nil)
	if err != nil {
		log.Printf("Unable to extract content from PDF file: %v", err)
		return nil, err
	}

	entries, err := ioutil.ReadDir("/tmp/noSurprise/out")
	if err != nil {
		log.Panicf("failed reading directory: %s", err)
	}
	return entries, err
}

func writePdfToDisk(msg *gmail.Message, srv *gmail.Service, user string, part *gmail.MessagePart) error {
	fmt.Printf("processing file: %s\n", part.Filename)
	err := os.MkdirAll("/tmp/noSurprise/out", 0755)
	if err != nil {
		log.Fatalf("Unable to create directory: %v", err)
	}
	attach, _ := srv.Users.Messages.Attachments.Get(user, msg.Id, part.Body.AttachmentId).Do()
	decoded, err := base64.URLEncoding.DecodeString(attach.Data)

	f, err := os.Create("/tmp/noSurprise/" + part.PartId + ".pdf")
	if err != nil {
		log.Fatalf("Unable to create file: %v", err)
	}

	defer f.Close()

	n2, err := f.Write(decoded)
	if err != nil {
		log.Fatalf("Unable to write file: %v", err)
	}
	fmt.Printf("wrote %d bytes\n", n2)
	return err
}
