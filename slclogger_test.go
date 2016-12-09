package slclogger

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SlcLoggerTestSuite struct {
	suite.Suite
}

func (suite *SlcLoggerTestSuite) SetupTest() {
	return
}

func (suite *SlcLoggerTestSuite) TestIsValidWebHookURL() {

	assert.True(suite.T(), isValidWebHookURL("https://hooks.slack.com/services/T"))
	assert.False(suite.T(), isValidWebHookURL("invaild.slack.com/services/T"))
	assert.False(suite.T(), isValidWebHookURL("https://hooks.slack.com/services/T invalid"))
}

func (suite *SlcLoggerTestSuite) TestNewSlcLogger() {

	logger, err := NewSlcLogger(&SlcLoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T",
	})

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "https://hooks.slack.com/services/T", logger.WebHookURL)
	assert.Equal(suite.T(), "Notification", logger.DefaultTitle)
	assert.Equal(suite.T(), LevelInfo, logger.LogLevel)
}

func (suite *SlcLoggerTestSuite) TestNewSlcLoggerFullOptions() {

	logger, err := NewSlcLogger(&SlcLoggerParams{
		WebHookURL:   "https://hooks.slack.com/services/T",
		DefaultTitle: "Default Title",
		Channel:      "general",
		LogLevel:     LevelWarn,
		IconURL:      "https://example.com",
		UserName:     "My Logger",
	})

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "https://hooks.slack.com/services/T", logger.WebHookURL)
	assert.Equal(suite.T(), "general", logger.Channel)
	assert.Equal(suite.T(), "Default Title", logger.DefaultTitle)
	assert.Equal(suite.T(), LevelWarn, logger.LogLevel)
	assert.Equal(suite.T(), "https://example.com", logger.IconURL)
	assert.Equal(suite.T(), "My Logger", logger.UserName)
}

func (suite *SlcLoggerTestSuite) TestValidateParams() {

	err := validateParams(&SlcLoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T/valid",
		IconURL:    "https://example.com/vaild",
	})
	assert.Nil(suite.T(), err)

	err1 := validateParams(&SlcLoggerParams{})
	assert.NotNil(suite.T(), err1)

	err2 := validateParams(&SlcLoggerParams{
		WebHookURL: "invaildUrl",
	})
	assert.NotNil(suite.T(), err2)

	err3 := validateParams(&SlcLoggerParams{
		WebHookURL: "https://example.com/invalidUrl/",
	})
	assert.NotNil(suite.T(), err3)

	err4 := validateParams(&SlcLoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T/valid",
		IconURL:    "invalid",
	})
	assert.NotNil(suite.T(), err4)
}

func (suite *SlcLoggerTestSuite) TestSetLogLevel() {

	logger, _ := NewSlcLogger(&SlcLoggerParams{
		WebHookURL: "https://hooks.slack.com/services/T",
	})
	logger.SetLogLevel(LevelWarn)

	assert.Equal(suite.T(), LevelWarn, logger.LogLevel)
}

func (suite *SlcLoggerTestSuite) TestBuildPayload() {

	logger, _ := NewSlcLogger(&SlcLoggerParams{
		DefaultTitle: "title",
		WebHookURL:   "https://hooks.slack.com/services/T",
	})

	a := &Attachment{Text: "body", Title: "title", Color: "good"}
	attachments := []Attachment{*a}

	expected, err := json.Marshal(Payload{

		Attachments: attachments,
	})

	actual, _ := logger.buildPayload("good", "body", []string{"title"})

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), expected, actual)
}

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) DummySendNotification() error {
	return nil
}

func (suite *SlcLoggerTestSuite) TestSendNotification() {

	// TODO
	// logger, _ := NewSlcLogger(&SlcLoggerParams{
	//	WebHookURL: "https://hooks.slack.com/services/T",
	//})

	// err := logger.sendNotification(LevelErr, "warning", "message", []string{"title"})
}

func (suite *SlcLoggerTestSuite) TestInfo() {

	// TODO
	// logger, _ := NewSlcLogger(&SlcLoggerParams{
	//	WebHookURL: "https://hooks.slack.com/services/T",
	//})

	// err := logger.Info("message", "title")
}

func TestModelRecipientTestSuite(t *testing.T) {
	suite.Run(t, new(SlcLoggerTestSuite))
}
