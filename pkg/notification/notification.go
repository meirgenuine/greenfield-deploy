package notification

type Notification struct {
	Message string
}

type User struct {
	ChatID int64
}

func (nt Notification) String() string {
	return nt.Message
}
