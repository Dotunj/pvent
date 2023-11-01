package util

import (
	"bytes"
	"fmt"
	"os"
)

func PublishReport(rate int, successMessages, failedMessages uint64) error {
	buf := bytes.Buffer{}

	v := "Total Messages: %v\n Sent: %v\n Failed: %v\n"
	report := fmt.Sprintf(v, rate, successMessages, failedMessages)

	_, err := buf.WriteString(report)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(os.Stdout)
	return err
}
