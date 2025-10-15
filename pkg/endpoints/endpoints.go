package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"time"

	"signinapc.app/api/pkg/events"
	"signinapc.app/api/pkg/services"
)

type Endpoint interface {
	Path() string
	Handle(w http.ResponseWriter, r *http.Request)
}

func OK(body any) func(http.ResponseWriter, *http.Request) {
	encoded, err := json.Marshal(body)

	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(encoded)

		if err != nil {
			log.Println("Failed to send response:", err.Error())
		}
	}
}

type memberSignIn struct {
	events services.EventService
	ids    events.IDs
}

type memberSignInRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Agree bool   `json:"agreeToMemberAgreement"`
}

func (m memberSignInRequest) Validate() error {
	errs := []error{}

	if m.Name == "" {
		errs = append(errs, errors.New("field 'name' must be non-empty"))
	}

	if m.Email == "" {
		errs = append(errs, errors.New("field 'email' must be non-empty"))
	}

	if !m.Agree {
		errs = append(errs, errors.New("member did not agree to the member agreement"))
	}

	return errors.Join(errs...)
}

type memberSignInResponse struct{}

func (memberSignIn) Path() string {
	return "/member/signin"
}

func (m memberSignIn) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	requestBody, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Unexpected error with reading body of request.", err.Error())
		http.Error(w, "Failed to read body of request.", http.StatusInternalServerError)
		return
	}

	var request memberSignInRequest
	err = json.Unmarshal(requestBody, &request)

	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request:\n%s", err.Error()), http.StatusBadRequest)
		return
	}

	err = request.Validate()

	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request:\n%s", err.Error()), http.StatusBadRequest)
		return
	}

	stream := events.MemberStreamID(request.Name, request.Email, m.ids)

	sequence, err := m.events.NextSequence(stream)

	if err != nil {
		log.Println("Failed retrieve sequence number:", err.Error())
		http.Error(w, "Failed to query database for necessary information.", http.StatusBadRequest)
		return
	}

	eventTime := time.Now().Round(time.Minute)

	eventData := events.MakeSignedInData(request.Name, request.Email, eventTime)

	persisted, err := m.events.Persist(events.MakeSignedIn(eventData, sequence))

	if err != nil {
		log.Println("Failed to persist event:", err.Error())
		return
	}

	if !persisted {
		http.Error(w, "A newer change for this member has been made! Please reload and try again.", http.StatusBadRequest)
		return
	}

	OK(memberSignInResponse{})(w, r)
}

func MemberSignIn(events services.EventService, ids events.IDs) Endpoint {
	return memberSignIn{
		events,
		ids,
	}
}

type memberAgreed struct {
	events services.EventService
	ids    events.IDs
}

type memberAgreedRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (m memberAgreedRequest) Validate() error {
	errs := []error{}

	if m.Name == "" {
		errs = append(errs, errors.New("field 'name' must be non-empty"))
	}

	if m.Email == "" {
		errs = append(errs, errors.New("field 'email' must be non-empty"))
	}

	return errors.Join(errs...)
}

func (memberAgreed) Path() string {
	return "/member/agreed"
}

func (m memberAgreed) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	requestBody, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Unexpected error with reading body of request.", err.Error())
		http.Error(w, "Failed to read body of request.", http.StatusInternalServerError)
		return
	}

	var request memberAgreedRequest
	err = json.Unmarshal(requestBody, &request)

	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request:\n%s", err.Error()), http.StatusBadRequest)
		return
	}

	err = request.Validate()

	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request:\n%s", err.Error()), http.StatusBadRequest)
		return
	}

	stream := events.MemberStreamID(request.Name, request.Email, m.ids)

	records, err := m.events.ReadAll(stream)

	if err != nil {
		log.Println("Failed to query database:", err.Error())
		http.Error(w, "Failed to query database.", http.StatusInternalServerError)
		return
	}

	signedIn := slices.ContainsFunc(
		records,
		func(e events.Event) bool {
			return e.Meta(m.ids).Type == events.SIGNED_IN_EVENT_TYPE
		},
	)

	OK(signedIn)(w, r)
}

func MemberAgreed(events services.EventService, ids events.IDs) Endpoint {
	return memberAgreed{
		events,
		ids,
	}
}
