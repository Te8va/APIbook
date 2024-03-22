package handler

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func ReadUserChoice() (int, error) {

	var action int

	fmt.Println("Выберите команду: ")
	fmt.Println("0 - Выход из программы")
	fmt.Println("1 - Получить книгу по ID")
	fmt.Println("2 - Добавить книгу")
	fmt.Println("3 - Обновить книгу")
	fmt.Println("4 - Удалить книгу")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return -1, fmt.Errorf("ошибка при считывании ввода: %v", scanner.Err())
	}

	action, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return -1, fmt.Errorf("ошибка при преобразовании введенной команды в число: %v", err)
	}

	return action, nil
}
