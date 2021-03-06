package storage

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/emitter-io/emitter/broker/message"
	"github.com/emitter-io/emitter/utils"
	"github.com/stretchr/testify/assert"
)

// awaiter represents a query awaiter.
type mockAwaiter struct {
	f func(timeout time.Duration) (r [][]byte)
}

func (a *mockAwaiter) Gather(timeout time.Duration) (r [][]byte) {
	return a.f(timeout)
}

type testStorageConfig struct {
	Provider string                 `json:"provider"`
	Config   map[string]interface{} `json:"config,omitempty"`
}

func newTestMemStore() *InMemory {
	s := new(InMemory)
	s.Configure(nil)

	s.Store(testMessage(1, 1, 1))
	s.Store(testMessage(1, 1, 2))
	s.Store(testMessage(1, 2, 1))
	s.Store(testMessage(1, 2, 2))
	s.Store(testMessage(1, 3, 1))
	s.Store(testMessage(1, 3, 2))
	return s
}

func TestInMemory_Name(t *testing.T) {
	s := NewInMemory(nil)
	assert.Equal(t, "inmemory", s.Name())
}

func TestInMemory_Configure(t *testing.T) {
	s := new(InMemory)
	cfg := map[string]interface{}{
		"maxsize": float64(1),
		"prune":   float64(1),
	}

	err := s.Configure(cfg)
	assert.NoError(t, err)

	errClose := s.Close()
	assert.NoError(t, errClose)
}

func TestInMemory_Store(t *testing.T) {
	s := new(InMemory)
	s.Configure(nil)

	err := s.Store(testMessage(1, 2, 3))
	assert.NoError(t, err)
	assert.Equal(t, []byte("1,2,3"), s.mem.Get("0000000000000001:1").Value().(message.Message).Payload)
}

func TestInMemory_QueryLast(t *testing.T) {
	s := newTestMemStore()
	const wildcard = uint32(1815237614)
	tests := []struct {
		query    []uint32
		limit    int
		count    int
		gathered []byte
	}{
		{query: []uint32{0, 10, 20, 50}, limit: 10, count: 0},
		{query: []uint32{0, 1, 1, 1}, limit: 10, count: 1},
		{query: []uint32{0, 1, 1, wildcard}, limit: 10, count: 2},
		{query: []uint32{0, 1}, limit: 10, count: 6},
		{query: []uint32{0, 2}, limit: 10, count: 0},
		{query: []uint32{0, 1, 2}, limit: 10, count: 2},
		{query: []uint32{0, 1}, limit: 5, count: 5},
		{query: []uint32{0, 1}, limit: 5, count: 5, gathered: []byte{0x6a, 0x77, 0xb4, 0x1, 0x7, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x59, 0x31, 0x2c, 0x33, 0x2c, 0x32, 0xb4, 0x2, 0x4, 0x53, 0x73, 0x69, 0x64, 0x40, 0x20, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0xb4, 0x3, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x13, 0x59, 0xaa, 0x23, 0x32, 0x77, 0xb0, 0x1, 0x59, 0x31, 0x2c, 0x33, 0x2c, 0x31, 0xb0, 0x2, 0x40, 0x20, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0xb0, 0x3, 0x13, 0x59, 0xaa, 0x23, 0x32, 0x77, 0xb0, 0x1, 0x59, 0x31, 0x2c, 0x32, 0x2c, 0x32, 0xb0, 0x2, 0x40, 0x20, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0xb0, 0x3, 0x13, 0x59, 0xaa, 0x23, 0x32, 0x77, 0xb0, 0x1, 0x59, 0x31, 0x2c, 0x32, 0x2c, 0x31, 0xb0, 0x2, 0x40, 0x20, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0xb0, 0x3, 0x13, 0x59, 0xaa, 0x23, 0x32, 0x77, 0xb0, 0x1, 0x59, 0x31, 0x2c, 0x31, 0x2c, 0x32, 0xb0, 0x2, 0x40, 0x20, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0xb0, 0x3, 0x13, 0x59, 0xaa, 0x23, 0x32, 0x77, 0xb0, 0x1, 0x59, 0x31, 0x2c, 0x31, 0x2c, 0x31, 0xb0, 0x2, 0x40, 0x20, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0xb0, 0x3, 0x13, 0x59, 0xaa, 0x23, 0x32}},
	}

	for _, tc := range tests {

		if tc.gathered == nil {
			s.Query = nil
		} else {
			s.Query = func(string, []byte) (message.Awaiter, error) {
				return &mockAwaiter{f: func(_ time.Duration) [][]byte { return [][]byte{tc.gathered} }}, nil
			}
		}

		out, err := s.QueryLast(tc.query, tc.limit)
		assert.NoError(t, err)

		count := 0
		for range out {
			count++
		}

		assert.Equal(t, tc.count, count)
	}
}

func TestInMemory_lookup(t *testing.T) {
	s := newTestMemStore()
	const wildcard = uint32(1815237614)
	tests := []struct {
		query []uint32
		limit int
		count int
	}{
		{query: []uint32{0, 10, 20, 50}, limit: 10, count: 0},
		{query: []uint32{0, 1, 1, 1}, limit: 10, count: 1},
		{query: []uint32{0, 1, 1, wildcard}, limit: 10, count: 2},
		{query: []uint32{0, 1}, limit: 10, count: 6},
		{query: []uint32{0, 2}, limit: 10, count: 0},
		{query: []uint32{0, 1, 2}, limit: 10, count: 2},
	}

	for _, tc := range tests {
		matches := s.lookup(lookupQuery{Ssid: tc.query, Limit: tc.limit})
		assert.Equal(t, tc.count, len(matches))
	}
}

func TestInMemory_OnRequest(t *testing.T) {
	s := newTestMemStore()
	tests := []struct {
		name        string
		query       lookupQuery
		expectOk    bool
		expectCount int
	}{
		{name: "dummy"},
		{name: "memstore"},
		{
			name:        "memstore",
			query:       lookupQuery{Ssid: []uint32{0, 1}, Limit: 1},
			expectOk:    true,
			expectCount: 1,
		},
		{
			name:        "memstore",
			query:       lookupQuery{Ssid: []uint32{0, 1}, Limit: 10},
			expectOk:    true,
			expectCount: 6,
		},
	}

	for _, tc := range tests {
		q, _ := utils.Encode(tc.query)
		resp, ok := s.OnRequest(tc.name, q)
		assert.Equal(t, tc.expectOk, ok)
		if tc.expectOk && ok {
			msgs, err := message.DecodeFrame(resp)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectCount, len(msgs))
		}
	}

	// Special, wrong payload case
	_, ok := s.OnRequest("memstore", []byte{})
	assert.Equal(t, false, ok)

}

func Test_param(t *testing.T) {
	raw := `{
	"provider": "memory",
	"config": {
		"maxsize": 99999999
	}
}`
	cfg := testStorageConfig{}
	json.Unmarshal([]byte(raw), &cfg)

	v := param(cfg.Config, "maxsize", 0)
	assert.Equal(t, int64(99999999), v)
}
