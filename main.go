// autoimage project main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

var (
	Rand        = rand.New(rand.NewSource(time.Now().UnixNano()))
	BotsConfigs = []BotConfig{}
)

type BotConfig struct {
	Token           string `yaml: token`
	Channel         string `yaml: channel`
	Root            string `yaml: root`
	Repeat          bool   `yaml: repeat`
	PostingInterval uint   `yaml: postingterval`
}

func main() {
	ConfigFileName := ""
	flag.StringVar(&ConfigFileName, "c", "config.txt", "path to configuration file")
	flag.Parse()

	ConfigBytes, err := ioutil.ReadFile(ConfigFileName)
	err = yaml.Unmarshal(ConfigBytes, &BotsConfigs)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error while opening configuration file"))
		os.Exit(1)
	}

	if len(BotsConfigs) == 0 {
		fmt.Println("No bots defined!\n")
		os.Exit(1)
	}

	logchan := make(chan error, len(BotsConfigs))
	for _, bc := range BotsConfigs {
		log.Printf("Starting new bot")
		go func(botconf BotConfig) {
			errbot := StartBot(botconf, logchan)
			if errbot != nil {
				logchan <- errbot
			}
		}(bc)
	}

	for l := range logchan {
		fmt.Println(l)
	}
}
