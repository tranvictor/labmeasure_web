package main

import (
	"io/ioutil"
	"strings"
)

type HistoryList []*History

func LoadHistories() (HistoryList, error) {
	fileInfos, err := ioutil.ReadDir("results/")
	if err != nil {
		return HistoryList{}, err
	}

	result := make(HistoryList, 0)
	for index := range fileInfos {
		fileInfo := fileInfos[len(fileInfos)-1-index]
		name := fileInfo.Name()
		if fileInfo.IsDir() || !strings.HasSuffix(name, ".json") {
			continue
		}
		if history, err := LoadFromJson(name); err == nil {
			result = append(result, history)
		}
	}
	return result, nil
}
