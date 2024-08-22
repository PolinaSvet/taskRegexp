package logger

import (
	"encoding/json"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const fileNameLog string = "ui/data/log.json"

type MessStructLog struct {
	Datetime time.Time `json:"datetime"`
	Database string    `json:"database"`
	Mess     string    `json:"mess"`
}

func SetupLogger() {
	// Настраиваем zerolog
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})
}

func SetLog(datetime time.Time, base string, mess string) {
	// Создаем новый экземпляр структуры
	r := MessStructLog{
		Datetime: datetime,
		Database: base,
		Mess:     mess,
	}

	// Записываем данные в JSON-формате
	logData, err := json.Marshal(r)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal log data")
		return
	}

	// Открываем файл для записи
	f, err := os.OpenFile(fileNameLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open log file")
		return
	}
	defer f.Close()

	// Записываем данные в файл
	_, err = f.Write(append(logData, '\n'))
	if err != nil {
		log.Error().Err(err).Msg("Failed to write to log file")
		return
	}

	// Записываем данные в консольный лог
	log.Info().
		Time("datetime", r.Datetime).
		Str("base", r.Database).
		Str("mess", r.Mess).
		Msg("mess from database")

}

func GetLog() {
	// Открываем файл для чтения
	f, err := os.Open(fileNameLog)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open log file")
		return
	}
	defer f.Close()

	// Читаем данные из файла
	var logs []MessStructLog
	decoder := json.NewDecoder(f)
	for {
		var r MessStructLog
		err := decoder.Decode(&r)
		if err != nil {
			break
		}
		logs = append(logs, r)
	}

	// Выводим данные в консоль
	for _, r := range logs {
		log.Info().
			Time("datetime", r.Datetime).
			Str("base", r.Database).
			Str("mess", r.Mess).
			Msg("mess from database")
	}
}
