package main
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)
const apiKey = "7f6c3041659eea0086fad104037d3285"
func main() {
	fmt.Println("Data:", time.Now())
	fmt.Println("Autor: Kiryl Shkabara")
	fmt.Println("Port: 3000")
	http.HandleFunc("/", home)
	http.HandleFunc("/weather", weather)
	http.HandleFunc("/health", health)
	http.ListenAndServe(":3000", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
	<!DOCTYPE html>
	<html>
		<title>Zad1</title>
		<h2>Pogoda</h2>
		<select id="country"></select>
		<select id="city"></select>
		<button onclick="getWeather()">Pobierz pogode</button>
		<p id="r"></p>

		<script>
			const cities = {
				Poland:["Warsaw","Krakow"],
				Germany:["Berlin","Munich"],
			}

			const c = document.getElementById("country")
			const m = document.getElementById("city")
			const r = document.getElementById("r")

			Object.keys(cities).forEach(x => c.innerHTML += "<option>" + x + "</option>")

			function loadCities() {
				m.innerHTML = ""
				cities[c.value].forEach(x => m.innerHTML += "<option>" + x + "</option>")
			}

			c.onchange = loadCities
			loadCities()

			async function getWeather() {
				const d = await fetch("/weather?city=" + m.value).then(r => r.json())
				r.innerHTML = "<b>" + m.value + "</b><br><br>Temperatura: " + d.temp + " °C<br>Opis: " + d.description
			}
		</script>
	</html>`
	w.Write([]byte(html))
}

func weather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "API error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	main := data["main"].(map[string]interface{})
	weather := data["weather"].([]interface{})[0].(map[string]interface{})

	result := map[string]interface{}{
		"temp":        main["temp"],
		"description": weather["description"],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}