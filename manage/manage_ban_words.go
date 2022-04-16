package manage

import (
	"VirtualGenshinServer/csvs"
	"fmt"
	"regexp"
	"time"
)

type ManagesBanWord struct {
	BanWordBase   []string
	BanWordExtend []string
}

var managesBanWord *ManagesBanWord

func GetBanWords() *ManagesBanWord {
	if managesBanWord == nil {
		managesBanWord = new(ManagesBanWord)
		managesBanWord.BanWordBase = []string{"出号", "外挂"}
		managesBanWord.BanWordExtend = []string{"刷单"}
	}
	return managesBanWord
}

// IsBanWord 匹配txt文本是否为禁词
func (w *ManagesBanWord) IsBanWord(txt string) bool {
	for _, s := range w.BanWordBase {
		match, _ := regexp.MatchString(s, txt)
		if match {
			fmt.Println("禁词匹配", s, txt)
			return true
		}
	}

	for _, s := range w.BanWordExtend {
		match, _ := regexp.MatchString(s, txt)
		if match {
			fmt.Println("禁词匹配", s, txt)
			return true
		}
	}
	return false
}

func (w *ManagesBanWord) Run() {
	w.BanWordBase = csvs.GetBanWordBase()

	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			//fmt.Println(time.Now().Unix(), "词库更新")
		}
	}
}
