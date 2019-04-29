package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

// File containing the credentials for the bot.
const accountFile = "account.txt"

// Subreddits to monitor.
var monitoredSubreddits = []string{
	"gamedeals",
}

// Users to alert when a free game deal is found.
var alertedUsers = []string{
	"YOUR.REDDIT.USERNAME",
}

// Words of interest.
var wantedWords = []string{
	"test01",
	"100% off",
	"100 % off",
	"100%off",
	"100%",
	"100",
	"free",
	"$0",
	"0$",
	"$0.00",
	"$0.0",
	"$0,00",
	"$0,0",
	"0.00$",
	"0.0$",
	"0,00$",
	"0,0$",
	"€0",
	"0€",
	"€0.00",
	"€0.0",
	"€0,00",
	"€0,0",
	"0.00€",
	"0.0€",
	"0,00€",
	"0,0€",
}

// Words to ignore.
var ignoredWords = []string{
	"drm-free",
	"free delivery",
	"plus free",
	"+ free",
	"free gift for redditors",
	"free us shipping",
	"free coin shop game",
	"100% orange juice",
}

// DealsBot checks for free game deals.
type DealsBot struct {
	bot reddit.Bot
}

// Post handles graw post events.
func (b *DealsBot) Post(p *reddit.Post) error {
	if b.isFreeGameDeal(p) {
		return b.sendAlerts(p)
	}
	return nil
}

func (b *DealsBot) isFreeGameDeal(p *reddit.Post) bool {
	title := strings.ToLower(p.Title)

	for _, w := range ignoredWords {
		title = strings.ReplaceAll(title, w, "")
	}

	for _, w := range wantedWords {
		if strings.Contains(title, w) {
			return true
		}
	}

	return false
}

func (b *DealsBot) sendAlerts(p *reddit.Post) error {
	subject := "Free game deal found"
	text := fmt.Sprintf("Title: %v\n\nURL: %v\n\nDomain: %v", p.Title, p.URL, p.Domain)

	log.Printf("\n%v\n%v", subject, text)

	for _, u := range alertedUsers {
		if err := b.bot.SendMessage(u, subject, text); err != nil {
			return err
		}
		time.Sleep(time.Minute)
	}

	return nil
}

func main() {
	log.Print("starting game deals bot")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	bot, err := reddit.NewBotFromAgentFile(accountFile, 15*time.Minute)
	if err != nil {
		return err
	}

	cfg := graw.Config{Subreddits: monitoredSubreddits}
	handler := &DealsBot{bot: bot}
	_, wait, err := graw.Run(handler, bot, cfg)
	if err != nil {
		return err
	}

	err = wait()
	if err != nil {
		return err
	}

	return nil
}
