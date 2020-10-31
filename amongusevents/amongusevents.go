package amongusevents

type amongUsEventState struct {
	eventTitle      string
	eventAttendees  []string
	eventCantAttend []string
}

func CreateEvent(title string) string {
	return title
}

func extractEventState(event string) amongUsEventState {
	return amongUsEventState{}
}

func removeUnTrackedReactions(event string) string {
	return event
}
