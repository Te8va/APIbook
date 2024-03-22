package main

import (
	"fmt"

	"github.com/Te8va/APIbook/internal/app/client/handler"
)

func main() {
	for {
		action, err := handler.ReadUserChoice()
		if err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}
		handler.СhoiceAction(action)
	}
}
