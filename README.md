# 🌦️ MQTT Weather Forecast App

## 📘 Description
This Go application connects to a public MQTT broker (HiveMQ) to exchange live weather forecast data.  
It also exposes a built-in web server that displays this data in a beautiful, interactive dashboard.  
The app combines MQTT messaging, REST APIs, and a frontend interface (HTML/CSS/JS).





---

## ⚙️ How It Works

### 🧩 Architecture Overview

**Backend (Go)**

- Connects to the HiveMQ MQTT broker (`tcp://broker.hivemq.com:1883`).
- Subscribes to a topic: `rmbtech/interview/rug/yasser/weather_forecast`.
- Fetches real-time weather data from OpenWeatherMap API.
- Publishes the weather forecast to the MQTT topic.
- Exposes two REST endpoints:
  - `/api/weather?city=...` → fetches and publishes new data.
  - `/data` → returns the latest MQTT message received.
- Hosts static frontend files under `/static`.

**Frontend (HTML + CSS + JS)**

- Allows users to enter any city name.
- Fetches live weather data through the backend.
- Displays temperature, humidity, weather description, and random wind speed.
- Automatically updates when the user searches for a new city.

**MQTT Communication Flow**

- The backend publishes weather JSON data on the HiveMQ broker.
- The backend also subscribes to the same topic to receive data.
- The received message is stored in memory and accessible through `/data`.

---

## 🛠️ How It’s Made

### 🔧 Technologies Used

| Component | Technology |
|-----------|------------|
| Backend   | Go (net/http, Eclipse Paho MQTT client) |
| Frontend  | HTML, CSS, JavaScript |
| Broker    | HiveMQ Public MQTT Broker |
| API       | OpenWeatherMap API (5-day forecast) |

### 📄 Folder Structure
```

/project-MQTT
│
├── main.go
├── go.mod
├── /static
│   ├── index.html
│   ├── style.css
│   └── script.js
└── README.md

````

### 🧱 Backend Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/weather?city=<name>` | GET | Fetches weather data for a city and publishes it to MQTT |
| `/data` | GET | Returns the latest data received from the MQTT broker |
| `/` | GET | Serves the web UI |

---

## 🚀 How to Run

### 🔹 Option 1 — Run from Source
```bash
# Clone repository
git clone https://github.com/yourusername/mqtt-weather-app.git
cd mqtt-weather-app

# Run the Go app
go run main.go
````

Then open [http://localhost:8080](http://localhost:8080) in your browser.

### 🔹 Option 2 — Run Compiled Binary

```bash
./mqtt-weather-app
```

---

## 💡 How to Use

1. Open the web interface.
2. Enter the name of a city (e.g., Paris, Tokyo, New York).
3. Click the search button.

The dashboard will display:

* 🌡️ Current temperature
* 💧 Humidity
* ☁️ Weather description
* 💨 Wind speed (simulated for visual appeal)

### 🧠 Example of JSON Data (MQTT Payload)

```json
{
  "city": { "name": "Jijel" },
  "list": [
    {
      "dt_txt": "2025-10-23 15:00:00",
      "main": { "temp": 25.4, "humidity": 64 },
      "weather": [{ "description": "clear sky" }]
    }
  ]
}
```

---

## 🧰 Code Quality and Modularity

* The backend separates MQTT logic, HTTP handling, and API fetching.
* Thread safety ensured using `sync.Mutex` when reading/writing shared data.
* Frontend separated into modular files (`index.html`, `style.css`, `script.js`).
* The Go app automatically serves static files without external servers.

---

## 🤖 Use of LLM

This documentation and parts of the code structure were improved and explained using ChatGPT (GPT-5) for clarity, formatting, and technical precision.

---

## 🧠 Possible Improvements

* Add WebSocket or SSE for true live streaming updates.
* Store MQTT data in a local database.
* Allow multi-city comparison or charts.

---

## 🏁 Author

**Yasser Boubessit**

