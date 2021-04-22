package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	var (
		datafile *string
		dir      *string
	)

	// создание файла для логов
	logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// создание логов
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	infoLogFile := log.New(logFile, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogFile := log.New(logFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// чтение аргументов
	datafile = flag.String("datafile", "urls.txt", `Path to datafile. Default: "urls.txt"`)
	dir = flag.String("dir", "./dir/", `Path to dir. Default: "./dir/"`)
	flag.Parse()

	// открытие файла
	urlFile, err := os.Open(*datafile)
	if err != nil {
		errorLog.Println(err)
		errorLogFile.Println(err)
		os.Exit(1)
		return
	}
	defer urlFile.Close()
	infoLog.Printf("Открытие файла: %s\n", *datafile)
	infoLogFile.Printf("Открытие файла: %s\n", *datafile)

	// создание директории
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		os.MkdirAll(*dir, 0777)
	}

	// чтение файла
	infoLog.Printf("Чтение файла: %s\n", *dir+*datafile)
	infoLogFile.Printf("Чтение файла: %s\n", *dir+*datafile)
	scanner := bufio.NewScanner(urlFile)

	for scanner.Scan() {
		address := string(scanner.Text())
		writeBodyToFile(address, dir, errorLog, infoLog, errorLogFile, infoLogFile)
	}

	if err := scanner.Err(); err != nil {
		errorLog.Println(err)
		errorLogFile.Println(err)
		os.Exit(1)
	}
}

// функция записывает в файл данные ответа
func writeBodyToFile(
	address string,
	dir *string,
	errorLog *log.Logger,
	infoLog *log.Logger,
	errorLogFile *log.Logger,
	infoLogFile *log.Logger) {

	body := MakeRequest(address, errorLog, infoLog, errorLogFile, infoLogFile)

	log.Println("body")

	fileName := strings.Replace(address, "https://", "", -1)
	fileName = strings.Replace(fileName, "http://", "", -1)
	fileName = strings.Replace(fileName, "/", ".", -1)

	// создание файла
	file := MakeFile(dir, fileName, errorLog, infoLog, errorLogFile, infoLogFile)
	log.Println("file")

	// запись в файл
	file.Write(body)
	infoLog.Printf("Запись в файл: %s\n", *dir+fileName+".html")
	infoLogFile.Printf("Запись в файл: %s\n", *dir+fileName+".html")

	return
}

// функция создает файл
func MakeFile(dir *string, fileName string, errorLog *log.Logger, infoLog *log.Logger, errorLogFile *log.Logger, infoLogFile *log.Logger) (file *os.File) {
	file, err := os.Create(*dir + fileName + ".html")
	if err != nil {
		errorLog.Println("Unable to create file:", err)
		errorLogFile.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	infoLog.Printf("Создание файла: %s\n", *dir+fileName+".html")
	infoLogFile.Printf("Создание файла: %s\n", *dir+fileName+".html")

	return
}

// функция отправляет запрос и получает данные
func MakeRequest(address string, errorLog *log.Logger, infoLog *log.Logger, errorLogFile *log.Logger, infoLogFile *log.Logger) (body []byte) {
	resp, err := http.Get(address)
	if err != nil {
		errorLog.Println(err)
		errorLogFile.Println(err)
		os.Exit(1)
	}
	infoLog.Printf("Отправка GET запроса на адрес: %s\n", address)
	infoLogFile.Printf("Отправка GET запроса на адрес: %s\n", address)

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		errorLog.Println(err)
		errorLogFile.Println(err)
		os.Exit(1)
	}
	infoLog.Printf("Получение данных с адреса: %s\n", address)
	infoLogFile.Printf("Получение данных с адреса: %s\n", address)

	return
}
