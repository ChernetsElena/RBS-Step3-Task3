package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

//var c = make(chan string, 15)
var logFlag *bool
var wg sync.WaitGroup

func main() {
	var (
		datafile *string
		dir      *string
	)

	// чтение аргументов
	datafile = flag.String("datafile", "urls.txt", `Path to datafile."`)
	dir = flag.String("dir", "dir", `Path to dir."`)
	logFlag = flag.Bool("log", false, `Write logs to file."`)

	flag.Parse()

	// открытие файла
	urlFile, err := os.Open(*datafile)
	if err != nil {
		writeError(err)
		os.Exit(1)
	}
	defer urlFile.Close()
	writeInfo("Открытие файла: " + *datafile)

	// создание директории
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		os.MkdirAll(*dir, 0777)
	}

	// чтение файла
	writeInfo("Чтение файла: " + *datafile)
	scanner := bufio.NewScanner(urlFile)

	// узнаем кол-во строк файла, указываем сколько горутин в группе надо ждать (кол-во строк)
	for scanner.Scan() {
		wg.Add(1)
	}

	// перемещаем курсор в начало файла
	urlFile.Seek(0, 0)

	scanner = bufio.NewScanner(urlFile)

	start := time.Now()
	for scanner.Scan() {
		go writeToFileResponce(scanner.Text(), *dir)
	}
	elapsedTime := time.Since(start)

	// for j := 0; j < i; j++ {
	// 	_ = <-c
	// }

	if err := scanner.Err(); err != nil {
		writeError(err)
	}

	wg.Wait()

	fmt.Println("Total Time: " + elapsedTime.String())
}

func writeToFileResponce(address string, dir string) {
	defer wg.Done()
	body := MakeRequest(address)

	fileName := strings.Replace(address, "https://", "", -1)
	fileName = strings.Replace(fileName, "http://", "", -1)
	fileName = strings.Replace(fileName, "/", ".", -1)

	// создание файла
	filePath := path.Join(dir, fileName+".html")
	file, err := os.Create(filePath)
	if err != nil {
		writeError(err)
	}
	defer file.Close()
	writeInfo("Создание файла: " + filePath)

	// запись в файл
	file.Write(body)
	writeInfo("Запись в файл: " + filePath)
	//c <- fileName

}

// функция отправляет запрос и получает данные
func MakeRequest(
	address string) (body []byte) {

	resp, err := http.Get(address)
	if r := recover(); r != nil {
		writeError(err)
	}

	writeInfo("Отправка GET запроса на адрес: " + address)

	if resp == nil {
		writeError(errors.New("Resp is nil"))
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if r := recover(); r != nil {
		writeError(err)
	}
	defer resp.Body.Close()

	writeInfo("Получение данных с адреса: " + address)

	return
}

func writeInfo(infoMessage string) {
	if *logFlag {
		// создание файла для логов
		logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()
		infoLogFile := log.New(logFile, "INFO\t", log.Ldate|log.Ltime)
		infoLogFile.Printf(infoMessage)
	}
	log.Println("INFO\t" + infoMessage)
}

func writeError(errorMessage error) {
	if *logFlag {
		logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()
		errorLogFile := log.New(logFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		errorLogFile.Println(errorMessage)
	}
	log.Println(errors.New("ERROR\t"), errorMessage)
}
