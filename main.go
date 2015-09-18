package main

import (
	"log"
	"time"
)

func CreateDateArray(t time.Time, count int) []string {
	dateArray := make([]string, count)
	for i := 0; i < count; i++ {
		after5mins := t.Add(time.Duration(i*5) * time.Minute)
		dateArray[i] = after5mins.Format("200601021504")
	}
	return dateArray
}

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	w := NewWeather(config)
	s, err := NewSlack(config)
	if err != nil {
		log.Fatal(err)
	}

	// 指定した分の日付リストを作成する
	dateArray := CreateDateArray(time.Now(), 7)

	// Gifアニメ作成
	if err := w.CreateGifImage(dateArray); err != nil {
		log.Fatal(err)
	}

	// Slackに投稿
	if err := s.UploadFile(dateArray[0]); err != nil {
		log.Fatal(err)
	}

	// 前日分の投稿画像を削除
	if err := s.DeleteFiles(dateArray[0]); err != nil {
		log.Fatal(err)
	}
}
