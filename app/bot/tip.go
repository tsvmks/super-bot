package bot

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/pkg/errors"
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
func NewTip(dataLocation string) (*Sys, error) {
	log.Printf("[INFO] created sys bot, data location=%s", dataLocation)
	res := Sys{dataLocation: dataLocation}
	if err := res.loadBasicData(); err != nil {
		return nil, err
	}
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

func (p *Tip) loadBasicData() error {
	bdata, err := readLines(p.dataLocation + "/basic.data")
	if err != nil {
		return errors.Wrap(err, "can't load basic.data")
	}

	for _, line := range bdata {
		elems := strings.Split(line, "|")
		if len(elems) != 3 {
			log.Printf("[DEBUG] bad format %s, ignored", line)
			continue
		}
		sysCommand := sysCommand{
			description: elems[1],
			message:     elems[2],
			triggers:    strings.Split(elems[0], ";"),
		}
		p.commands = append(p.commands, sysCommand)
		log.Printf("[DEBUG] loaded basic response, %v, %s", sysCommand.triggers, sysCommand.message)
	}
	return nil
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
