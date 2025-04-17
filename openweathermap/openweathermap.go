package openweathermap

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"slices"
	"strings"
	"time"
	"weather/common"
)

func New(cfgPath string) *Openweathermap {
	return &Openweathermap{
		cfg: NewConfig(cfgPath),
	}
}

type Openweathermap struct {
	cfg *Config
}

func (o Openweathermap) getCurrentWeather() WeatherWrapper {
	URL := o.cfg.Base + "/weather?q=" + o.cfg.City + "&units=metric&appid=" + o.cfg.Appid

	body, err := common.GetContent(URL)

	var ww WeatherWrapper

	if err != nil {
		log.Printf("openweathermap: %v", err.Error())
		return ww
	}

	jerr := json.Unmarshal(body, &ww)

	if jerr != nil {
		log.Printf("openweathermap: %v", jerr.Error())
		return ww
	}

	return ww
}

func (o Openweathermap) getWeatherForecast() ForecastList {
	URL := o.cfg.Base + "/forecast?q=" + o.cfg.City + "&cnt=5&&units=metric&appid=" + o.cfg.Appid

	body, err := common.GetContent(URL)

	if err != nil {
		log.Printf("openweathermap: %v", err.Error())
		return nil
	}

	var jsonObjMap map[string]json.RawMessage

	jerr := json.Unmarshal(body, &jsonObjMap)

	if jerr != nil {
		log.Printf("openweathermap: %v", jerr.Error())
		return nil
	}

	var fl ForecastList

	json.Unmarshal(jsonObjMap["list"], &fl)

	return fl
}

func (o Openweathermap) PrepareData() (common.Apidata, error) {

	currentWW := o.getCurrentWeather()
	forcastList := o.getWeatherForecast()

	var apidata = common.Apidata{
		Main:      "N/A",
		Phenomena: "N/A",
		Timestamp: 0,
	}

	if len(currentWW.Weather) == 0 || len(forcastList) == 0 {
		return apidata, fmt.Errorf("openweathermap: can't read data from API (zero len)")
	}

	//PHENOMENA -->
	ShowPhenomenaWW := currentWW
	closestForecastWW := forcastList.getClosestWeather()

	currentWeatherClass := int(currentWW.Weather[0].ID / 100)
	badWeatherClassList := []int{2, 5, 6} //2xx - Thunderstorm, 5xx - Rain, 6xx - Snow...

	if int(closestForecastWW.Weather[0].ID/100) != currentWeatherClass {
		if !slices.Contains(badWeatherClassList, currentWeatherClass) {
			ShowPhenomenaWW = closestForecastWW
			ShowPhenomenaWW.Weather[0].Main = "++" + ShowPhenomenaWW.Weather[0].Main
		}
	}

	apidata.Phenomena = ShowPhenomenaWW.Weather[0].Main
	weatherID := ShowPhenomenaWW.Weather[0].ID

	if weatherID >= 200 && weatherID < 300 {
		apidata.Phenomena = "THND"
	}

	apidata.Phenomena += o.addIntensityMarker(weatherID)
	apidata.Phenomena = fmt.Sprintf("%-10s", apidata.Phenomena)
	//<-- PHENOMENA

	temperature := o.formatTemperature(int(math.Round(currentWW.Main.Temp)))
	feeilsLike := o.formatTemperature(int(math.Round(currentWW.Main.FeelsLike)))

	minForecastTemp, maxForecastTemp := forcastList.getMinMaxTemperature()

	forecastMaxMin := fmt.Sprintf(
		"%s/%s",
		o.formatTemperature(int(math.Round(maxForecastTemp))),
		o.formatTemperature(int(math.Round(minForecastTemp))),
	)

	spacesRequired := common.LCDRowLength - (len(temperature) + len(feeilsLike) + len(forecastMaxMin))

	var leftSpaces int = spacesRequired / 2
	var rightSpaces int = leftSpaces

	if spacesRequired%2 != 0 {
		rightSpaces++
	}

	apidata.Main = fmt.Sprintf(
		"%s%s%s%s%s",
		temperature,
		strings.Repeat(" ", leftSpaces),
		feeilsLike,
		strings.Repeat(" ", rightSpaces),
		forecastMaxMin,
	)

	apidata.Timestamp = time.Now().Unix()

	return apidata, nil
}

func (o Openweathermap) addIntensityMarker(remoteWeatherID uint16) string {

	var result string

	// https://openweathermap.org/weather-conditions#Weather-Condition-Codes-2
	switch remoteWeatherID {
	case 200: // 2xx - Thunderstorm
		result = "#"
	case 201:
		result = "##"
	case 202:
		result = "###"
	case 210:
		result = "!#"
	case 211:
		result = "!##"
	case 212:
		result = "!###"
	case 500: // 5xx - Rain
		result = "#"
	case 501:
		result = "##"
	case 502:
		result = "###"
	case 503:
		result = "####"
	case 600: // 60x - Snow
		result = "*"
	case 601:
		result = "**"
	case 602:
		result = "***"
	case 611: // 61x - Sleet
		result = "*#"
	case 612:
		result = "**#"
	case 613:
		result = "**##"
	case 801: // 80x - Clouds
		result = "#"
	case 802:
		result = "##"
	case 803:
		result = "###"
	case 804:
		result = "####"
	}

	return result
}

func (o Openweathermap) formatTemperature(num int) string {
	if num > 0 {
		return fmt.Sprintf("+%d", num)
	}

	if num < 0 {
		return fmt.Sprintf("%d", num)
	}

	return "0"
}
