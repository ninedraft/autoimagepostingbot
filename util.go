package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetImages(root string) ([]string, error) {
	ext := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".bmp": true}
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
	defer func() {
		er, is_error := recover().(error)
		if is_error {
			log <- fmt.Errorf("%s: FATAL ERROR\n\t%v", config.Name, er)
		}
	}()
	BotRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	log <- fmt.Errorf("%s: creating new bot\n\tposting interval %ds\n\troot:%q", config.Name, config.PostingInterval, config.Root)
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return err
	}
	log <- fmt.Errorf("%s: collecting images", config.Name)
	Images, err := GetImages(config.Root)
	if err != nil {
		return err
	}
	selecter := NewImageSelecter(Images, config.Repeat)
	log <- fmt.Errorf("%s: %d images found", config.Name, len(Images))
	image, ok := selecter.Next()
	if !ok {
		log <- fmt.Errorf("%s: Can't get image", config.Name)
	}
	photo := tgbotapi.NewPhotoUpload(0, image)
	photo.ChannelUsername = config.Channel
	time.Sleep(time.Second * time.Duration(BotRand.Intn(5)))
	_, err = bot.Send(photo)
	if err != nil {
		log <- err
	}
	go func() {
		log <- fmt.Errorf("%s: starting new bot at channel %q", config.Name, config.Channel)
		dur := time.Duration(config.PostingInterval * uint(time.Second))
		log <- fmt.Errorf("%s: posting insterval %s", config.Name, dur)
		var image string
		var ok bool
		for range time.Tick(dur) {
			time.Sleep(time.Second * time.Duration(BotRand.Intn(5)))
			image, ok = selecter.Next()
			photo = tgbotapi.NewPhotoUpload(0, image)
			photo.ChannelUsername = config.Channel
			_, err = bot.Send(photo)
			if err != nil {
				log <- err
			}
			if !ok {
				log <- fmt.Errorf("%s: end image loop", config.Name)
				return
			}
			log <- fmt.Errorf("%s: img %q", config.Name, filepath.Base(image))
		}
	}()
	return nil
}
