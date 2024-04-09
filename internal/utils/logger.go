package utils

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	logger    = initLogger()
	errLogger = initErrorLogger()
)

func initLogger() *log.Logger {
	styles := log.DefaultStyles()
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBUG").
		Bold(true).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#414868")).
		Foreground(lipgloss.Color("#FFFFFF"))

	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Bold(true).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#485E30")).
		Foreground(lipgloss.Color("#FFFFFF"))
	
	newLogger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		ReportCaller:    true,
		Level:           log.DebugLevel,
	})
	newLogger.SetStyles(styles)
	return newLogger
}

func initErrorLogger() *log.Logger {
	styles := log.DefaultStyles()

	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR!").
		Bold(true).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#F7768E")).
		Foreground(lipgloss.Color("#000000"))

	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARNING!").
		Bold(true).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#FF9E64")).
		Foreground(lipgloss.Color("#000000"))

	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("FATAL!").
		Bold(true).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#BB9AF7")).
		Foreground(lipgloss.Color("#000000"))

	newLogger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		ReportCaller:    true,
	})
	newLogger.SetStyles(styles)
	return newLogger
}

// Wrapper function to log errors if they exist
func LogError(err error) error {
	if err != nil {
		errLogger.Error(err)
	}
	return err
}

func LogFatal(err error) error {
	if err != nil {
		errLogger.Fatal(err)
	}
	return err
}
