package handler

import (
	"fmt"
	"os"
)

func СhoiceAction(action int) {
	switch action {
	case 1:
		getBookByID()
	case 2:
		addBook()
	case 3:
		updateBook()
	case 4:
		deleteBook()
	case 0:
		fmt.Println("Выход из программы.")
		os.Exit(0)
	default:
		fmt.Println("Неверная команда:", action)
	}

}
