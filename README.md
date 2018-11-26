# SlcLogger

[![Build Status](https://travis-ci.org/iktakahiro/slclogger.svg?branch=master)](https://travis-ci.org/iktakahiro/slclogger)

**Simple and Human Friendly Slack Client for Logging/Notification written in Go**

## Install

```bash
go get "github.com/iktakahiro/slclogger"
```

## How to Use

### Basic Usage

```go
package main

import (
	"errors"

	"github.com/iktakahiro/slclogger"
)

func something() error {
	return errors.New("an error has occurred")
}

func main() {

	logger, _ := slclogger.NewSlcLogger(&slclogger.LoggerParams{
		WebHookURL: "https://hooks.slack.com/services/YOUR_WEBHOOK_URL",
	})

	if err := something(); err != nil {
		logger.Error(err, "Error Notification")
	}
}
```

When you execute the above sample code, your Slack channel will receive the message.

![](./example/example-slack1.png)

### Log Levels

The default log level is *Info*. You can set it when initializing a SlcLogger struct.

```go
package main

import (
	"github.com/iktakahiro/slclogger"
)

func main() {
    logger, _ := slclogger.NewSlcLogger(&slclogger.LoggerParams{
        WebHookURL: "https://hooks.slack.com/services/YOUR_WEBHOOK_URL",
        LogLevel: slclogger.LevelDebug,
    })

    logger.Debug("Debug Message")
    logger.Info("Info Message")
    logger.Warn("Warn Message")
    logger.Error("Error Message")
}
```

![](./example/example-slack2.png)

You can also change the level at any time.

```go
package main

import (
	"github.com/iktakahiro/slclogger"
)

func main() {
	logger, _ := slclogger.NewSlcLogger(&slclogger.LoggerParams{
		WebHookURL: "https://hooks.slack.com/services/YOUR_WEBHOOK_URL",
	})

    logger.SetLogLevel(slclogger.LevelWarn)

    // The following notification will be ignored.
    logger.Debug("Debug Message")
}
```

### Configure Options

All options are shown below.

```go
package main

import (
	"github.com/iktakahiro/slclogger"
)

func main() {

    logger, err := slclogger.NewSlcLogger(&slclogger.LoggerParams{
        WebHookURL:         "https://hooks.slack.com/services/YOUR_WEBHOOK_URL",
        DefaultTitle:       "Default Title",
        DefaultChannel:     "general",
        DebugChannel:       "debug-channel",
        InfoChannel:        "info-channel",
        WarnChannel :       "warn-channel",
        ErrorChannel :      "error-channel",
        LogLevel:           slclogger.LevelWarn,
        IconURL:            "https://example.com",
        UserName:           "My Logger",
    })
}
```

Param | Default Value
------ | ------------
WebHookURL (*require*) | --
DefaultTitle | "Notification"
DefaultChannel | "" (When this param is omitted, the default channel of specified WebHook is used.)
DebugChannel | "" (When this param is omitted, the value of DefaultChannel is used.)
InfoChannel | "" (When this param is omitted,  the value of DefaultChannel is used.)
WarnChannel | "" (When this param is omitted,  the value of DefaultChannel  is used.)
ErrorChannel | "" (When this param is omitted, the value of DefaultChannel  is used.)
LogLevel | Info
IconURL | ""
UseName | ""

## Error Handling

If you want to handle errors, use SlcErr.

```go
if err := logger.Info("info message"); err != nil {
    if slcErr, ok := err.(*slclogger.SlcErr); ok {
        fmt.Println(slcErr)
        fmt.Println(slcErr.Code)
    }
}
```

## Test

```bash
make test
```

## Documents

- [slclogger \- GoDoc](https://godoc.org/github.com/iktakahiro/slclogger)
