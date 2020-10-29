package bot

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	// MsgBroadcastStarted defines text to be sent by the bot when the broadcast started
	MsgBroadcastStarted = "Вещание началось. Приобщиться можно тут: https://stream.radio-t.com/"

	// MsgBroadcastFinished defines text to be sent by the bot when the broadcast finished
	MsgBroadcastFinished = "Вещание завершилось"
)

// BroadcastParams defines parameters for broadcast detection
type BroadcastParams struct {
	URL          string        // URL for "ping"
	PingInterval time.Duration // Ping interval
	DelayToOff   time.Duration // State will be switched to off in no ok replies from URL in this intrval
	Client       http.Client   // http client
}

// BroadcastStatus bot replies with current broadcast status
type BroadcastStatus struct {
	version         string // current broadcast status
	lastSentVersion string // last status sent with OnMessage
	statusMx       sync.Mutex
}

// NewBroadcastStatus starts status checking goroutine and returns bot instance
func NewBroadcastStatus(ctx context.Context, params BroadcastParams) *BroadcastStatus {
	log.Printf("[INFO] BroadcastStatus bot with %v", params.URL)
	b := &BroadcastStatus{}
	go b.checker(ctx, params)
	return b
}

// Help returns help message
func (b *BroadcastStatus) Help() string {
	return ""
}

// OnMessage returns current broadcast status if it was changed
func (b *BroadcastStatus) OnMessage(_ Message) (response Response) {
	b.statusMx.Lock()
	defer b.statusMx.Unlock()

	response.Pin = false
	if b.lastSentVersion != b.version {
		response.Send = true
		response.Text = b.version
		b.lastSentVersion = b.version
	}
	return
}

func (b *BroadcastStatus) checker(ctx context.Context, params BroadcastParams) {
	lastOn := time.Time{}
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(params.PingInterval):
			lastOn = b.check(ctx, lastOn, params)
		}
	}
}

// check do ping to url and change current state
func (b *BroadcastStatus) check(ctx context.Context, lastOn time.Time, params BroadcastParams) time.Time {
	b.statusMx.Lock()
	defer b.statusMx.Unlock()

	newStatus := ping(ctx, params.Client, params.URL)

	b.version = newStatus

	// 1 -> 1
	return time.Now()
}

// ping do get request to https://stream.radio-t.com and returns true on OK status and false for all other statuses
func ping(ctx context.Context, client http.Client, url string) (status string) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("[WARN] unable to create %v request, %v", url, err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[DEBUG] unable to execute %v request, %v", url, err)
		return
	}
	defer resp.Body.Close()

	responseData,err  := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	status = string(responseData)
	return
}

// nolint
func (b *BroadcastStatus) getStatus() string {
	b.statusMx.Lock()
	defer b.statusMx.Unlock()
	return b.version
}

// ReactOn keys
func (b *BroadcastStatus) ReactOn() []string {
	return []string{}
}
