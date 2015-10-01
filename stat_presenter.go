package main

type StatPresenter struct {
	LastHistory *History
	Filename    string
	Histories   HistoryList
}

func (sp StatPresenter) NotNil() bool {
	return sp.LastHistory != nil
}
