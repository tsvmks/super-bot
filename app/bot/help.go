package bot

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

// Sys implements basic bot function to respond on ping and others from basic.data file.
// also, reacts on say! with keys/values from say.data file
type Help struct {
	say          []string
	dataLocation string
	commands     []sysCommand
}

// sysCommand hold one type triggers from basic.data
type helpCommand struct {
	triggers    []string
	description string
	message     string
}

// NewSys makes new sys bot and load data to []say and basic map
func NewHelp(dataLocation string) (*Help, error) {
	log.Printf("[INFO] created sys bot, data location=%s", dataLocation)
	res := Help{dataLocation: dataLocation}
	if err := res.loadSayData(); err != nil {
		return nil, err
	}
	rand.Seed(0)
	return &res, nil
}

// Help returns help message
func (p Help) Help() (line string) {
	for _, c := range p.commands {
		line += genHelpMsg(c.triggers, c.description)
	}
	return line
}

// OnMessage implements bot.Interface
func (p Help) OnMessage(msg Message) (response Response) {
	if strings.EqualFold(msg.Text, "help!") {
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

func (p *Help) loadSayData() error {
	say, err := readLines(p.dataLocation + "/helps.data")
	if err != nil {
		return err
	}
	p.say = say
	log.Printf("[DEBUG] loaded say.data, %d records", len(say))
	return nil
}

// ReactOn keys
func (p Help) ReactOn() []string {
	res := make([]string, 0)
	for _, bot := range p.commands {
		res = append(bot.triggers, res...)
	}
	return res
}
