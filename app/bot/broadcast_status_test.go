package bot

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestBroadcast_OnMessage(t *testing.T) {
	tbl := []struct {
		lastSentStatus   string
		status           string
		expectedResponse Response
	}{
		{"Версия 1.0.0", "Версия 1.0.0", Response{}},
		{"Версия 1.0.0", "Версия 2.0.0", Response{Text: "Версия 2.0.0", Send: true}},
		{"Версия 2.0.0", "Версия 3.0.0", Response{Text: "Версия 3.0.0", Send: true, Unpin: false}},
		{"Версия 3.0.0", "Версия 3.0.0", Response{}},
	}

	b := &BroadcastStatus{}

	for i, tt := range tbl {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b.lastSentVersion = tt.lastSentStatus
			b.version = tt.status
			response := b.OnMessage(Message{})

			require.Equal(t, tt.expectedResponse, response)
		})
	}
}
/*
// Kind of integration test to check all workflow
func TestBroadcast_StatusTransitions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	status := true // hold here ping status
	statusMx := sync.Mutex{}
	setStatus := func(s bool) {
		statusMx.Lock()
		defer statusMx.Unlock()
		status = s
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusMx.Lock()
		defer statusMx.Unlock()
		if status {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	b := NewBroadcastStatus(ctx, BroadcastParams{
		URL:          ts.URL,
		PingInterval: time.Millisecond,
		DelayToOff:   100 * time.Millisecond,
		Client:       http.Client{},
	})

	// Test reacts on first message
	require.Equal(t, Response{}, b.OnMessage(Message{}))

	// Test do not react on second message because status not changed
	require.Equal(t, Response{}, b.OnMessage(Message{}))

	// Wait for off->on
	time.Sleep(20 * time.Millisecond)
	require.Equal(t, Response{Text: MsgBroadcastStarted, Send: true, Pin: false}, b.OnMessage(Message{}))
	//require.Equal(t., b.getStatus())

	// off
	setStatus(false)
	// Still on, no deadline reached
	time.Sleep(20 * time.Millisecond)
	require.Equal(t, Response{}, b.OnMessage(Message{}))
	//require.True(t, b.getStatus())

	// Deadline reached on->off
	time.Sleep(110 * time.Millisecond)
	require.Equal(t, Response{Text: MsgBroadcastFinished, Send: true, Unpin: true}, b.OnMessage(Message{}))
	//require.False(t, b.getStatus())
}

func TestBroadcast_StatusOffToOn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := &BroadcastStatus{}
	b.check(ctx, time.Time{}, BroadcastParams{
		URL:    ts.URL,
		Client: http.Client{},
	})

	//require.True(t, b.status)
}

func TestBroadcast_StatusOffToOff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	b := &BroadcastStatus{}
	//b.status = false
	b.check(ctx, time.Time{}, BroadcastParams{
		URL:    ts.URL,
		Client: http.Client{},
	})

	//require.False(t, b.status)
}

func TestBroadcast_StatusOnToOffNoDeadline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	b := &BroadcastStatus{}
	//b.status = true
	b.check(ctx, time.Now(), BroadcastParams{
		URL:        ts.URL,
		DelayToOff: time.Second,
		Client:     http.Client{},
	})

	//require.True(t, b.status)
}

func TestBroadcast_StatusOnToOffWithDeadline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	b := &BroadcastStatus{}
	//b.status = true
	b.check(ctx, time.Now().Add(-2*time.Second), BroadcastParams{
		URL:        ts.URL,
		DelayToOff: time.Second,
		Client:     http.Client{},
	})

	//require.False(t, b.status)
}

func TestBroadcast_FirstOnMessageReturnsCurrentState(t *testing.T) {
	b := &BroadcastStatus{}
	response := b.OnMessage(Message{})
	require.False(t, response.Send)
}

func TestBroadcast_OnMessageReturnsNothingIfStateNotChanged(t *testing.T) {
	b := &BroadcastStatus{}
	response := b.OnMessage(Message{})
	require.False(t, response.Send)

	//b = &BroadcastStatus{status: true, lastSentStatus: true}
	response = b.OnMessage(Message{})
	require.False(t, response.Send)
}

func TestBroadcast_OnMessageReturnsReplyOnChange(t *testing.T) {
	//b := &BroadcastStatus{lastSentStatus: false, status: true} // OFF ->ON
	//resp := b.OnMessage(Message{})
	//require.True(t, resp.Send)
	//require.Equal(t, MsgBroadcastStarted, resp.Text)

	//b = &BroadcastStatus{lastSentStatus: true, status: false} // ON -> OFF
	//resp = b.OnMessage(Message{})
	//require.True(t, resp.Send)
	//require.Equal(t, MsgBroadcastFinished, resp.Text)
}

func TestBroadcast_PingReturnsTrueOn200Status(t *testing.T) {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	//require.True(t, ping(ctx, http.Client{}, ts.URL))
}

func TestBroadcast_PingReturnsFalseOnNot200Status(t *testing.T) {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	//require.False(t, ping(ctx, http.Client{}, ts.URL))
}

func TestBroadcast_PingReturnsFalseOnUnableToDoRequest(t *testing.T) {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//require.False(t, ping(ctx, http.Client{}, "http://localhost:9873"))
}

func TestBroadcast_PingReturnsFalseOnUnableToCreateReq(t *testing.T) {
	//require.False(t, ping(nil, http.Client{}, "http://localhost:9873"))
}

func TestBroadcast_ReactOnAnyMessage(t *testing.T) {
	require.Empty(t, (&BroadcastStatus{}).ReactOn())
}
 */
