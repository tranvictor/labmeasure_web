package main

import "github.com/tranvictor/labmeasure"

type StatPresenter struct {
	*labmeasure.Stat
	Filename  string
	Histories HistoryList
}

func (sp StatPresenter) NotNil() bool {
	return sp.Stat != nil
}
