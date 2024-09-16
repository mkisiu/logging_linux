package logging_linux

type SyslogLinux struct {
	Networks    string
	Server      string
	ProcessTag  string
	Action      string
	LogFileName string
	LogChannel  chan LogEntry
	LogEntry    LogEntry
}

type LogEntry struct {
	LogLevel string
	Msg      string
}
