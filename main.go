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
	w := NewWeather()

	// 指定した分の日付リストを作成する
	dateArray := CreateDateArray(time.Now(), 7)

	// Gifアニメ作成
	w.CreateGifImage(dateArray)

	// Slackに投稿
	if err := w.PostSlack(dateArray[0]); err != nil {
		log.Fatal(err)
	}
}
