package openweathermap

type ForecastList []Forecast

func (fcl ForecastList) getMinMaxTemperature() (float64, float64) {

	if fcl.len() == 0 {
		return 0.0, 0.0
	}

	min := fcl[0].Main.Temp
	max := min

	for _, fc := range fcl {
		if fc.Main.Temp > max {
			max = fc.Main.Temp
		}

		if fc.Main.Temp < min {
			min = fc.Main.Temp
		}
	}

	return min, max
}

func (fcl ForecastList) len() int {
	return len(fcl)
}

func (fcl ForecastList) getClosestWeather() WeatherWrapper {
	ww := WeatherWrapper{
		Weather: fcl[0].Weather,
		Main:    fcl[0].Main,
		Wind:    fcl[0].Wind,
		Name:    "undefined",
	}

	return ww
}
