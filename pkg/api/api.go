package api

import (
	"GoRegexp/pkg/calcregexp"
	"GoRegexp/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

// Программный интерфейс сервера GoNews
type API struct {
	reg        calcregexp.Calc
	router     *mux.Router
	inputFile  string
	outputFile string
}

// Конструктор объекта API
func New(reg calcregexp.Calc, inputFile, outputFile string) *API {
	api := API{
		reg:        reg,
		inputFile:  inputFile,
		outputFile: outputFile,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {

	api.router.HandleFunc("/", api.templateHandler).Methods(http.MethodGet, http.MethodOptions)

	api.router.HandleFunc("/expLine", api.addExpLineHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/expFile", api.addExpFileHandler).Methods(http.MethodPost, http.MethodOptions)

	// Регистрация обработчика для статических файлов (шаблонов)
	api.router.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("ui"))))
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.router
}

// Базовый маршрут.
func (api *API) templateHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("ui/html/base.html", "ui/html/routes.html"))

	// Отправляем HTML страницу с данными
	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Расчет одного выражения.
func (api *API) addExpLineHandler(w http.ResponseWriter, r *http.Request) {
	//принимаем данные
	var jsonDataMap map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonDataMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//преобразуем в нужный формат
	var line []string
	for _, v := range jsonDataMap {
		line = append(line, fmt.Sprintf("%v", v))
	}
	//вычисляем
	data := api.reg.Сalculate(line, true)

	//отдаем данные
	bytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

// Расчет выражений из файла.
func (api *API) addExpFileHandler(w http.ResponseWriter, r *http.Request) {
	//принимаем данные
	var jsonDataMap map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonDataMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//преобразуем в нужный формат
	inputFile := api.inputFile
	name, ok := jsonDataMap["inputFile"]
	if ok && name != "" {
		inputFile = fmt.Sprintf("%v", name)
	}

	outputFile := api.outputFile
	name, ok = jsonDataMap["outputFile"]
	if ok && name != "" {
		outputFile = fmt.Sprintf("%v", name)
	}

	//читаем все данные из файла
	lines, err := api.reg.ReadLinesFromFile(inputFile)
	data := make([]string, 0)
	if err != nil {
		logger.SetLog(time.Now(), "", fmt.Sprintf("ошибка чтения входоного файла: %v", err))
		data = append(data, fmt.Sprintf("ошибка чтения входоного файла: %v", err))
	} else {
		//вычисляем
		data = api.reg.Сalculate(lines, true)
		//записываем все данные в файла
		err = api.reg.WriteLinesToFile(outputFile, data)
		if err != nil {
			logger.SetLog(time.Now(), "", fmt.Sprintf("ошибка записи выходного файла: %v", err))
			data = append(data, fmt.Sprintf("ошибка чтения входоного файла: %v", err))
		}
	}

	//отдаем данные
	bytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}
