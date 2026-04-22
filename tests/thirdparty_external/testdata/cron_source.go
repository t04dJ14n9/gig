package main

import "github.com/robfig/cron/v3"

func CronNew() string {
	c := cron.New()
	if c == nil {
		return "nil"
	}
	return "ok"
}

func CronParserParse() string {
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sched, err := p.Parse("0 0 * * *")
	if err != nil {
		return "ERR"
	}
	_ = sched
	return "ok"
}
