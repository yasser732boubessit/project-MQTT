// Function to update the weather card with data
function updateWeather(data) {
  const info = data.list[0];
  document.getElementById("city").textContent = `📍 ${data.city.name}`;
  document.getElementById("temp").textContent = `${info.main.temp.toFixed(1)}°C`;
  document.getElementById("desc").textContent = info.weather[0].description;
  document.getElementById("humidity").textContent = `${info.main.humidity}%`;
  document.getElementById("wind").textContent = `${(Math.random() * 10 + 2).toFixed(1)} km/h`;
}

// Fetch latest MQTT data from /data endpoint
function fetchLatestData() {
  fetch("/data")
    .then(res => res.json())
    .then(data => {
      if (data.message) {
        console.log("No MQTT data yet.");
      } else {
        updateWeather(data);
      }
    })
    .catch(err => console.error("Error fetching MQTT data:", err));
}

// Fetch new weather data from OpenWeather API (via Go backend)
function fetchWeatherByCity(city) {
  fetch(`/api/weather?city=${city}`)
    .then(res => res.json())
    .then(data => updateWeather(data))
    .catch(() => alert("An error occurred while fetching data 😢"));
}

// When the page loads → get the latest MQTT data automatically
window.addEventListener("load", () => {
  const city =  "jijel";
  fetchWeatherByCity(city);
});

// When the user clicks the search button → request new data
document.getElementById("search-btn").addEventListener("click", () => {
  const city = document.getElementById("city-input").value || "jijel";
  fetchWeatherByCity(city);
});
