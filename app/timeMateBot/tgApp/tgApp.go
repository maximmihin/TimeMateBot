package tgApp

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

type TgTimeMate struct {
	Bot    *tgbotapi.BotAPI
	Config *tgbotapi.UpdateConfig
}

func (t *TgTimeMate) StartPolling() tgbotapi.UpdatesChannel {
	return t.Bot.GetUpdatesChan(*t.Config)
}

type Config struct {
	TgToken    string
	UpdateConf *tgbotapi.UpdateConfig
	DebugMode  bool
}

type configBuilder struct {
	tgToken    string
	updateConf *tgbotapi.UpdateConfig
	debugMode  bool
}

func NewConfigBuilder() *configBuilder {
	return &configBuilder{}
}

func (c *configBuilder) UpdateConfig(UpdateConf *tgbotapi.UpdateConfig) *configBuilder {
	c.updateConf = UpdateConf
	return c
}

func (c *configBuilder) BotToken(botToken string) *configBuilder {
	c.tgToken = botToken
	return c
}

func (c *configBuilder) DebugMode(debugMode bool) *configBuilder {
	c.debugMode = debugMode
	return c
}

func (c *configBuilder) Build() *Config {
	conf := &Config{}

	conf.TgToken = c.tgToken

	if c.updateConf == nil {
		conf.UpdateConf = &tgbotapi.UpdateConfig{
			Offset:  0,
			Limit:   0,
			Timeout: 60,
			AllowedUpdates: []string{
				"message",
				"channel_post",
				"edited_channel_post",
				"callback_query",
			},
		}
	} else {
		conf.UpdateConf = c.updateConf
	}

	conf.DebugMode = c.debugMode

	return conf

}

func NewApp(config *Config) (*TgTimeMate, error) {
	TimeMate := TgTimeMate{
		Config: config.UpdateConf,
	}

	var err error
	TimeMate.Bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	TimeMate.Bot.Debug = config.DebugMode
	return &TimeMate, nil
}
