package slclogger

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
)

type logLevel int

// LogLevel Constants
const (
	LevelDebug logLevel = iota + 1
	LevelInfo
	LevelWarn
	LevelError
)

// SlcErr is a struct for Slack response error code.
type SlcErr struct {
	Err  error
	Code int
}

// Error is an implementation of the error interface.
func (s *SlcErr) Error() string {
	return s.Err.Error()
}

// SlcLogger is the structure detailing logger.
type SlcLogger struct {
	WebHookURL string
	LoggerParams
	LogLevel     logLevel
	DebugChannel string
	InfoChannel  string
	WarnChannel  string
	ErrorChannel string
	UserName     string
	DefaultTitle string
	IconURL      string
}

// LoggerParams is the set of parameters that can be used when creating a new SlcLogger.
type LoggerParams struct {
	WebHookURL     string
	LogLevel       logLevel
	DefaultChannel string
	DebugChannel   string
	InfoChannel    string
	WarnChannel    string
	ErrorChannel   string
	UserName       string
	DefaultTitle   string
	IconURL        string
}

// isValidWebHookURL validates whether the value is valid as a URL and returns bool.
func isValidWebHookURL(url string) bool {
	if !govalidator.IsURL(url) {
		return false
	}
	return strings.HasPrefix(url, "https://hooks.slack.com/services/")
}

// validateParams validates LoggerParams.
func validateParams(params *LoggerParams) error {

	if params.WebHookURL == "" {
		return &SlcErr{errors.New("WebHookURL is a required parameter"), 0}
	}
	if !isValidWebHookURL(params.WebHookURL) {
		return &SlcErr{errors.New("WebHookURL must be a valid url"), 0}
	}
	if params.IconURL != "" && !govalidator.IsURL(params.IconURL) {
		return &SlcErr{errors.New("IconURL must be a valid url"), 0}
	}
	return nil
}

// NewSlcLogger returns a new SLcLogger.
func NewSlcLogger(params *LoggerParams) (*SlcLogger, error) {

	if err := validateParams(params); err != nil {
		return nil, err
	}

	var level logLevel
	if params.LogLevel == 0 {
		level = LevelInfo
	} else {
		level = params.LogLevel
	}

	var defaultTitle string
	if params.DefaultTitle == "" {
		defaultTitle = "Notification"
	} else {
		defaultTitle = params.DefaultTitle
	}
	debugChannel := params.DebugChannel
	infoChannel := params.InfoChannel
	warnChannel := params.WarnChannel
	errorChannel := params.ErrorChannel

	if debugChannel == "" {
		debugChannel = params.DefaultChannel
	}
	if infoChannel == "" {
		infoChannel = params.DefaultChannel
	}
	if warnChannel == "" {
		warnChannel = params.DefaultChannel
	}
	if errorChannel == "" {
		errorChannel = params.DefaultChannel
	}

	return &SlcLogger{
		WebHookURL:   params.WebHookURL,
		LogLevel:     level,
		DebugChannel: debugChannel,
		InfoChannel:  infoChannel,
		WarnChannel:  warnChannel,
		ErrorChannel: errorChannel,
		DefaultTitle: defaultTitle,
		UserName:     params.UserName,
		IconURL:      params.IconURL,
	}, nil
}

// SetLogLevel sets a logLevel to SlcLogger.LogLevel.
func (s *SlcLogger) SetLogLevel(level logLevel) {
	s.LogLevel = level
}

type payload struct {
	Channel     string       `json:"channel"`
	UserName    string       `json:"username"`
	IconURL     string       `json:"icon_url"`
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	Color     string `json:"color"`
}

// buildPayload returns an encoded JSON Object to post Slack API.
func (s *SlcLogger) buildPayload(slackChannel, color, message string, titleParam []string) ([]byte, error) {

	var title string
	if len(titleParam) == 0 {
		title = s.DefaultTitle
	} else {
		title = titleParam[0]
	}

	a := &attachment{Text: message, Title: title, Color: color}
	attachments := []attachment{*a}

	return json.Marshal(payload{
		Channel:     slackChannel,
		UserName:    s.UserName,
		IconURL:     s.IconURL,
		Attachments: attachments,
	})
}

func (s *SlcLogger) getTargetChannel(logLevel logLevel) string {
	switch logLevel {
	case 1:
		return s.DebugChannel
	case 2:
		return s.InfoChannel
	case 3:
		return s.WarnChannel
	case 4:
		return s.ErrorChannel
	}
	return s.DefaultChannel
}

// sendNotification posts a message to the WebHookURL of Slack.
func (s *SlcLogger) sendNotification(logLevel logLevel, color string, message interface{}, titleParam []string) error {

	var text string
	if t, ok := message.(error); ok {
		text = t.Error()
	} else if t, ok := message.(string); ok {
		text = t
	} else {
		return &SlcErr{errors.New("the type of message parameter should be string or error"), 0}
	}

	if logLevel < s.LogLevel {
		return nil
	}
	slackChannel := s.getTargetChannel(logLevel)

	payload, err := s.buildPayload(slackChannel, color, text, titleParam)
	if err != nil {
		return &SlcErr{err, 0}
	}

	resp, err := http.Post(s.WebHookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return &SlcErr{err, 0}
	}
	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		return &SlcErr{errors.New(string(body)), resp.StatusCode}
	}

	return nil
}

const (
	colorDebug string = "#03A9F4"
	colorInfo  string = "good"
	colorWarn  string = "warning"
	colorErr   string = "danger"
)

// Debug is a wrapper function of sendNotification function that implicitly sets the logLevel and color.
func (s *SlcLogger) Debug(message interface{}, title ...string) error {
	return s.sendNotification(LevelDebug, colorDebug, message, title)
}

// Info is a wrapper function of sendNotification function that implicitly sets the logLevel and color.
func (s *SlcLogger) Info(message interface{}, title ...string) error {
	return s.sendNotification(LevelInfo, colorInfo, message, title)
}

// Warn is a wrapper function of sendNotification function that implicitly sets the logLevel and color.
func (s *SlcLogger) Warn(message interface{}, title ...string) error {
	return s.sendNotification(LevelWarn, colorWarn, message, title)
}

// Error is a wrapper function of sendNotification function that implicitly sets the logLevel and color.
func (s *SlcLogger) Error(message interface{}, title ...string) error {
	return s.sendNotification(LevelError, colorErr, message, title)
}
