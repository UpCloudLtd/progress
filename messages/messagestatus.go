package messages

type MessageStatus string

const (
	MessageStatusPending MessageStatus = "pending"
	MessageStatusStarted MessageStatus = "started"
	MessageStatusSuccess MessageStatus = "success"
	MessageStatusWarning MessageStatus = "warning"
	MessageStatusError   MessageStatus = "error"
	MessageStatusSkipped MessageStatus = "skipped"
	MessageStatusUnknown MessageStatus = "unknown"
)

var validUpdateStatuses = map[MessageStatus]bool{
	MessageStatusPending: true,
	MessageStatusStarted: true,
	MessageStatusSuccess: true,
	MessageStatusWarning: true,
	MessageStatusError:   true,
	MessageStatusSkipped: true,
	MessageStatusUnknown: true,
}

var finishedUpdateStatuses = map[MessageStatus]bool{
	MessageStatusSuccess: true,
	MessageStatusWarning: true,
	MessageStatusError:   true,
	MessageStatusSkipped: true,
	MessageStatusUnknown: true,
}

func (ms MessageStatus) IsValid() bool {
	return validUpdateStatuses[ms]
}

func (ms MessageStatus) IsInProgress() bool {
	return ms == MessageStatusStarted
}

func (ms MessageStatus) IsFinished() bool {
	return finishedUpdateStatuses[ms]
}
