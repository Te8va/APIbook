package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/Te8va/APIbook/internal/app/client/domain"
)

func getBookByID() {
	fmt.Print("Введите ID книги: ")

	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании ID книги:", sc.Err())
		return
	}
	id := sc.Text()

	response, err := http.Get(domain.BaseURL + "/" + id)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Ошибка при создании запроса к сереверу:", response.Status)
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

	fmt.Print("Введите название книги:")

	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании названия книги:", sc.Err())
		return
	}
	title := sc.Text()

	fmt.Print("Введите автора книги: ")
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании автора книги:", sc.Err())
		return
	}
	author := sc.Text()

	var year int
	var err error
	for {
		fmt.Print("Введите год книги: ")
		if !sc.Scan() {
			fmt.Println("Ошибка при считывании года книги:", sc.Err())
			return
		}

		year, err = strconv.Atoi(sc.Text())
		if err != nil {
			fmt.Println("Ошибка при преобразовании года книги в число:", err)
			continue
		}
		break
	}

	book := domain.Book{
		Title:  title,
		Author: author,
		Year:   year,
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	if err := encoder.Encode(book); err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}

	response, err := http.Post(domain.BaseURL, "application/json", &buf)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		fmt.Println("Ошибка при создании запроса к сереверу:", response.Status)
		return
	}

	fmt.Println("Книга успешно создана:", response.Status)
}

func updateBook() {
	fmt.Print("Введите ID книги: ")
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании ID книги:", sc.Err())
		return
	}
	id := sc.Text()

	fmt.Print("Введите название книги:")
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании названия книги:", sc.Err())
		return
	}
	title := sc.Text()

	fmt.Print("Введите автора книги:")
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании автора книги:", sc.Err())
		return
	}
	author := sc.Text()

	var year int
	var err error
	for {
		fmt.Print("Введите год книги: ")
		if !sc.Scan() {
			fmt.Println("Ошибка при считывании года книги:", sc.Err())
			return
		}

		var err error

		year, err = strconv.Atoi(sc.Text())
		if err != nil {
			fmt.Println("Ошибка при преобразовании года книги в число:", err)
			continue
		}
		break
	}

	book := domain.Book{
		Title:  title,
		Author: author,
		Year:   year,
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	if err := encoder.Encode(book); err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}

	request, err := http.NewRequest(http.MethodPut, domain.BaseURL+"/"+id, &buf)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		fmt.Println("Ошибка при создании запроса к сереверу:", response.Status)
		return
	}

	fmt.Println("Книга успешно обновлена:", response.Status)
}

func deleteBook() {
	fmt.Print("Введите ID книги:")

	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		fmt.Println("Ошибка при считывании ID книги:", sc.Err())
		return
	}
	id := sc.Text()

	request, err := http.NewRequest(http.MethodDelete, domain.BaseURL+"/"+id, nil)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		fmt.Println("Ошибка при создании запроса к сереверу:", response.Status)
		return
	}

	fmt.Println("Книга успешно удалена:", response.Status)
}
