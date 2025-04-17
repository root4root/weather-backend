package openweathermap

type Forecast struct {
	Date    uint64    `json:"dt"`
	Main    Main      `json:"main"`
	Weather []Weather `json:"weather"`
	Wind    Wind      `json:"wind"`
	DateTxt string    `json:"dt_txt"`
}

type WeatherWrapper struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
	Wind    Wind      `json:"wind"`
	Name    string    `json:"name"` //City name, ex. "London"
}

type Weather struct {
	ID          uint16 `json:"id"` //2xx, 5xx, 80x, etc..
	Main        string `json:"main"`
	Description string `json:"description"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Pressure  uint16  `json:"grnd_level"`
	Humidity  uint8   `json:"humidity"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   uint16  `json:"deg"`
	Gust  float64 `json:"gust"`
}
