package logger

type Records struct {
	log      *Logger
	tagName  string
	dataName string
}

type Record struct {
	data    map[string]any
	event   *Event
	records *Records
}

func (s *Records) SetTagName(tagName string) *Records {
	s.tagName = tagName
	return s
}

func (s *Records) SetDataName(dataName string) *Records {
	s.dataName = dataName
	return s
}

func (s *Records) Trace(events ...func(event *Event)) *Record {
	record := &Record{
		event:   s.log.Trace(),
		records: s,
	}
	for _, f := range events {
		f(record.event)
	}
	return record
}

func (s *Records) Debug(events ...func(event *Event)) *Record {
	record := &Record{
		event:   s.log.Debug(),
		records: s,
	}
	for _, f := range events {
		f(record.event)
	}
	return record
}

func (s *Records) Info(events ...func(event *Event)) *Record {
	record := &Record{
		event:   s.log.Info(),
		records: s,
	}
	for _, f := range events {
		f(record.event)
	}
	return record
}

func (s *Records) Warn(events ...func(event *Event)) *Record {
	record := &Record{
		event:   s.log.Warn(),
		records: s,
	}
	for _, f := range events {
		f(record.event)
	}
	return record
}

func (s *Records) Error(events ...func(event *Event)) *Record {
	record := &Record{
		event:   s.log.Error(),
		records: s,
	}
	for _, f := range events {
		f(record.event)
	}
	return record
}

func (s *Records) Fatal(events ...func(event *Event)) *Record {
	record := &Record{
		event:   s.log.Fatal(),
		records: s,
	}
	for _, f := range events {
		f(record.event)
	}
	return record
}

func (s *Records) Panic(events ...func(event *Event)) *Record {
	record := &Record{
		event:   s.log.Panic(),
		records: s,
	}
	for _, f := range events {
		f(record.event)
	}
	return record
}

// Event Handling custom logic.
func (s *Record) Event(f func(event *Event)) *Record {
	if f != nil {
		f(s.event)
	}
	return s
}

func (s *Record) Tag(tag string) *Record {
	return s.Event(func(event *Event) { event.Str(s.records.tagName, tag) })
}

// Set The Set method should be called before the Msg or Err method.
// Make sure the `value` parameter value can be serialized to json.
func (s *Record) Set(key string, value any) *Record {
	if s.data == nil {
		s.data = make(map[string]any, 1<<3)
	}
	if key != "" {
		s.data[key] = value
	}
	return s
}

func (s *Record) log() *Event {
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
	log.AddEvent(func(event *Event) {
		event.CallerSkipFrame(1)
	})
	return &Records{
		log:      log,
		tagName:  "tag",
		dataName: "data",
	}
}
