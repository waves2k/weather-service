package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
	"github.com/waves2k/weather-service/internal/client/http/geocoding"
	"github.com/waves2k/weather-service/internal/client/http/open_meteo"
)

// http://geocoding-api.open-meteo.com/v1/search?name=Moscow&count=1&language=ru&format=json
//https://api.open-meteo.com/v1/forecast?latitude=55.75&longitude=37.62&current_weather=true
// http://api.open-meteo.com/v1/forecast?latitude=55&longitude=37&current=temperature2_m

const httpPort = ":3000"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	geocodingClient := geocoding.NewClient(httpClient)
	openMeteoClient := open_meteo.NewClient(httpClient)
	r.Get("/{city}", func(w http.ResponseWriter, r *http.Request) {
		city := chi.URLParam(r, "city")

		coordsRes, err := geocodingClient.GetCoords(city)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		fmt.Println(coordsRes)

		tempRes, err := openMeteoClient.GetTemperature(coordsRes.Latitude, coordsRes.Longitude)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		fmt.Println(tempRes)

		raw, err := json.Marshal(tempRes)
		if err != nil {
			log.Printf(err.Error())
		}

		w.Write([]byte(raw))
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	jobs, err := initJobs(s)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		fmt.Println("starting server on port " + httpPort)
		if err := http.ListenAndServe(httpPort, r); err != nil {
			panic(err)
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		fmt.Printf("starting job: %v\n", jobs[0].ID())
		s.Start()
	}(&wg)

	wg.Wait()
}

func initJobs(scheduler gocron.Scheduler) ([]gocron.Job, error) {

	j, err := scheduler.NewJob(
		gocron.DurationJob(
			30*time.Minute,
		),
		gocron.NewTask(
			func() {
				fmt.Println("Hello")
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return []gocron.Job{j}, nil
}

// func runCron() {

// 	// each job has a unique id
// 	fmt.Println(j.ID())

// 	// start the scheduler
// 	s.Start()
// }
