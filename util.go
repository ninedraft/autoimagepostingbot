package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetImages(root string) ([]string, error) {
	ext := map[string]bool{"jpg": true, "jpeg": true, "png": true, "bmp": true}
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	Images := []string{}
	for _, f := range files {
		switch {
		case f.IsDir():
			new_images, err := GetImages(root + "/" + f.Name())
			if err != nil {
				return nil, err
			}
			Images = append(Images, new_images...)
		case !f.IsDir() && ext[strings.ToLower(filepath.Ext(f.Name()))]:
			Images = append(Images, root+"/"+f.Name())
		}
	}
	if len(Images) == 0 {
		return nil, fmt.Errorf("Can't find images in %q\n", root)
	}
	return Images, err
}

func StartBot(config BotConfig, log chan error) error {
	log <- fmt.Errorf("Creating new bot")
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return err
	}
	log <- fmt.Errorf("Collecting images")
	Images, err := GetImages(config.Root)
	selecter := NewImageSelecter(Images, config.Repeat)
	if err != nil {
		return err
	}
	log <- fmt.Errorf("%d images found", len(Images))
	image, ok := selecter.Next()
	if !ok {
		log <- fmt.Errorf("Can't get image")
	}
	photo := tgbotapi.NewPhotoUpload(0, image)
	photo.ChannelUsername = config.Channel
	_, err = bot.Send(photo)
	if err != nil {
		log <- err
	}
	log <- fmt.Errorf("Starting new posting cycle")
	go func() {
		log <- fmt.Errorf("Starting new bot at channel %q", config.Channel)
		dur := time.Duration(config.PostingInterval * uint(time.Second))
		var image string
		var ok bool
		for range time.Tick(dur) {
			time.Sleep(time.Second * time.Duration(Rand.Intn(5)))
			image, ok = selecter.Next()
			photo = tgbotapi.NewPhotoUpload(0, image)
			photo.ChannelUsername = config.Channel
			_, err = bot.Send(photo)
			if err != nil {
				log <- err
			}
			if !ok {
				log <- fmt.Errorf("End image loop")
				return
			}
		}
	}()
	return nil
}
