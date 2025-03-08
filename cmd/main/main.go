package main

import (
	"fmt"

	"github.com/dorik33/TgBot/internal/api"
)

func main() {
	data, _ := api.GetInfo()
	fmt.Println(string(data))
}
