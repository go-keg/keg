package alert

import "fmt"

type Report func(record *ErrorRecord)

var defaultReport Report = func(record *ErrorRecord) {
	if record.Count == 1 {
		fmt.Println("[ALERT] " + record.RawMsg)
	}
}
