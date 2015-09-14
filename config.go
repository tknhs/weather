package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Yahoo YahooConfig `toml:"yahoo"`
	Slack SlackConfig `toml:"slack"`
}
type YahooConfig struct {
	Token string `toml:"token"`
}

type SlackConfig struct {
	Token   string `toml:"token"`
	Channel string `toml:"channel"`
}

const configFile = "config.toml"

// 設定ファイルの存在チェック
func (c *Config) FileExists() bool {
	_, err := os.Stat(configFile)
	return err == nil
}

func (c *Config) scan() string {
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	input := sc.Text()
	return input
}

func (c *Config) LoadConfig() {
	_, err := toml.DecodeFile(configFile, &c)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Config) CreateConfig() {
	var yahooConfig YahooConfig
	fmt.Print("Input YahooToken: ")
	textYahooToken := c.scan()
	yahooConfig.Token = textYahooToken

	var slackConfig SlackConfig
	fmt.Print("Input SlackToken: ")
	textSlackToken := c.scan()
	slackConfig.Token = textSlackToken
	fmt.Print("Input SlackChannel: ")
	textSlackChannel := c.scan()
	slackConfig.Channel = textSlackChannel

	c.Yahoo = yahooConfig
	c.Slack = slackConfig

	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(c)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(configFile, []byte(buffer.String()), 0644)
}

func NewConfig() *Config {
	c := &Config{}
	if isFile := c.FileExists(); isFile {
		c.LoadConfig()
	} else {
		c.CreateConfig()
		c.LoadConfig()
	}
	return c
}
