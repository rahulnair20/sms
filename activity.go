package sms

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"net/http"
	"net/url"
	"strings"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

var log = logger.GetLogger("activity-helloworld")

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

type twilio struct {
	accountSid string
	authToken  string
	urlString  string
	msgData    string
	to         string
	from       string
}

func sendSMS(accountSid string, authToken string, urlString string, smsText string, to string, from string) (string, error) {

	twilio := twilio{accountSid, authToken, urlString + accountSid + "/Messages.json", smsText, to, from}

	msgData := url.Values{}
	msgData.Set("To", twilio.to)
	msgData.Set("From", twilio.from)
	msgData.Set("Body", twilio.msgData)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// Create HTTP request client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", twilio.urlString, &msgDataReader)
	req.SetBasicAuth(twilio.accountSid, twilio.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			//fmt.Println(data["sid"])
			response := []string{"Message sent successfully to", twilio.to}
			return strings.Join(response, ":"), nil
		}
	} else {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		//fmt.Println(resp.Status, "Response", decoder, "Error message:", err)
		panic("test error")
		return "", err
	}
	return "", nil
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	accountSid := context.GetInput("accountSid").(string)
	authToken := context.GetInput("authToken").(string)
	urlString := context.GetInput("urlString").(string)
	smsText := context.GetInput("msgData").(string)
	to := context.GetInput("to").(string)
	from := context.GetInput("from").(string)

	log.Debugf("The Flogo engine says [%s] to [%s]", accountSid+authToken+urlString, smsText+to+from)
	result, err := sendSMS(accountSid, authToken, urlString, smsText, to, from)
	log.Debugf("The Flogo engine says [%s] to [%s]", result, err)
	context.SetOutput("result", result)
	context.SetOutput("error", err)

	return true, err
}
