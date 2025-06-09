package logger

import (
	"github.com/rs/zerolog"
)

type Records struct {
	log      *Logger
	fc       func(event *zerolog.Event)
	tagName  string
	dataName string
}

type Record struct {
	records *Records
	event   *zerolog.Event
	data    map[string]interface{}
}

// Func Public custom processing Event.
func (s *Records) Func(fc func(event *zerolog.Event)) *Records {
	s.fc = fc
	return s
}

func (s *Records) SetTagName(tagName string) *Records {
	s.tagName = tagName
	return s
}

func (s *Records) SetDataName(dataName string) *Records {
	s.dataName = dataName
	return s
}

func (s *Records) Trace() *Record {
	record := &Record{
		records: s,
		event:   s.log.Trace(),
	}
	if s.fc != nil {
		record.event.Func(s.fc)
	}
	return record
}

func (s *Records) Debug() *Record {
	record := &Record{
		records: s,
		event:   s.log.Debug(),
	}
	if s.fc != nil {
		record.event.Func(s.fc)
	}
	return record
}

func (s *Records) Info() *Record {
	record := &Record{
		records: s,
		event:   s.log.Info(),
	}
	if s.fc != nil {
		record.event.Func(s.fc)
	}
	return record
}

func (s *Records) Warn() *Record {
	record := &Record{
		records: s,
		event:   s.log.Warn(),
	}
	if s.fc != nil {
		record.event.Func(s.fc)
	}
	return record
}

func (s *Records) Error() *Record {
	record := &Record{
		records: s,
		event:   s.log.Error(),
	}
	if s.fc != nil {
		record.event.Func(s.fc)
	}
	return record
}

func (s *Records) Fatal() *Record {
	record := &Record{
		records: s,
		event:   s.log.Fatal(),
	}
	if s.fc != nil {
		record.event.Func(s.fc)
	}
	return record
}

func (s *Records) Panic() *Record {
	record := &Record{
		records: s,
		event:   s.log.Panic(),
	}
	if s.fc != nil {
		record.event.Func(s.fc)
	}
	return record
}

// Func Add custom key-value pairs to the log object.
func (s *Record) Func(fc func(event *zerolog.Event)) *Record {
	if fc != nil {
		fc(s.event)
	}
	return s
}

func (s *Record) Tag(tag string) *Record {
	return s.Func(func(event *zerolog.Event) { event.Str(s.records.tagName, tag) })
}

// Set The Set method should be called before the Msg or Err method.
// Make sure the `value` parameter value can be serialized to json.
func (s *Record) Set(key string, value interface{}) *Record {
	if s.data == nil {
		s.data = make(map[string]interface{}, 1<<3)
	}
	if key != "" {
		s.data[key] = value
	}
	return s
}

func (s *Record) log() *zerolog.Event {
	tmp := s.event
	if s.data != nil {
		tmp.Any(s.records.dataName, s.data)
	}
	return tmp
}

func (s *Record) Msg(msg string) {
	s.log().Msg(msg)
}

func (s *Record) Err(err error) {
	s.log().Msg(err.Error())
}

func NewRecords(log *Logger) *Records {
	return &Records{
		log:      log,
		fc:       func(event *zerolog.Event) { event.CallerSkipFrame(1) },
		tagName:  "tag",
		dataName: "data",
	}
}
