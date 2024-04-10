package log

import (
	"os"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	logger      = initLogger()
	errorLogger = initErrorLogger()
	loggerStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1, 0, 1)
)

func initLogger() *log.Logger {
	styles := log.DefaultStyles()
	styles.Levels[log.DebugLevel] = loggerStyle.Copy().
		SetString("DEBUG").
		Background(lipgloss.Color("#414868")).
		Foreground(lipgloss.Color("#FFFFFF"))

	styles.Levels[log.InfoLevel] = loggerStyle.Copy().
		SetString("INFO").
		Background(lipgloss.Color("#485E30")).
		Foreground(lipgloss.Color("#FFFFFF"))

	newLogger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		//ReportCaller:    true,
		Level: log.DebugLevel,
	})
	newLogger.SetStyles(styles)
	return newLogger
}

func initErrorLogger() *log.Logger {
	styles := log.DefaultStyles()

	styles.Levels[log.ErrorLevel] = loggerStyle.Copy().
		SetString("ERROR!").
		Background(lipgloss.Color("#F7768E")).
		Foreground(lipgloss.Color("#000000"))

	styles.Levels[log.WarnLevel] = loggerStyle.Copy().
		SetString("WARNING!").
		Background(lipgloss.Color("#FF9E64")).
		Foreground(lipgloss.Color("#000000"))

	styles.Levels[log.FatalLevel] = loggerStyle.Copy().
		SetString("FATAL!").
		Background(lipgloss.Color("#BB9AF7")).
		Foreground(lipgloss.Color("#000000"))

	newLogger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		//ReportCaller:    true,
	})
	newLogger.SetStyles(styles)
	return newLogger
}

// Wrapper function to log errors if they exist
func ErrorWrapper(err error) error {
	if err != nil {
		programCounter, file, line, ok := runtime.Caller(1)
		if !ok {
			errorLogger.Error("Error retrieving stack trace.")
		} else {
			errorLogger.Error("Stack trace:",
				"function", runtime.FuncForPC(programCounter).Name(),
				"file", file,
				"line", line,
			)
		}
		errorLogger.Error(err)
	}
	return err
}

func FatalWrapper(err error) error {
	if err != nil {
		programCounter, file, line, ok := runtime.Caller(1)
		if !ok {
			errorLogger.Error("Error retrieving stack trace.")
		} else {
			errorLogger.Error("Stack trace:",
				"function", runtime.FuncForPC(programCounter).Name(),
				"file", file,
				"line", line,
			)
		}
		errorLogger.Fatal(err)
	}
	return err
}

func DebugCaller(message interface{}, keyValue ...interface{}) {

	logger.Debug(message, keyValue...)
	programCounter, file, line, ok := runtime.Caller(1)
	if !ok {
		logger.Debug("Error retrieving stack trace.")
	} else {
		logger.Debug("Stack trace:",
			"function", runtime.FuncForPC(programCounter).Name(),
			"file", file,
			"line", line,
		)
	}
}

func Error(message interface{}, keyValue ...interface{}) {
	errorLogger.Error(message, keyValue...)
}

func Fatal(message interface{}, keyValue ...interface{}) {
	errorLogger.Fatal(message, keyValue...)
}

func Warn(message interface{}, keyValue ...interface{}) {
	logger.Warn(message, keyValue...)
}

func Debug(message interface{}, keyValue ...interface{}) {
	logger.Debug(message, keyValue...)
}

func Info(message interface{}, keyValue ...interface{}) {
	logger.Info(message, keyValue...)
}
