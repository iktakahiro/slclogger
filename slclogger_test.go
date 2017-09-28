package slclogger

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestIsValidWebHookURL(t *testing.T) {
	isValid := isValidWebHookURL("https://hooks.slack.com/services/T")
	if !isValid {
		t.Errorf("isValid should be '%v' actual: %v", true, isValid)
	}

	isValid = isValidWebHookURL("invalid.slack.com/services/T")
	if isValid {
		t.Errorf("isValid should be '%v' actual: %v", false, isValid)
	}

	isValid = isValidWebHookURL("https://hooks.slack.com/services/T invalid")
	if isValid {
		t.Errorf("isValid should be '%v' actual: %v", false, isValid)
	}
}

func TestNewSlcLogger(t *testing.T) {

	logger, err := NewSlcLogger(&SlcLoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T",
	})

	if err != nil {
		t.Errorf("err should be 'nil' actual: %v", err.Error())
	}
	var expectedURL = "https://hooks.slack.com/services/T"
	if logger.WebHookURL != expectedURL {
		t.Errorf("WebHookURL should be %v actual: %v", expectedURL, logger.WebHookURL)
	}
	var expectedTitle = "Notification"
	if logger.DefaultTitle != expectedTitle {
		t.Errorf("DefaultTitle should be %v actual: %v", expectedTitle, logger.DefaultTitle)
	}
	if logger.LogLevel != LevelInfo {
		t.Errorf("LogLevel should be %v actual: %v", LevelInfo, logger.LogLevel)
	}
}

func TestNewSlcLoggerFullOptions(t *testing.T) {

	logger, err := NewSlcLogger(&SlcLoggerParams{
		WebHookURL:   "https://hooks.slack.com/services/T",
		DefaultTitle: "Custom Title",
		Channel:      "general",
		LogLevel:     LevelWarn,
		IconURL:      "https://example.com",
		UserName:     "My Logger",
	})

	if err != nil {
		t.Errorf("err should be 'nil' actual: %v", err.Error())
	}
	var expectedURL = "https://hooks.slack.com/services/T"
	if logger.WebHookURL != expectedURL {
		t.Errorf("WebHookURL should be %v actual: %v", expectedURL, logger.WebHookURL)
	}
	var expectedTitle = "Custom Title"
	if logger.DefaultTitle != expectedTitle {
		t.Errorf("DefaultTitle should be %v actual: %v", expectedTitle, logger.DefaultTitle)
	}
	var expectedChannel = "general"
	if logger.Channel != expectedChannel {
		t.Errorf("Channel should be %v actual: %v", expectedChannel, logger.Channel)
	}
	if logger.LogLevel != LevelWarn {
		t.Errorf("LogLevel should be %v actual: %v", LevelWarn, logger.LogLevel)
	}
	var expectedIconURL = "https://example.com"
	if logger.IconURL != expectedIconURL {
		t.Errorf("IconURL should be %v actual: %v", expectedIconURL, logger.IconURL)
	}
	var expectedUserName = "My Logger"
	if logger.UserName != expectedUserName {
		t.Errorf("UserName should be %v actual: %v", expectedUserName, logger.UserName)
	}
}

func TestValidateParams(t *testing.T) {

	err := validateParams(&SlcLoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T/valid",
		IconURL:    "https://example.com/vaild",
	})
	if err != nil {
		t.Errorf("err should be 'nil' actual: %v", err.Error())
	}

	err1 := validateParams(&SlcLoggerParams{})

	expectedErrMessage := "WebHookURL is a required parameter"
	if err1.Error() != expectedErrMessage {
		t.Errorf("error message should be '%v' actual: %v", expectedErrMessage, err1.Error())
	}

	invalidURLErr := validateParams(&SlcLoggerParams{
		WebHookURL: "invalidUrl",
	})

	expectedErrMessage = "WebHookURL must be a valid url"
	if invalidURLErr.Error() != expectedErrMessage {
		t.Errorf("error message should be '%v' actual: %v", expectedErrMessage, invalidURLErr.Error())
	}

	invalidDomainErr := validateParams(&SlcLoggerParams{
		WebHookURL: "invalidUrl",
	})
	if invalidDomainErr.Error() != expectedErrMessage {
		t.Errorf("error message should be '%v' actual: %v", expectedErrMessage, invalidDomainErr.Error())
	}

	invalidIconURLErr := validateParams(&SlcLoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T/valid",
		IconURL:    "invalid",
	})
	expectedErrMessage = "IconURL must be a valid url"
	if invalidIconURLErr.Error() != expectedErrMessage {
		t.Errorf("error message should be '%v' actual: %v", expectedErrMessage, invalidIconURLErr.Error())
	}
}

func TestSetLogLevel(t *testing.T) {

	logger, _ := NewSlcLogger(&SlcLoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T",
	})

	logger.SetLogLevel(LevelErr)

	if logger.LogLevel != LevelErr {
		t.Errorf("LogLevel should be %v actual: %v", LevelErr, logger.LogLevel)
	}
}

func TestBuildPayload(t *testing.T) {

	logger, _ := NewSlcLogger(&SlcLoggerParams{
		DefaultTitle: "title",
		WebHookURL:   "https://hooks.slack.com/services/T",
	})

	a := &attachment{Text: "body", Title: "title", Color: "good"}
	attachments := []attachment{*a}

	expected, _ := json.Marshal(payload{
		Attachments: attachments,
	})

	actual, err := logger.buildPayload("good", "body", []string{"title"})

	if err != nil {
		t.Errorf("err should be 'nil' actual: %v", err.Error())
	}

	if !bytes.Equal(actual, expected) {
		t.Errorf("built payload should be %v equall actual: %v", expected, actual)
	}
}
