package main

import (
	"encoding/json"
	"fmt"
	"github.com/tranvictor/labmeasure"
	"io/ioutil"
)

type History struct {
	Examined           int
	Correct            int
	Incorrect          int
	Accuracy           float32
	AveragePrecision   float32
	AverageRecall      float32
	PrecisionThreshold float32
	RecallThreshold    float32
	ID                 string
}

func (h *History) RecordFromStat(stat labmeasure.Stat) {
	h.Examined = stat.Examined
	h.Correct = stat.Correct
	h.Incorrect = stat.Incorrect
	h.Accuracy = stat.Accuracy()
	h.AveragePrecision = stat.AveragePrecision()
	h.AverageRecall = stat.AverageRecall()
	h.PrecisionThreshold = stat.PrecisionThreshold()
	h.RecallThreshold = stat.RecallThreshold()
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
