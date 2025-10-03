package main

import (
	"fmt"

	"github.com/simonjwhitlock/bootdev-go-blogaggregator/internal/config"
)

func main() {
	jsonConf, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	jsonConf.UserName = "Simon"

	err = jsonConf.SetUser()
	if err != nil {
		fmt.Println(err)
	}

	jsonConf, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DB url:", jsonConf.DbUrl)
	fmt.Println("UserName:", jsonConf.UserName)
}
