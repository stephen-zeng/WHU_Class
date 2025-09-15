package main

import (
	"fmt"
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
