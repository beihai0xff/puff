package log

var defaultTestConfig = &OutputConfig{

	Writer:    OutputConsole,
	Level:     "info",
	Formatter: "json",
}
