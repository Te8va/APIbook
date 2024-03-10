package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const BaseURL = "http://localhost:8080/books"

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"data"`
}

func main() {
	for {
		var action int

		fmt.Println("Выберите команду: ")
		fmt.Println("0 - Выход из программы")
		fmt.Println("1 - Получить книгу по ID")
		fmt.Println("2 - Добавить книгу")
		fmt.Println("3 - Обновить книгу")
		fmt.Println("4 - Удалить книгу")

		sc := bufio.NewScanner(os.Stdin)
		if !sc.Scan() {
			fmt.Println("Ошибка при считывании ввода:", sc.Err())
			return
		}

		action, err := strconv.Atoi(sc.Text())
		if err != nil {
			fmt.Println("Ошибка при преобразовании введенной команды в число:", err)
			return
		}

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

}

func getBookByID() {
	var id string

	fmt.Print("Введите ID книги: ")

	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании ID книги:", sc.Err())
		return
	}
	id = sc.Text()

	response, err := http.Get(BaseURL + "/" + id)
	if err != nil {
		fmt.Println("Ошибка при выполнении GET запроса:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Ошибка: Неверный статус ответа:", response.Status)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	fmt.Println("Ответ от сервера:", string(body))

}

func addBook() {
	var title, author string
	var year int

	fmt.Print("Введите название книги: ")

	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании названия книги:", sc.Err())
		return
	}
	title = sc.Text()

	fmt.Print("Введите автора книги: ")
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании автора книги:", sc.Err())
		return
	}
	author = sc.Text()

	fmt.Print("Введите год книги: ")
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании года книги:", sc.Err())
		return
	}

	year, err := strconv.Atoi(sc.Text())
	if err != nil {
		fmt.Println("Ошибка при преобразовании года книги в число:", err)
		return
	}

	book := Book{
		Title:  title,
		Author: author,
		Year:   year,
	}

	jsonData, err := json.Marshal(book)
	if err != nil {
		fmt.Println("Ошибка кодирования JSON:", err)
		return
	}

	response, err := http.Post(BaseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при выполнении POST запроса:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		fmt.Println("Ошибка: Неверный статус ответа при создании книги:", response.Status)
		return
	}

	fmt.Println("Статус создания книги:", response.Status)
}

func updateBook() {
	var title, author, id string
	var year int

	fmt.Print("Введите ID книги: ")
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании ID книги:", sc.Err())
		return
	}
	id = sc.Text()

	fmt.Print("Введите название книги: ")
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании названия книги:", sc.Err())
		return
	}
	title = sc.Text()

	fmt.Print("Введите автора книги: ")
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании автора книги:", sc.Err())
		return
	}
	author = sc.Text()

	fmt.Print("Введите год книги: ")
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании года книги:", sc.Err())
		return
	}

	year, err := strconv.Atoi(sc.Text())
	if err != nil {
		fmt.Println("Ошибка при преобразовании года книги в число:", err)
		return
	}

	book := Book{
		Title:  title,
		Author: author,
		Year:   year,
	}

	jsonData, err := json.Marshal(book)
	if err != nil {
		fmt.Println("Ошибка при кодирования JSON", err)
	}
	request, err := http.NewRequest(http.MethodPut, BaseURL+"/"+id, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при создании PUT запроса:", err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка при выполнении PUT запроса:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		fmt.Println("Ошибка: Неверный статус ответа при обновлении книги:", response.Status)
		return
	}

	fmt.Println("Статус обновления книги:", response.Status)
}

func deleteBook() {
	var id string

	fmt.Print("Введите ID книги: ")

	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании ID книги:", sc.Err())
		return
	}
	id = sc.Text()

	request, err := http.NewRequest(http.MethodDelete, BaseURL+"/"+id, nil)
	if err != nil {
		fmt.Println("Ошибка при создании Delete запроса:", err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка при выполнении Delete запроса:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		fmt.Println("Ошибка: Неверный статус ответа при создании книги:", response.Status)
		return
	}

	fmt.Println("Статус удаления книги:", response.Status)
}
