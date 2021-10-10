package main

import (
	"container/list"
	"fmt"
	"os"
	"time"

	dapnet "github.com/michaelpeterswa/godapnet"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var settings DapnetNewsSettings
var hc HashCache
var me dapnet.Sender

func (hc *HashCache) buildCache(feed *gofeed.Feed) {
	for _, x := range feed.Items {
		hash := Hash{
			Hash: createHash(x.Title, x.PublishedParsed),
			Time: x.PublishedParsed,
		}
		hc.Set(hash)
	}
}

func (hc *HashCache) checkCurrentArticles(p *gofeed.Parser) error {
	queue := list.New()

	feed, err := p.ParseURL(settings.NewsEndpoint)
	if err != nil {
		logger.Error("Couldn't create feed object from URL", zap.Error(err))
		return err
	}

	for _, article := range feed.Items {
		articleHash := createHash(article.Title, article.PublishedParsed)
		if !hc.Exists(articleHash) {
			hash := Hash{
				Hash: articleHash,
				Time: article.PublishedParsed,
			}
			hc.Set(hash)
			queue.PushBack(article.Title)
		}
	}

	ticker := time.NewTicker(time.Duration(settings.DeliveryDelay) * time.Second)
	for range ticker.C {
		front := queue.Front()
		if front != nil {
			logger.Info("Sending Entry to DAPNet...", zap.Any("entry", front.Value))
			msg := fmt.Sprintf("News - %s", front.Value)
			callsigns := settings.CallsignNames
			txGps := settings.TransmitterGroupNames
			emerg := false

			messages := dapnet.CreateMessage(me.Callsign, msg, callsigns, txGps, emerg)
			payloads := dapnet.GeneratePayload(messages)
			dapnet.SendMessage(payloads, me.Username, me.Password)
			queue.Remove(front)
		} else {
			ticker.Stop()
			logger.Info("Delivery Queue is Empty...")
		}
	}
	return nil
}

func main() {
	fileSettings, err := os.ReadFile("config/settings.yaml")
	if err != nil {
		logger.Fatal("Error loading settings.yaml file")
	}

	err = yaml.Unmarshal(fileSettings, &settings)
	if err != nil {
		logger.Error("YAML failed to unmarshal to DapnetNewsSettings", zap.Error(err))
	}

	me = dapnet.Sender{
		Callsign: settings.DapnetCallsign,
		Username: settings.DapnetUsername,
		Password: settings.DapnetPassword,
	}

	fp := gofeed.NewParser()
	_, err = fp.ParseURL(settings.NewsEndpoint)
	if err != nil {
		logger.Error("Couldn't create feed object from URL", zap.Error(err))
	}

	hc.Init(settings.TTL)

	feed, err := fp.ParseURL(settings.NewsEndpoint)
	if err != nil {
		logger.Error("Couldn't create feed object from URL", zap.Error(err))
	} else {
		hc.buildCache(feed)
	}

	c := cron.New()
	c.AddFunc(settings.CheckInterval, func() {
		err := hc.checkCurrentArticles(fp)
		if err != nil {
			logger.Info("checkCurrentArticles() failed... see above")
		}
	})
	c.AddFunc(settings.CleanInterval, func() { hc.Clean() })
	c.Start()
	select {}
}
