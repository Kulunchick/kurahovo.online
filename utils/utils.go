package utils

import (
	"database/sql"
	"fmt"
	"sync"
)

type LatestData struct {
	AQI  Data `json:"aqi"`
	PM25 Data `json:"pm25"`
	PM10 Data `json:"pm10"`
	Temp Data `json:"temp"`
	Rel  Data `json:"rel"`
	Luft Data `json:"luft"`
}

type Data struct {
	Value string `json:"value"`
	Date  string `json:"date"`
}

func GetJatestData(db *sql.DB) (LatestData, []error) {
	var data LatestData
	var errors []error

	wg := sync.WaitGroup{}
	wg.Add(6)

	cmd := "SELECT ROUND(value, 1), date FROM %s WHERE date >= DATE_SUB(NOW(), INTERVAL %d MINUTE) ORDER BY date DESC LIMIT 0, 1"

	go func() {
		defer wg.Done()
		err := db.QueryRow(fmt.Sprintf(cmd, "AvgAQI", 30)).Scan(&data.AQI.Value, &data.AQI.Date)
		switch {
		case err == sql.ErrNoRows:
			data.AQI.Value, data.AQI.Date = "N/A", "N/A"
		case err != nil:
			errors = append(errors, err)
		}
	}()
	go func() {
		defer wg.Done()
		err := db.QueryRow(fmt.Sprintf(cmd, "AvgPM25", 15)).Scan(&data.PM25.Value, &data.PM25.Date)
		switch {
		case err == sql.ErrNoRows:
			data.PM25.Value, data.PM25.Date = "N/A", "N/A"
		case err != nil:
			errors = append(errors, err)
		}
	}()
	go func() {
		defer wg.Done()
		err := db.QueryRow(fmt.Sprintf(cmd, "AvgPM10", 15)).Scan(&data.PM10.Value, &data.PM10.Date)
		switch {
		case err == sql.ErrNoRows:
			data.PM10.Value, data.PM10.Date = "N/A", "N/A"
		case err != nil:
			errors = append(errors, err)
		}
	}()
	go func() {
		defer wg.Done()
		err := db.QueryRow(fmt.Sprintf(cmd, "AvgTem", 15)).Scan(&data.Temp.Value, &data.Temp.Date)
		switch {
		case err == sql.ErrNoRows:
			data.Temp.Value, data.Temp.Date = "N/A", "N/A"
		case err != nil:
			errors = append(errors, err)
		}
	}()
	go func() {
		defer wg.Done()
		err := db.QueryRow(fmt.Sprintf(cmd, "AvgRel", 15)).Scan(&data.Rel.Value, &data.Rel.Date)
		switch {
		case err == sql.ErrNoRows:
			data.Rel.Value, data.Rel.Date = "N/A", "N/A"
		case err != nil:
			errors = append(errors, err)
		}
	}()
	go func() {
		defer wg.Done()
		err := db.QueryRow(fmt.Sprintf(cmd, "AvgLuft", 15)).Scan(&data.Luft.Value, &data.Luft.Date)
		switch {
		case err == sql.ErrNoRows:
			data.Luft.Value, data.Luft.Date = "N/A", "N/A"
		case err != nil:
			errors = append(errors, err)
		}
	}()
	wg.Wait()

	return data, errors
}
