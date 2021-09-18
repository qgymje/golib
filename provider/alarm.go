package provider

// IAlarm is an alarmr when things going wrong
type IAlarm interface {
	Send(subject, body string) error
	SendTo(to []string, subject, body string) error
}
