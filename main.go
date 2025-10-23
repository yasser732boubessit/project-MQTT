package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ForecastData struct {
	City struct {
		Name string `json:"name"`
	} `json:"city"`
	List []struct {
		DtTxt string `json:"dt_txt"`
		Main  struct {
			Temp     float64 `json:"temp"`
			Humidity int     `json:"humidity"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
	} `json:"list"`
}

// Variable to store the latest data received from the MQTT broker
var (
	latestData string
	dataMutex  sync.Mutex
)

// ---- Function to fetch weather forecast data from OpenWeatherMap API ----
func getForecastData(apiKey, city string) (*ForecastData, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric&lang=en", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println("⚠️ API returned status:", resp.StatusCode)
		fmt.Println("Response:", string(body))
	}

	var data ForecastData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func main() {
	broker := "tcp://broker.hivemq.com:1883"
	topic := "rmbtech/interview/rug/yasser/weather_forecast"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("golang-weather-client")

	// Callback when the connection to the MQTT broker is established
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("✅ Connected to MQTT broker")

		// Subscribe to the topic and handle incoming messages
		if token := c.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			fmt.Println("📩 Received new MQTT message")
			dataMutex.Lock()
			latestData = string(msg.Payload()) // Store the latest received message
			dataMutex.Unlock()
		}); token.Wait() && token.Error() != nil {
			fmt.Println("❌ Subscription error:", token.Error())
		}
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	apiKey := "ce0c141507f08d8e2787a3b636ac2827"

	// ---- REST API endpoint: fetch forecast for a city and publish to MQTT ----
	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("city")
		if city == "" {
			city = "jijel"
		}

		data, err := getForecastData(apiKey, city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(data.List) > 0 {
			current := data.List[0]
			fmt.Printf("\n🏙️ City: %s\n🌡️ %.1f°C\n💧 %d%%\n☁️ %s\n🕒 %s\n",
				data.City.Name, current.Main.Temp, current.Main.Humidity, current.Weather[0].Description, current.DtTxt)
		}

		payload, _ := json.Marshal(data)
		client.Publish(topic, 0, false, payload)
		fmt.Printf("✅ Published forecast for %s\n", data.City.Name)

		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	})

	// ---- REST endpoint to get the latest MQTT data ----
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dataMutex.Lock()
		defer dataMutex.Unlock()
		if latestData == "" {
			w.Write([]byte(`{"message":"No data received yet"}`))
			return
		}
		w.Write([]byte(latestData))
	})

	// ---- Serve static frontend files (HTML, CSS, JS) ----
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("🌍 Web server running at: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
