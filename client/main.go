package main

import (
	"github.com/gordonklaus/portaudio"
	"log"
	"time"
)

func main() {
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, 1024, func(in, out []float32) {
		for i := range in {
			out[i] = in[i] // Копируем входные данные в выходной буфер
		}
	})
	if err != nil {
		log.Fatal("Ошибка открытия потока:", err)
	}
	defer stream.Close()

	// Запускаем поток
	err = stream.Start()
	if err != nil {
		log.Fatal("Ошибка запуска потока:", err)
	}
	defer stream.Stop()

	log.Println("Запись и воспроизведение начаты. Говори в микрофон! Остановка через 10 секунд...")

	// Даём время для записи и воспроизведения (10 секунд)
	time.Sleep(10 * time.Second)

	log.Println("Остановка записи и воспроизведения")
}
