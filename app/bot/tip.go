package bot

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

// Sys implements basic bot function to respond on ping and others from basic.data file.
// also, reacts on say! with keys/values from say.data file
type Tip struct {
	say          []string
	dataLocation string
	commands     []sysCommand
}

// sysCommand hold one type triggers from basic.data
type tipCommand struct {
	triggers    []string
	description string
	message     string
}

// NewSys makes new sys bot and load data to []say and basic map
func NewTip(dataLocation string) (*Tip, error) {
	log.Printf("[INFO] created sys bot, data location=%s", dataLocation)
	res := Tip{dataLocation: dataLocation}
	if err := res.loadSayData(); err != nil {
		return nil, err
	}
	rand.Seed(0)
	return &res, nil
}

// Help returns help message
func (p Tip) Help() (line string) {
	for _, c := range p.commands {
		line += genHelpMsg(c.triggers, c.description)
	}
	return line
}

// OnMessage implements bot.Interface
func (p Tip) OnMessage(msg Message) (response Response) {
	if !contains(p.ReactOn(), msg.Text) {
		return Response{}
	}

	if strings.EqualFold(msg.Text, "tip!") {
		if p.say != nil && len(p.say) > 0 {
			return Response{
				Text: fmt.Sprintf("_%s_", p.say[rand.Intn(len(p.say))]),
				Send: true,
			}
		}
		return Response{}
	}

	for _, bot := range p.commands {
		if found := contains(bot.triggers, strings.ToLower(msg.Text)); found {
			return Response{Text: bot.message, Send: true}
		}
	}

	return Response{}
}

func (p *Tip) loadSayData() error {
	say, err := readLines(p.dataLocation + "/tips.data")
	if err != nil {
		return err
	}
	p.say = say
	log.Printf("[DEBUG] loaded say.data, %d records", len(say))
	return nil
}

// ReactOn keys
func (p Tip) ReactOn() []string {
	res := make([]string, 0)
	for _, bot := range p.commands {
		res = append(bot.triggers, res...)
	}
	return res
}
