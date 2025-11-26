package alert

import (
	"log"
)

type Report func(record *ErrorRecord)

var defaultReport Report = func(record *ErrorRecord) {
	if record.Count == 1 {
		log.Println("[ALERT] " + record.RawMsg)
	}
}
