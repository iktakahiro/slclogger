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

const (
	LevelDebug logLevel = iota + 1
	LevelInfo
	LevelWarn
	LevelErr
)

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
	SlcLoggerParams
	LogLevel     logLevel
	Channel      string
	UserName     string
	DefaultTitle string
	IconURL      string
}

// SlcLoggerParams is the set of parameters that can be used when creating a new SlcLogger.
type SlcLoggerParams struct {
	WebHookURL   string
	LogLevel     logLevel
	Channel      string
	UserName     string
	DefaultTitle string
	IconURL      string
}

// isValidWebHookURL validates whether the value is valid as a URL and returns bool.
func isValidWebHookURL(url string) bool {

	if !govalidator.IsURL(url) {
		return false
	}
	return strings.HasPrefix(url, "https://hooks.slack.com/services/")
}

func validateParams(params *SlcLoggerParams) error {
	if params.WebHookURL == "" {
		return &SlcErr{errors.New("WebHookUrl is a required parameter."), 0}
	}
	if !isValidWebHookURL(params.WebHookURL) {
		return &SlcErr{errors.New("WebHookUrl must be a valid webhook url."), 0}
	}

	if params.IconURL != "" && !govalidator.IsURL(params.IconURL) {
		return &SlcErr{errors.New("WebHookUrl must be a valid url."), 0}
	}
	return nil
}

// NewSlcLogger returns a new SLcLogger.
func NewSlcLogger(params *SlcLoggerParams) (*SlcLogger, error) {

	if err := validateParams(params); err != nil {
		return nil, err
	}

	var logLevel logLevel
	if params.LogLevel == 0 {
		logLevel = LevelInfo
	} else {
		logLevel = params.LogLevel
	}
	var defaultTitle string
	if params.DefaultTitle == "" {
		defaultTitle = "Notification"
	} else {
		defaultTitle = params.DefaultTitle
	}

	return &SlcLogger{
		WebHookURL:   params.WebHookURL,
		LogLevel:     logLevel,
		Channel:      params.Channel,
		DefaultTitle: defaultTitle,
		UserName:     params.UserName,
		IconURL:      params.IconURL,
	}, nil
}

type Payload struct {
	Channel     string       `json:"channel"`
	UserName    string       `json:"username"`
	IconURL     string       `json:"iconUrl"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"titleLink"`
	Text      string `json:"text"`
	Color     string `json:"color"`
}

func (s *SlcLogger) buildPayload(color, message string, titleParam []string) ([]byte, error) {

	var title string
	if len(titleParam) == 0 {
		title = s.DefaultTitle
	} else {
		title = titleParam[0]
	}

	a := &Attachment{Text: message, Title: title, Color: color}
	attachments := []Attachment{*a}

	return json.Marshal(Payload{
		Channel:     s.Channel,
		UserName:    s.UserName,
		IconURL:     s.IconURL,
		Attachments: attachments,
	})
}

// SetLogLevel sets a logLevel to SlcLogger.LogLevel.
func (s *SlcLogger) SetLogLevel(level logLevel) {
	s.LogLevel = level
}

// sendNotification posts a message to WebHookURL.
func (s *SlcLogger) sendNotification(logLevel logLevel, color, message string, titleParam []string) error {

	if logLevel < s.LogLevel {
		return nil
	}

	payload, err := s.buildPayload(color, message, titleParam)
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
	ColorDebug string = "#03A9F4"
	ColorInfo  string = "good"
	ColorWarn  string = "warning"
	ColorErr   string = "danger"
)

// Debug is a wrapper function of sendNotification function that implicitly sets the logLevel and color.
func (s *SlcLogger) Debug(message string, title ...string) error {
	return s.sendNotification(LevelDebug, ColorDebug, message, title)
}

// Info is a wrapper function of sendNotification function that implicitly sets the logLevel and color.
func (s *SlcLogger) Info(message string, title ...string) error {
	return s.sendNotification(LevelInfo, ColorInfo, message, title)
}

// Warn is a wrapper function of sendNotification function that implicitly sets the logLevel and color.
func (s *SlcLogger) Warn(message string, title ...string) error {
	return s.sendNotification(LevelWarn, ColorWarn, message, title)
}

// Err is a wrapper function of sendNotification function that implicitly sets the logLevel and color.
func (s *SlcLogger) Err(message string, title ...string) error {
	return s.sendNotification(LevelErr, ColorErr, message, title)
}
