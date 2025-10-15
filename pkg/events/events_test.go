package events_test

import (
	"testing"
	"time"

	"signinapc.app/api/pkg/events"
)

func TestSignedInBody(t *testing.T) {
	eventTime, err := time.Parse(time.Layout, "10/11 04:10:35AM '28 -0400")

	if err != nil {
		t.Error("Failed to parse time. ", err.Error())
	}

	expected := `{"name":"Test Person","email":"test.email@gmail.com","time":"2028-10-11T04:10:35-04:00"}`

	value := events.MakeSignedIn(
		events.MakeSignedInData("Test Person", "test.email@gmail.com", eventTime),
		1,
	)

	result := value.Body()

	if result != expected {
		t.Errorf("Expected %s to be %s.", result, expected)
	}
}

func TestSignedInShow(t *testing.T) {
	eventTime, err := time.Parse(time.Layout, "06/05 10:05:44PM '24 +0300")

	if err != nil {
		t.Errorf("Failed to parse time. %s", err.Error())
	}

	expected := "SignedIn(Data = SignedInData(Name = Test Person, Email = test.email@gmail.com, Time = 06/05 10:05:44PM '24 +0300), Seqeuence = 1)"

	value := events.MakeSignedIn(events.MakeSignedInData("Test Person", "test.email@gmail.com", eventTime), 1)

	result := value.Show()

	if result != expected {
		t.Errorf("Expected %s to be %s.", result, expected)
	}
}

func TestRegexIDsToID(t *testing.T) {
	ids, err := events.MakeRegexIDs()

	if err != nil {
		t.Errorf("Failed to create regex IDs instance: %s", err.Error())
	}

	value := "Test User"
	expected := "test-user"

	result := ids.ToID(value)

	if result != expected {
		t.Errorf("Expected %s to be %s", result, expected)
	}
}
