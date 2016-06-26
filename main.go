// autoimage project main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/pkg/errors"
)

var (
	Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type BotConfig struct {
	Token             string `json:"Token"`
	Channel           string `json:"Channel"`
	ChangeDirInterval uint   `json:"ChangeDirInterval"`
	Root              string `json:"Root"`
	Repeat            bool   `json:"Repeat"`
	//in minutes
	PostingInterval uint `json:"PostInterval"`
	//in minutes
}

func main() {
	ConfigFileName := ""
	flag.StringVar(&ConfigFileName, "c", "config.txt", "path to configuration file")
	flag.Parse()
	ConfigBytes, err := ioutil.ReadFile(ConfigFileName)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error while opening configuration file"))
		os.Exit(1)
	}
	Config := []BotConfig{}
	err = json.Unmarshal(ConfigBytes, &Config)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error while parsing configuration file"))
		os.Exit(1)
	}
	if len(Config) == 0 {
		fmt.Println("No bots defined!\n")
		os.Exit(1)
	}
	logchan := make(chan error, len(Config))
	for _, bc := range Config {
		StartBot(bc, logchan)
	}
	for l := range logchan {
		fmt.Println(l)
	}
}
