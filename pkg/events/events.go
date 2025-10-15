package events

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type IDs interface {
	ToID(value string) string
	Combine(ids []string) string
}

type regexIDs struct {
	matchWS regexp.Regexp
}

func (r regexIDs) ToID(value string) string {
	removedSpaces := r.matchWS.ReplaceAllString(value, "-")
	lowercase := strings.ToLower(removedSpaces)

	return lowercase
}

func (regexIDs) Combine(ids []string) string {
	return strings.Join(ids, "///")
}

func MakeRegexIDs() (IDs, error) {
	matchWS, err := regexp.Compile(`\s+`)

	if err != nil {
		return nil, errors.New("Failed to compile whitespace regex: " + err.Error())
	}

	ids := regexIDs{
		matchWS: *matchWS,
	}

	return ids, nil
}

type EventMeta struct {
	Stream      string
	Type        string
	SequenceNum int64
}

type Event interface {
	Meta(ids IDs) EventMeta
	Body() string
}

func MemberStreamID(name string, email string, ids IDs) string {
	nameID := ids.ToID(name)
	emailID := ids.ToID(email)

	return ids.Combine([]string{nameID, emailID})
}

type SignedInData struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Time  time.Time `json:"time"`
}

func (s SignedInData) Show() string {
	return fmt.Sprintf("SignedInData(Name = %s, Email = %s, Time = %s)", s.Name, s.Email, s.Time.Format(time.Layout))
}

func MakeSignedInData(name string, email string, time time.Time) SignedInData {
	return SignedInData{
		Name:  name,
		Email: email,
		Time:  time,
	}
}

type SignedIn struct {
	Data     SignedInData
	Sequence int64
}

const SIGNED_IN_EVENT_TYPE = "SignedIn"

func MakeSignedIn(data SignedInData, sequence int64) SignedIn {
	return SignedIn{
		Data:     data,
		Sequence: sequence,
	}
}

func (s SignedIn) Show() string {
	return fmt.Sprintf("SignedIn(Data = %s, Seqeuence = %d)", s.Data.Show(), s.Sequence)
}

func (s SignedIn) Meta(ids IDs) EventMeta {
	return EventMeta{
		Stream:      MemberStreamID(s.Data.Name, s.Data.Email, ids),
		Type:        SIGNED_IN_EVENT_TYPE,
		SequenceNum: s.Sequence,
	}
}

func (s SignedIn) Body() string {
	encoded, err := json.Marshal(s.Data)

	if err != nil {
		fmt.Println("Failed to encode ", s.Show(), ". Using empty JSON for output.")
		encoded = []byte("{}")
	}

	var compacted bytes.Buffer
	json.Compact(&compacted, encoded)

	return compacted.String()
}
