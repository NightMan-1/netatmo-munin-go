package main

import (
	"fmt"
	"os"
	"strconv"
	netatmo "github.com/exzz/netatmo-api-go"
	toml "github.com/BurntSushi/toml"
)

type configGlobalStruct struct {
	ClientID, ClientSecret, Username, Password (string)
	Sensors (configSensorsStruct)

}

type configSensorsStruct struct {
	CO2, Noise, Pressure, Humidity, Temperature, WindStrength, Rain, WindAngle, BatteryPercent, WifiStatus, RFStatus (bool)
}

var configGlobal (configGlobalStruct)


type modulesDataStruct struct {
	name (string)
	data (float64)
}
var modulesData = make(map[string][]modulesDataStruct)

var DisplayConfig = false


func main() {
	cfg_file := ""
	if _, err := os.Stat("/etc/munin/netatmo.cfg"); err == nil {
		cfg_file = "/etc/munin/netatmo.cfg"
	}
	if _, err := os.Stat("netatmo.cfg"); err == nil {
		cfg_file = "netatmo.cfg"
	}
	if cfg_file == "" {
		fmt.Println("Configuration file not found")
		os.Exit(1)
	}

	configGlobal.Sensors = configSensorsStruct{true, true, true, true, true, true, true, false, true, false, true}
	if _, err := toml.DecodeFile(cfg_file, &configGlobal); err != nil {
		fmt.Printf("Cannot parse config file: %s\n", err)
		os.Exit(1)
	}
	if (configGlobal.ClientID == "" || configGlobal.ClientSecret == "" || configGlobal.Username == "" || configGlobal.Password == "") {
		fmt.Println("Wrong configuration")
		os.Exit(1)
	}
	if (configGlobal.ClientID == "ClientID" || configGlobal.ClientSecret == "ClientSecret" || configGlobal.Username == "Username" || configGlobal.Password == "Password") {
		fmt.Println("Do not forget change example configuration settings :)")
		os.Exit(1)
	}



	//help
	if (len(os.Args) == 2 && os.Args[1] == "help") {
		fmt.Println("Netatmo Munin PlugIn")
		fmt.Println("https://github.com/NightMan-1/netatmo-munin-go")
		fmt.Println("(c)2017 Sergey Gurinovich")
		os.Exit(0)
	}

	//connect NetAtmo
	n, err := netatmo.NewClient(netatmo.Config{
		ClientID:     configGlobal.ClientID,
		ClientSecret: configGlobal.ClientSecret,
		Username:     configGlobal.Username,
		Password:     configGlobal.Password,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dc, err := n.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//parce NetAtmo data
	for _, station := range dc.Stations() {
		//fmt.Printf("Station : %v\n", station.StationName)
		for _, module := range station.Modules() {
			//fmt.Printf("\tModule : %v\n", module.ModuleName)
			_, data := module.Data()
			for dataType, value := range data {
				//fmt.Printf("\t\t%s : %v\n", dataType, value)
				if (dataType == "LastMesure") {
					continue
				}
				t, _ := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
				if (len(dc.Stations()) > 1){
					modulesData[dataType] = append(modulesData[dataType], modulesDataStruct{module.ModuleName + " (" + station.StationName + ")", t})
				}else{
					modulesData[dataType] = append(modulesData[dataType], modulesDataStruct{module.ModuleName, t})
				}
			}
			_, info := module.Info()
			for dataType, value := range info {
				//fmt.Printf("\t\t%s : %v\n", dataType, value)
				t, _ := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
				if (len(dc.Stations()) > 1){
					modulesData[dataType] = append(modulesData[dataType], modulesDataStruct{module.ModuleName + " (" + station.StationName + ")", t})
				}else{
					modulesData[dataType] = append(modulesData[dataType], modulesDataStruct{module.ModuleName, t})
				}
			}

		}
	}

	if (len(os.Args) == 2 && os.Args[1] == "config") { DisplayConfig = true }

	for m_type, m_data := range modulesData { // type name + module data
		switch m_type {
		case "BatteryPercent":
			if (configGlobal.Sensors.BatteryPercent) {
				fmt.Println("multigraph netatmo_battery")
				if (DisplayConfig){
					fmt.Println("graph_title BatteryPercent level")
					//fmt.Println("graph_vlabel level")
					fmt.Println("graph_category netatmo")
					for key, value := range m_data {
						fmt.Printf("battery_%v.label %v\n", key, value.name)
					}
				}else{
					for key, value := range m_data {
						fmt.Printf("battery_%v.value %v\n", key, value.data)
					}
				}
			}
		case "WifiStatus":
			if (configGlobal.Sensors.BatteryPercent) {
				fmt.Println("multigraph netatmo_WifiStatus")
				if (DisplayConfig){
					fmt.Println("graph_title WifiStatus level")
					//fmt.Println("graph_vlabel level")
					fmt.Println("graph_category netatmo")
					for key, value := range m_data {
						fmt.Printf("WifiStatus_%v.label %v\n", key, value.name)
					}
				}else{
					for key, value := range m_data {
						fmt.Printf("WifiStatus_%v.value %v\n", key, value.data)
					}
				}
			}
		case "RFStatus":
			if (configGlobal.Sensors.BatteryPercent) {
				fmt.Println("multigraph netatmo_RFStatus")
				if (DisplayConfig){
					fmt.Println("graph_title RFStatus level")
					//fmt.Println("graph_vlabel level")
					fmt.Println("graph_category netatmo")
					for key, value := range m_data {
						fmt.Printf("RFStatus_%v.label %v\n", key, value.name)
					}
				}else{
					for key, value := range m_data {
						fmt.Printf("RFStatus_%v.value %v\n", key, value.data)
					}
				}
			}
		case "CO2":
			if (configGlobal.Sensors.CO2) {
				fmt.Println("multigraph netatmo_co2")
				if (DisplayConfig){
					fmt.Println("graph_title CO2 level")
					fmt.Println("graph_vlabel ppm")
					fmt.Println("graph_category netatmo")
					for key, value := range m_data {
						fmt.Printf("co2_%v.label %v\n", key, value.name)
						fmt.Printf("co2_%v.warning  1000\n", key)
						fmt.Printf("co2_%v.critical 1500\n", key)
					}
				}else{
					for key, value := range m_data {
						fmt.Printf("co2_%v.value %v\n", key, value.data)
					}
				}
			}
		case "Noise":
			if (configGlobal.Sensors.Noise) {
				fmt.Println("multigraph netatmo_noise")
				if (DisplayConfig){
					fmt.Println("graph_title Noise level")
					fmt.Println("graph_vlabel dB")
					fmt.Println("graph_category netatmo")
					for key, value := range m_data {
						fmt.Printf("noise_%v.label %v\n", key, value.name)
					}
				}else{
					for key, value := range m_data {
						fmt.Printf("noise_%v.value %v\n", key, value.data)
					}
				}
			}
		case "Pressure":
			if (configGlobal.Sensors.Pressure) {
				fmt.Println("multigraph netatmo_pressure")
				if (DisplayConfig){
					fmt.Println("graph_title Pressure level")
					fmt.Println("graph_vlabel mmHg")
					fmt.Println("graph_category netatmo")
					for key, value := range m_data {
						fmt.Printf("pressure_%v.label %v\n", key, value.name)
					}
				}else{
					for key, value := range m_data {
						fmt.Printf("pressure_%v.value %00.2f\n", key, value.data*0.75006375541921)
					}

				}
			}
		case "Humidity":
			if (configGlobal.Sensors.Humidity) {
				fmt.Println("multigraph netatmo_humidity")
				if (DisplayConfig){
					fmt.Println("graph_title Humidity level")
					fmt.Println("graph_vlabel %")
					fmt.Println("graph_category netatmo")
					for key, value := range m_data {
						fmt.Printf("hum_%v.label %v\n", key, value.name)
					}
				}else{
					for key, value := range m_data {
						fmt.Printf("hum_%v.value %v\n", key, value.data)
					}
				}
			}
		case "Temperature":
			if (configGlobal.Sensors.Temperature) {
				fmt.Println("multigraph netatmo_temp")
				if (DisplayConfig){
					fmt.Println("graph_title Temperature level")
					fmt.Println("graph_vlabel °C")
					fmt.Println("graph_category netatmo")
					for key, value := range m_data {
						fmt.Printf("temp_%v.label %v\n", key, value.name)
					}
				}else{
					for key, value := range m_data {
						fmt.Printf("temp_%v.value %v\n", key, value.data)
					}
				}
			}
		case "WindStrength":
			if (configGlobal.Sensors.WindStrength) {
				fmt.Println("multigraph netatmo_wind")
				if (DisplayConfig){
					fmt.Println("graph_title Wind speed")
					fmt.Println("graph_vlabel km/h")
					fmt.Println("graph_category netatmo")
					for key, value := range m_data {
						fmt.Printf("wind_speed_%v.label %v (speed)\n", key, value.name)
						fmt.Printf("wind_gust_%v.label %v (gust)\n", key, value.name)
					}
				}else{
					for key, _ := range m_data {
						fmt.Printf("wind_speed_%v.value %v\n", key, modulesData["WindStrength"][key].data)
						fmt.Printf("wind_gust_%v.value %v\n", key, modulesData["GustStrength"][key].data)
					}
				}
			}
		case "Rain1Day":
			if (configGlobal.Sensors.Rain) {
				fmt.Println("multigraph netatmo_rain")
				if (DisplayConfig){
					fmt.Println("graph_title Rain info")
					fmt.Println("graph_vlabel mm")
					fmt.Println("graph_scale no")
					fmt.Println("graph_category netatmo")
					for key, value := range m_data {
						fmt.Printf("rain_daily_%v.label %v (daily)\n", key, value.name)
						fmt.Printf("rain_hourly_%v.label %v (hourly)\n", key, value.name)
						fmt.Printf("rain_daily_%v.type COUNTER\n", key)
						fmt.Printf("rain_hourly_%v.type COUNTER\n", key)
						fmt.Printf("rain_daily_%v.min 0\n", key)
						fmt.Printf("rain_hourly_%v.min 0\n", key)
					}
				}else{
					for key, _ := range m_data {
						fmt.Printf("rain_daily_%v.value %v\n", key, modulesData["Rain1Day"][key].data)
						fmt.Printf("rain_hourly_%v.value %v\n", key, modulesData["Rain1Hour"][key].data)
					}
				}
			}
		case "WindAngle":
			if (configGlobal.Sensors.WindAngle) {
				fmt.Println("multigraph netatmo_wind_angl")
				if (DisplayConfig){
					fmt.Println("graph_title Wind direction")
					fmt.Println("graph_vlabel degrees")
					fmt.Println("graph_category netatmo")
					for key, value := range m_data {
						fmt.Printf("wind_angle_%v.label %v\n", key, value.name)
						fmt.Printf("gust_angle_%v.label %v\n", key, value.name)
					}
				}else{
					for key, _ := range m_data {
						fmt.Printf("wind_angle_%v.value %v\n", key, modulesData["WindAngle"][key].data)
						fmt.Printf("gust_angle_%v.value %v\n", key, modulesData["GustAngle"][key].data)
					}
				}
			}
		}
		fmt.Println("")
	}

}
