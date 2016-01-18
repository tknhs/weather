package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nlopes/slack"
)

type Slack struct {
	Api      *slack.Client
	Channel  string
	Filename string
	NDaysAgo int
	TeamId   string
	UserId   string
}

func (s *Slack) getFileList(date string) ([]slack.File, error) {
	loc := time.FixedZone("Asia/Tokyo", 9*60*60)
	timeTo, err := time.ParseInLocation("20060102", date[0:8], loc)
	if err != nil {
		return nil, err
	}

	var fileArray []slack.File
	nDaysAgo := time.Duration(-24 * s.NDaysAgo)
	timeFrom := timeTo.Add(nDaysAgo * time.Hour)
	page := 1
	for {
		params := slack.GetFilesParameters{
			User:          s.UserId,
			TimestampFrom: slack.JSONTime(int(timeFrom.Unix())),
			TimestampTo:   slack.JSONTime(int(timeTo.Unix())),
			Types:         "GIF",
			Page:          page,
		}
		files, p, err := s.Api.GetFiles(params)
		if err != nil {
			return nil, err
		}

		if len(files) == 0 {
			break
		} else {
			fileArray = append(fileArray, files...)
		}

		if p.Page >= p.Pages {
			break
		} else {
			page += 1
		}
	}
	return fileArray, nil
}

func (s *Slack) UploadFile(date string) error {
	filename := s.Filename
	params := slack.FileUploadParameters{
		Title:          filename,
		Filetype:       "gif",
		File:           filename + ".gif",
		InitialComment: date,
		Channels:       []string{s.Channel},
	}
	_, err := s.Api.UploadFile(params)
	if err != nil {
		return err
	}
	return nil
}

func (s *Slack) DeleteFiles(date string) error {
	files, err := s.getFileList(date)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := s.Api.DeleteFile(file.ID); err != nil {
			return err
		}
		log.Println(fmt.Sprint(file.InitialComment) + ": Deleted.")
	}
	return nil
}

func NewSlack(config *Config) (*Slack, error) {
	slackApi := slack.New(config.Slack.Token)
	resp, err := slackApi.AuthTest()
	if err != nil {
		return nil, err
	}

	s := &Slack{
		Api:      slackApi,
		Channel:  config.Slack.Channel,
		Filename: config.General.Filename,
		NDaysAgo: config.Slack.NDaysAgo,
		TeamId:   resp.TeamID,
		UserId:   resp.UserID,
	}
	return s, nil
}
