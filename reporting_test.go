package main

import (
	"testing"
	"time"
)

func reportingTestSetup(input string) (reporter *ReportHandler, tester *TestHandler) {
	var broker DataBroker
	broker.Init()
	reporter = &ReportHandler{&broker}

	tester = new(TestHandler)
	tester.setUpMockInput(input)

	return
}

func reportingTestCleanup(reporter *ReportHandler, tester *TestHandler) {
	reporter.broker.Close()
	tester.cleanUp()
}

func TestReportHandler_GetDateFromUser(t *testing.T) {
	expectedDate, _ := time.Parse("Jan 2 2006", "Apr 27 2019")
	reporter, tester := reportingTestSetup("04/27/2019")
	defer reportingTestCleanup(reporter, tester)

	if date, err := reporter.GetDateFromUser(); err != nil {
		t.Error("Error when parsing user date, got", err)
	} else if !date.Equal(expectedDate) {
		t.Error("Wrong date returned, expected", expectedDate, "got", date)
	}
}

func TestReportHandler_GetTimeFromUser(t *testing.T) {
	expectedTime, _ := time.Parse("3 pm", "7 pm")
	reporter, tester := reportingTestSetup("7 pm")
	defer reportingTestCleanup(reporter, tester)

	if time, err := reporter.GetTimeFromUser(); err != nil {
		t.Error("Error when parsing user time, got", err)
	} else if !time.Equal(expectedTime) {
		t.Error("Wrong date returned, expected", expectedTime, "got", time)
	}
}

func TestReportHandler_GetSpecificShowFromUser(t *testing.T) {
	date, _ := time.Parse("Jan 2 2006 3 pm", "Apr 28 2019 8 pm")
	reporter, tester := reportingTestSetup("4")
	defer reportingTestCleanup(reporter, tester)

	if show, title, err := reporter.GetSpecificShowFromUser(&date); err != nil {
		t.Error("Error while getting show from user, got", err)
	} else if title != "Rogue One: A Star Wars Story" {
		t.Error("Wrong movie title, expected \"Rogue One: A Star Wars Story\", got", title)
	} else if show.ShowID != 28 {
		t.Error("Shhow with wrong ID received, expected 28, got", show.ShowID)
	}
}
