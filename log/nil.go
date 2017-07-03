package log

// Nil is a nil logger
var Nil = NilLogger{}

// NilLogger is an empty logger mainly used for tests
type NilLogger struct{}

func (l NilLogger) Debug(...interface{})                           {}
func (l NilLogger) Debugln(...interface{})                         {}
func (l NilLogger) Debugf(string, ...interface{})                  {}
func (l NilLogger) Info(...interface{})                            {}
func (l NilLogger) Infoln(...interface{})                          {}
func (l NilLogger) Infof(string, ...interface{})                   {}
func (l NilLogger) Warn(...interface{})                            {}
func (l NilLogger) Warnln(...interface{})                          {}
func (l NilLogger) Warnf(string, ...interface{})                   {}
func (l NilLogger) Error(...interface{})                           {}
func (l NilLogger) Errorln(...interface{})                         {}
func (l NilLogger) Errorf(string, ...interface{})                  {}
func (l NilLogger) Fatal(...interface{})                           {}
func (l NilLogger) Fatalln(...interface{})                         {}
func (l NilLogger) Fatalf(string, ...interface{})                  {}
func (l NilLogger) Panic(...interface{})                           {}
func (l NilLogger) Panicln(...interface{})                         {}
func (l NilLogger) Panicf(string, ...interface{})                  {}
func (l NilLogger) With(key string, value interface{}) Logger      { return l }
func (l NilLogger) WithField(key string, value interface{}) Logger { return l }
func (l NilLogger) Set(level Level) error                          { return nil }
