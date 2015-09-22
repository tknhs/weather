package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Yahoo   YahooConfig   `toml:"Yahoo"`
	Slack   SlackConfig   `toml:"Slack"`
	General GeneralConfig `toml:"General"`
}
type YahooConfig struct {
	Token string `toml:"token"`
}

type SlackConfig struct {
	Token    string `toml:"token"`
	Channel  string `toml:"channel"`
	NDaysAgo int    `toml:"n_days_ago"`
}

type GeneralConfig struct {
	Filename string `toml:"filename"`
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

func (c *Config) LoadConfig() error {
	_, err := toml.DecodeFile(configFile, &c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) CreateConfig() error {
	var yahooConfig YahooConfig
	fmt.Print("Input YahooToken: ")
	yahooConfig.Token = c.scan()

	var slackConfig SlackConfig
	fmt.Print("Input SlackToken: ")
	slackConfig.Token = c.scan()
	fmt.Print("Input SlackChannel: ")
	slackConfig.Channel = c.scan()
	fmt.Print("Input SlackNDaysAgo: ")
	nDaysAgo, err := strconv.Atoi(c.scan())
	if err != nil {
		return err
	}
	slackConfig.NDaysAgo = nDaysAgo

	var generalConfig GeneralConfig
	fmt.Print("Input Upload Filename (.gif): ")
	generalConfig.Filename = c.scan()

	c.Yahoo = yahooConfig
	c.Slack = slackConfig
	c.General = generalConfig

	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(c); err != nil {
		return err
	}
	ioutil.WriteFile(configFile, []byte(buffer.String()), 0644)
	return nil
}

func NewConfig() (*Config, error) {
	c := &Config{}
	if isFile := c.FileExists(); isFile {
		if err := c.LoadConfig(); err != nil {
			return nil, err
		}
	} else {
		if err := c.CreateConfig(); err != nil {
			return nil, err
		}
		if err := c.LoadConfig(); err != nil {
			return nil, err
		}
	}
	return c, nil
}
