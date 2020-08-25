package view

import "log"

// PublicError is used to distinguish between user and system errors
type PublicError interface {
	error
	Public() string
}

// Data is the top level structure that views expect for data
type Data struct {
	Alert *Alert
	Yield interface{}
}

// SetAlert sets the Alert field on Data
func (d *Data) SetAlert(err error) {
	var msg string

	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

// AlertError provides a method to create custom alert messages
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}
