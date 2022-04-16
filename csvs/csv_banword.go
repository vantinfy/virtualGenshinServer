package csvs

import "fmt"

type ConfigBanWord struct {
	Id  int
	Txt string
}

var ConfigBanWordSlice []*ConfigBanWord

func init() {
	// load base csv

	ConfigBanWordSlice = append(ConfigBanWordSlice,
		&ConfigBanWord{1, "外挂"},
		&ConfigBanWord{2, "微信"},
		&ConfigBanWord{3, "+v"},
		&ConfigBanWord{4, "收号"},
		&ConfigBanWord{5, "出号"},
	)
	fmt.Println("csv ban_word initialized")
}

func GetBanWordBase() []string {
	slice := make([]string, 0)
	for _, word := range ConfigBanWordSlice {
		slice = append(slice, word.Txt)
	}
	return slice
}
