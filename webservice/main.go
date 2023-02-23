package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

func main() {

	// Initialize variables
	page := "<h1>ForgeRock Assignment</h1>"
	symbol := os.Getenv("SYMBOL")
	api_key := os.Getenv("API_KEY")
	closing_sum := 0.0
	i := 0
	var dates []string
	num_days, err := strconv.Atoi(os.Getenv("NUM_DAYS"))
	if err != nil {
		log.Fatal(err)
	}

	// Create http client
	endpoint := fmt.Sprintf("https://www.alphavantage.co/query?apikey=%s&function=TIME_SERIES_DAILY_ADJUSTED&symbol=%s", api_key, symbol)

	resp, err := http.Get(endpoint)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// Get data from API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Convert JSON data to maps
	var result map[string]interface{}

	json.Unmarshal([]byte(body), &result)
	result = result["Time Series (Daily)"].(map[string]interface{})

	// Sort dates in API data
	for d := range result {
		dates = append(dates, d)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	// Iterate through dates
	for _, date_key := range dates {
		if i > num_days-1 {
			break
		}

		nest3 := result[date_key].(map[string]interface{})

		// Sort keys in API data
		var keys []string
		for k := range nest3 {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Build HTML
		page += fmt.Sprintf("<h2>For the date of %s</h2>", date_key)

		// Iterate through values in API data
		for _, sub_val := range keys {

			// Remove first 3 chars from string
			sub_str := sub_val[3:]

			// Convert string to float
			data := fmt.Sprintf("%s", nest3[sub_val])
			data_val, err := strconv.ParseFloat(data, 64)
			if err != nil {
				log.Fatal(err)
			}

			// Sum up data values
			if sub_str == "close" {
				closing_sum = closing_sum + data_val

			}

			page += fmt.Sprintf("<p>The %s value is: %s</p>", sub_str, nest3[sub_val])

		}

		i++

	}

	// Compute average
	page += fmt.Sprintf("<h3>The average closing value for the selected days was %.2f</h3>", closing_sum/float64(num_days))

	// Generate index.html file
	write_file(page)

	// Start web server
	render_page()
}

func write_file(resp string) {

	_ = os.Mkdir("/tmp/lab", 0755)

	f, err := os.Create("/tmp/lab/index.html")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.WriteString(resp)

	if err != nil {
		log.Fatal(err)
	}

}

func render_page() {
	p := http.FileServer(http.Dir("/tmp/lab"))
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", p))
}
