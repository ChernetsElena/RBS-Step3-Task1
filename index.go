package main

import (
	"bufio"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	// чтение аргументов
	var (
		datafile *string
		dir      *string
		logFlag  *bool
	)

	datafile = flag.String("datafile", "urls.txt", `Path to datafile."`)
	dir = flag.String("dir", "dir", `Path to dir."`)
	logFlag = flag.Bool("log", false, `Write logs to file."`)

	flag.Parse()

	// открытие файла
	urlFile, err := os.Open(*datafile)
	if err != nil {
		writeError(err, *logFlag)
		os.Exit(1)
	}
	defer urlFile.Close()
	writeInfo("Открытие файла: "+*datafile, *logFlag)

	// создание директории
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		os.MkdirAll(*dir, 0777)
	}

	// чтение файла
	writeInfo("Чтение файла: "+*datafile, *logFlag)
	scanner := bufio.NewScanner(urlFile)
	//log.Println(scanner)

	for scanner.Scan() {
		//log.Println(scanner)
		//log.Println(scanner.Scan())
		address := string(scanner.Text())
		body := MakeRequest(address, *logFlag)

		fileName := strings.Replace(address, "https://", "", -1)
		fileName = strings.Replace(fileName, "http://", "", -1)
		fileName = strings.Replace(fileName, "/", ".", -1)

		// создание файла
		//pathjoin добавить
		filePath := path.Join(*dir, fileName+".html")
		file, err := os.Create(filePath)
		if err != nil {
			writeError(err, *logFlag)
		}
		defer file.Close()
		writeInfo("Создание файла: "+filePath, *logFlag)

		// запись в файл
		file.Write(body)
		writeInfo("Запись в файл: "+filePath, *logFlag)
	}

	if err := scanner.Err(); err != nil {
		writeError(err, *logFlag)
	}
}

// функция отправляет запрос и получает данные
func MakeRequest(
	address string,
	logFlag bool) (body []byte) {

	resp, err := http.Get(address)
	if r := recover(); r != nil {
		writeError(err, logFlag)
	}

	writeInfo("Отправка GET запроса на адрес: "+address, logFlag)

	if resp == nil {
		writeError(errors.New("Resp is nil"), logFlag)
		return
	}

	// if resp.StatusCode != 200 {
	// 	writeError(errors.New("Resp status = "+string(resp.StatusCode)), logFlag)
	// 	return
	// }

	body, err = ioutil.ReadAll(resp.Body)
	if r := recover(); r != nil {
		writeError(err, logFlag)
	}
	defer resp.Body.Close()

	writeInfo("Получение данных с адреса: "+address, logFlag)

	return
}

func writeInfo(infoMessage string, logFlag bool) {
	if logFlag {
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

func writeError(errorMessage error, logFlag bool) {
	logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	errorLogFile := log.New(logFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if logFlag {
		errorLogFile.Println(errorMessage)
	}
	log.Println(errors.New("ERROR\t"), errorMessage)
}

//нужно ли закрывать resp
//сделать функции info, error будут выводить ошибку и инфо в консоль и файл, либо только в консоль , в зависимости от установленного флага
//при запуске программы
//defer scanner как работают
//обработать ошибку, когда один из юрл может быть не корректным, а остальные корректны, программа не должна вылетать
// в функции makeRequest также сделать чтобы не вылетала, а просто возвращала ошибку

// infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
// 	infoLogFile := log.New(logFile, "INFO\t", log.Ldate|log.Ltime)
// 	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
// 	errorLogFile := log.New(logFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
