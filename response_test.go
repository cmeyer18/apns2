package apns2

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Unit Tests

func TestResponseSent(t *testing.T) {
	assert.Equal(t, http.StatusOK, StatusSent)
	assert.Equal(t, true, (&Response{StatusCode: 200}).Sent())
	assert.Equal(t, false, (&Response{StatusCode: 400}).Sent())
}

func TestIntTimestampParse(t *testing.T) {
	response := &Response{}
	payload := "{\"reason\":\"Unregistered\", \"timestamp\":1458114061260}"
	err := json.Unmarshal([]byte(payload), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1458114061260)/1000, response.Timestamp.Unix())
}

func TestInvalidTimestampParse(t *testing.T) {
	response := &Response{}
	payload := "{\"reason\":\"Unregistered\", \"timestamp\": \"2016-01-16 17:44:04 +1300\"}"
	err := json.Unmarshal([]byte(payload), &response)
	assert.Error(t, err)
}
