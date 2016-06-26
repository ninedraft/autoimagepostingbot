package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetSubDirs(root string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	SubDirs := []os.FileInfo{}
	for _, f := range files {
		if f.IsDir() {
			SubDirs = append(SubDirs, f)
		}
	}
	if len(SubDirs) == 0 {
		return nil, fmt.Errorf("Can't find subdirs in root %q\n", root)
	}
	return SubDirs, nil
}

func StartBot(config BotConfig, log chan error) error {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return err
	}
	dirs, err := GetSubDirs(config.Root)
	if err != nil {
		return err
	}
	CurrDir := dirs[Rand.Intn(len(dirs))]
	Selecter, err := NewSelecter(config, CurrDir)
	if err != nil {
		return err
	}
	image, NextExist := Selecter.Next()
	photo := tgbotapi.NewPhotoUpload(0, image)
	photo.ChannelUsername = config.Channel
	_, err = bot.Send(photo)
	if err != nil {
		log <- err
	}
	go func() {
		dur := time.Duration(config.PostingInterval * uint(time.Minute))
		var s *ImageSelecter
		for range time.Tick(dur) {
			image, NextExist = Selecter.Next()
			photo = tgbotapi.NewPhotoUpload(0, image)
			photo.ChannelUsername = config.Channel
			_, err = bot.Send(photo)
			if err != nil {
				log <- err
			}
			if !NextExist {
				CurrDir = dirs[Rand.Intn(len(dirs))]
				s, err = NewSelecter(config, CurrDir)
				if err != nil {
					log <- err
				} else {
					log <- fmt.Errorf("Changing dir to %q by depletion", CurrDir.Name())
					Selecter = s
				}
			}
		}
	}()
	go func() {
		dur := time.Duration(config.ChangeDirInterval * uint(time.Minute))
		var s *ImageSelecter
		for range time.Tick(dur) {
			CurrDir = dirs[Rand.Intn(len(dirs))]
			s, err = NewSelecter(config, CurrDir)
			if err != nil {
				log <- err
			} else {
				log <- fmt.Errorf("Changing dir to %q by timer", CurrDir.Name())
				Selecter = s
			}
		}
	}()
	return nil
}
