package main

import (
	"GoRegexp/pkg/api"
	"GoRegexp/pkg/calcregexp"

	"flag"
	"fmt"
	"net/http"
)

// Сервер GoRegexp.
type server struct {
	api *api.API
}

func main() {

	// Обрабатываем флаги при запуске программы
	// go run server.go -inpfile "./ui/data/input.txt" -outfile "./ui/data/output.txt"
	var inpfile string
	var outfile string

	flag.StringVar(&inpfile, "inpfile", "./ui/data/input.txt", "File for reading data")
	flag.StringVar(&outfile, "outfile", "./ui/data/output.txt", "File for writting data")
	flag.Parse()

	fmt.Println("flags: inpfile->", inpfile, "; outfile->", outfile)

	// Создаём объект сервера.
	var srv server

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(calcregexp.NewCalc(), inpfile, outfile)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	fmt.Println("Запуск веб-сервера на http://127.0.0.1:8080 ...")
	http.ListenAndServe(":8080", srv.api.Router())
}
