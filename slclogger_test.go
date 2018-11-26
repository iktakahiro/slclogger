package slclogger

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidWebHookURL(t *testing.T) {

	assert.True(t, isValidWebHookURL("https://hooks.slack.com/services/T"))

	assert.False(t, isValidWebHookURL("invalid.slack.com/services/T"))

}

func TestNewSlcLogger(t *testing.T) {

	logger, err := NewSlcLogger(&LoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T",
	})
	assert.Nil(t, err)

	assert.Equal(t, "https://hooks.slack.com/services/T", logger.WebHookURL)

	assert.Equal(t, "Notification", logger.DefaultTitle)

	assert.Equal(t, LevelInfo, logger.LogLevel)

}

func TestNewSlcLoggerFullOptions(t *testing.T) {

	logger, err := NewSlcLogger(&LoggerParams{
		WebHookURL:     "https://hooks.slack.com/services/T",
		DefaultTitle:   "Custom Title",
		DefaultChannel: "general",
		WarnChannel:    "warn-channel",
		ErrorChannel:   "error-channel",
		LogLevel:       LevelWarn,
		IconURL:        "https://example.com/icon.png",
		UserName:       "My Logger",
	})

	assert.Nil(t, err)

	assert.Equal(t, "https://hooks.slack.com/services/T", logger.WebHookURL)

	assert.Equal(t, "Custom Title", logger.DefaultTitle)

	assert.Equal(t, "general", logger.DebugChannel)
	assert.Equal(t, "general", logger.InfoChannel)
	assert.Equal(t, "warn-channel", logger.WarnChannel)
	assert.Equal(t, "error-channel", logger.ErrorChannel)

	assert.Equal(t, LevelWarn, logger.LogLevel)

	assert.Equal(t, "https://example.com/icon.png", logger.IconURL)

	assert.Equal(t, "My Logger", logger.UserName)
}

func TestValidateParams(t *testing.T) {

	err := validateParams(&LoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T/valid",
		IconURL:    "https://example.com/vaild",
	})
	if err != nil {
		t.Errorf("err should be 'nil' actual: %v", err.Error())
	}

	err1 := validateParams(&LoggerParams{})

	expectedErrMessage := "WebHookURL is a required parameter"
	if err1.Error() != expectedErrMessage {
		t.Errorf("error message should be '%v' actual: %v", expectedErrMessage, err1.Error())
	}

	invalidURLErr := validateParams(&LoggerParams{
		WebHookURL: "invalidUrl",
	})

	expectedErrMessage = "WebHookURL must be a valid url"
	if invalidURLErr.Error() != expectedErrMessage {
		t.Errorf("error message should be '%v' actual: %v", expectedErrMessage, invalidURLErr.Error())
	}

	invalidDomainErr := validateParams(&LoggerParams{
		WebHookURL: "invalidUrl",
	})
	if invalidDomainErr.Error() != expectedErrMessage {
		t.Errorf("error message should be '%v' actual: %v", expectedErrMessage, invalidDomainErr.Error())
	}

	invalidIconURLErr := validateParams(&LoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T/valid",
		IconURL:    "invalid",
	})
	expectedErrMessage = "IconURL must be a valid url"
	if invalidIconURLErr.Error() != expectedErrMessage {
		t.Errorf("error message should be '%v' actual: %v", expectedErrMessage, invalidIconURLErr.Error())
	}
}

func TestSetLogLevel(t *testing.T) {

	logger, _ := NewSlcLogger(&LoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T",
	})

	// before
	assert.Equal(t, LevelInfo, logger.LogLevel)

	logger.SetLogLevel(LevelError)

	// after
	assert.Equal(t, LevelError, logger.LogLevel)
}

func TestBuildPayload(t *testing.T) {

	logger, _ := NewSlcLogger(&LoggerParams{
		DefaultTitle: "title",
		WebHookURL:   "https://hooks.slack.com/services/T",
	})

	a := &attachment{Text: "body", Title: "title", Color: "good"}
	attachments := []attachment{*a}

	expected, _ := json.Marshal(payload{
		Channel:     "random",
		Attachments: attachments,
	})

	actual, err := logger.buildPayload("random", "good", "body", []string{"title"})

	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}
