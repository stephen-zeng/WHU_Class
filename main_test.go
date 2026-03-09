package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	ics "github.com/arran4/golang-ical"
)

func Test(t *testing.T) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	event := cal.AddEvent(fmt.Sprintln("id@domain"))
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(time.Now())
	event.SetEndAt(time.Now())
	event.SetSummary("Summary")
	event.SetLocation("Address")
	event.SetDescription("Description")
	event.SetURL("https://URL/")
	event.AddRrule(fmt.Sprintf("FREQ=YEARLY;BYMONTH=%d;BYMONTHDAY=%d", time.Now().Month(), time.Now().Day()))
	event.SetOrganizer("sender@domain", ics.WithCN("This Machine"))
	event.AddAttendee("reciever or participant", ics.CalendarUserTypeIndividual, ics.ParticipationStatusNeedsAction, ics.ParticipationRoleReqParticipant, ics.WithRSVP(true))
	fmt.Println(cal.Serialize())
}

func TestGetKBListSafe_HTMLEntityDecoding(t *testing.T) {
	// Track what query parameters the server receives
	var receivedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(KBListResponse{})
	}))
	defer server.Close()

	// Simulate a cURL command with &amp; (HTML entity) in the URL
	curlWithAmp := fmt.Sprintf("curl '%s/test?gnmkdm=N2151&amp;layout=default' -H 'Accept: application/json'", server.URL)

	resp, err := getKBListSafe(curlWithAmp)
	if err != nil {
		t.Fatalf("getKBListSafe returned unexpected error: %v", err)
	}

	// Verify the response was parsed (empty is fine, just no error)
	_ = resp

	// Verify the server received properly decoded query parameters
	if strings.Contains(receivedQuery, "amp;") {
		t.Errorf("URL was not properly decoded: query contains 'amp;': %s", receivedQuery)
	}
	if !strings.Contains(receivedQuery, "gnmkdm=N2151") {
		t.Errorf("expected query to contain 'gnmkdm=N2151', got: %s", receivedQuery)
	}
	if !strings.Contains(receivedQuery, "layout=default") {
		t.Errorf("expected query to contain 'layout=default', got: %s", receivedQuery)
	}
}
