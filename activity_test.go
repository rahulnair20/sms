package sms

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs

	tc.SetInput("accountSid", "AC8112820c2969c0b9ba6abac8ee6a4062")
	tc.SetInput("authToken", "f2542c88dbfb58c494e642bf10af4140")
	tc.SetInput("urlString", "https://api.twilio.com/2010-04-01/Accounts/")
	tc.SetInput("msgData", "Hello from flogo")
	tc.SetInput("to", "+919177623444")
	tc.SetInput("from", "+14437433811")

	act.Eval(tc)

	//check result attr

	result := tc.GetOutput("result")
	assert.Equal(t, result, "Message sent successfully to:+919177623444")

}
