package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Lookup the Page on Warwick DC for CV357TT
		resp, err := http.Get("https://estates7.warwickdc.gov.uk/PropertyPortal/Property/Recycling/10003790863")
		if err != nil {
			log.Fatalln(err)
		}

		// Load in Environment Variables
		godotenv.Load("local.env")

		// Read the response body
		page, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		sb := string(page)

		// Start with Today's Date
		startDate := time.Now()

		// Use this to test other dates
		// startDate := time.Date(2022, time.January, 15, 0, 0, 0, 0, time.UTC)

		// loop forward form today for 7 days in to the next week to find a collection Date
		for i := 0; i < 7; i++ {

			date := startDate.AddDate(0, 0, i)
			dateString := date.Format(os.Getenv("DATE_FORMAT_DDMMYYYY"))

			// Is there a Collection on this day?
			if strings.Contains(sb, dateString) {

				// We have a date, find what type of collection it is
				collection := getCollectionForDate(sb, dateString)
				collection = cleanString(collection)

				// Send details of which collection was found
				sendEmail(collection)

				fmt.Fprintf(w, "We've found a collection date for %s, it's %s this week\n", dateString, collection)
			}
		}
	})

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
}

/// function to Parse the Collection Type based on the date found
func getCollectionForDate(searchText string, date string) string {

	// Where is the date in the HTML source?
	datePos := strings.Index(searchText, date)

	// Find and parse the preceeding Strong Tag with the Collection Type
	closingTagPos := strings.LastIndex(searchText[:datePos], "</strong>")
	openingTagPos := strings.LastIndex(searchText[:closingTagPos], "<strong>")
	openingTagPos += len("<strong>")

	res := searchText[openingTagPos:closingTagPos]
	return (res)
}

// Cleans line breaks and carriage returns etc
func cleanString(str string) string {

	// remove <br /> tags, trim collection and remove any spaces
	str = strings.Replace(str, "<br />", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.TrimSpace(str)

	re := regexp.MustCompile(`\r?\n`)
	str = re.ReplaceAllString(str, " ")

	return str
}

// Emails the results to anyone who has subscribed
func sendEmail(collectionType string) {

	recipients := strings.Split(os.Getenv("EMAIL_RECIPIENTS"), ";")

	// loop through the recipients and send an email for each
	for _, recipient := range recipients {
		fmt.Println(recipient)

		from := mail.NewEmail("Jason", "samjas73@gmail.com")
		subject := "It's " + collectionType + " this week"
		to := mail.NewEmail("Jason", recipient)
		plainTextContent := subject
		htmlContent := subject
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		response, err := client.Send(message)

		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response.StatusCode)
		}

		fmt.Println("Email Sent!")
	}
}
