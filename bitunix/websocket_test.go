package bitunix

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHeartbeatMessage(t *testing.T) {
	timestamp := time.Now().Unix()
	message := heartbeatMessage{
		Op:   "ping",
		Ping: timestamp,
	}

	bytes, err := json.Marshal(message)
	require.NoError(t, err)

	var decoded heartbeatMessage
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "ping", decoded.Op)
	assert.Equal(t, timestamp, decoded.Ping)
}
