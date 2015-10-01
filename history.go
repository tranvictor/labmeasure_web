package main

import (
	"encoding/json"
	"fmt"
	"github.com/tranvictor/labmeasure"
	"io/ioutil"
)

type History struct {
	BodyExamined           int
	BodyCorrect            int
	BodyIncorrect          int
	BodyAccuracy           float32
	BodyAveragePrecision   float32
	BodyAverageRecall      float32
	BodyPrecisionThreshold float32
	BodyRecallThreshold    float32
	TitleExamined          int
	TitleCorrect           int
	TitleIncorrect         int
	TitleAccuracy          float32
	ImageExamined          int
	ImageQualified         int
	ImageCorrect           int
	ImageIncorrect         int
	ImageAccuracy          float32
	ID                     string
}

func (h *History) RecordFromStat(stat labmeasure.FinalStat) {
	bodyStats := stat.Stats()["BodyComparer"].(labmeasure.BodyStat)
	h.BodyExamined = bodyStats.Examined
	h.BodyCorrect = bodyStats.Correct
	h.BodyIncorrect = bodyStats.Incorrect
	h.BodyAccuracy = bodyStats.Accuracy()
	h.BodyAveragePrecision = bodyStats.AveragePrecision()
	h.BodyAverageRecall = bodyStats.AverageRecall()
	h.BodyPrecisionThreshold = bodyStats.Configuration.PrecisionThreshold
	h.BodyRecallThreshold = bodyStats.Configuration.RecallThreshold
	titleStats := stat.Stats()["TitleComparer"].(labmeasure.TitleStat)
	h.TitleExamined = titleStats.Examined
	h.TitleCorrect = titleStats.Correct
	h.TitleIncorrect = titleStats.Incorrect
	h.TitleAccuracy = titleStats.Accuracy()
	imageStats := stat.Stats()["ImageComparer"].(labmeasure.ImageStat)
	h.ImageExamined = imageStats.Examined
	h.ImageCorrect = imageStats.Correct
	h.ImageQualified = imageStats.Qualified
	h.ImageIncorrect = imageStats.Incorrect
	h.ImageAccuracy = imageStats.Accuracy()
}

func (h History) WriteToJson(filename string) error {
	b, _ := json.MarshalIndent(h, "", "  ")
	err := ioutil.WriteFile(filename, b, 0666)
	return err
}

func LoadFromJson(filename string) (*History, error) {
	jsstring, err := ioutil.ReadFile("results/" + filename)
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}
	json_bytes := []byte(jsstring)
	var history History
	json.Unmarshal(json_bytes, &history)
	return &history, nil
}
