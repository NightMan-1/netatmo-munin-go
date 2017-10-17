package main

import (
	"fmt"
	"os"
	"github.com/go-ini/ini"
	netatmo "github.com/exzz/netatmo-api-go"
	"strconv"
)

type configGlobalStruct struct {
	ClientID, ClientSecret, Username, Password (string)
}
var configGlobal (configGlobalStruct)

type modulesDataStruct struct {
	name (string)
	data (float64)
}
var modulesData = make(map[string][]modulesDataStruct)



func main() {
	cfg_file := ""
	if _, err := os.Stat("netatmo.cfg"); err == nil {
		cfg_file = "netatmo.cfg"
	}
	if _, err := os.Stat("/etc/munin/netatmo.cfg"); err == nil {
		cfg_file = "/etc/munin/netatmo.cfg"
	}
	if cfg_file == "" {
		fmt.Println("Configuration file not found")
		os.Exit(1)
	}
	cfg, err := ini.Load(cfg_file)
	if err != nil {
		fmt.Println("Can not read", cfg_file)
		os.Exit(1)
	}

	CfgSection, err := cfg.GetSection("")

	configGlobal.ClientID = CfgSection.Key("ClientID").String()
	configGlobal.ClientSecret = CfgSection.Key("ClientSecret").String()
	configGlobal.Username = CfgSection.Key("Username").String()
	configGlobal.Password = CfgSection.Key("Password").String()

	if (configGlobal.ClientID == "" || configGlobal.ClientSecret == "" || configGlobal.Username == "" || configGlobal.Password == ""){
		fmt.Println("Wrong configuration")
		os.Exit(1)
	}

	//help
	if (len(os.Args) == 2 && os.Args[1] == "help"){
		fmt.Println("Netatmo Munin PlugIn v2.0")
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
		for _, module := range station.Modules() {
			_, data := module.Data()
			for dataType, value := range data {
				if (dataType == "LastMesure"){continue}
				t, _ := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
				modulesData[dataType] = append(modulesData[dataType], modulesDataStruct{module.ModuleName, t})
			}

		}
	}

	if (len(os.Args) == 2 && os.Args[1] == "config"){
		// plugin info
		for m_type, m_data := range modulesData { // type name + module data
			switch m_type {
			case "CO2":
				fmt.Println("multigraph netatmo_co2")
				fmt.Println("graph_title CO2 level")
				fmt.Println("graph_vlabel ppm")
				fmt.Println("graph_category netatmo")
				for key, value := range m_data {
					fmt.Printf("co2_%v.label %v\n", key, value.name)
					fmt.Printf("co2_%v.warning  1000\n", key)
					fmt.Printf("co2_%v.critical 1500\n", key)
				}
				fmt.Println("")
			case "Noise":
				fmt.Println("multigraph netatmo_noise")
				fmt.Println("graph_title Noise level")
				fmt.Println("graph_vlabel dB")
				fmt.Println("graph_category netatmo")
				for key, value := range m_data {
					fmt.Printf("noise_%v.label %v\n", key, value.name)
				}
				fmt.Println("")
			case "Pressure":
				fmt.Println("multigraph netatmo_pressure")
				fmt.Println("graph_title Pressure level")
				fmt.Println("graph_vlabel mmHg")
				fmt.Println("graph_category netatmo")
				for key, value := range m_data {
					fmt.Printf("pressure_%v.label %v\n", key, value.name)
				}
				fmt.Println("")
			case "Humidity":
				fmt.Println("multigraph netatmo_humidity")
				fmt.Println("graph_title Humidity level")
				fmt.Println("graph_vlabel %")
				fmt.Println("graph_category netatmo")
				for key, value := range m_data {
					fmt.Printf("hum_%v.label %v\n", key, value.name)
				}
				fmt.Println("")
			case "Temperature":
				fmt.Println("multigraph netatmo_temp")
				fmt.Println("graph_title Temperature level")
				fmt.Println("graph_vlabel °C")
				fmt.Println("graph_category netatmo")
				for key, value := range m_data {
					fmt.Printf("temp_%v.label %v\n", key, value.name)
				}
				fmt.Println("")
			case "WindStrength":
				fmt.Println("multigraph netatmo_wind")
				fmt.Println("graph_title Wind speed")
				fmt.Println("graph_vlabel km/h")
				fmt.Println("graph_category netatmo")
				for key, value := range m_data {
					fmt.Printf("wind_speed_%v.label %v (speed)\n", key, value.name)
					fmt.Printf("wind_gust_%v.label %v (gust)\n", key, value.name)
				}
				fmt.Println("")
			case "Rain1Day":
				fmt.Println("multigraph netatmo_rain")
				fmt.Println("graph_title Rain info")
				fmt.Println("graph_vlabel mm")
				fmt.Println("graph_category netatmo")
				for key, value := range m_data {
					fmt.Printf("rain_daily_%v.label %v (daily)\n", key, value.name)
					fmt.Printf("rain_hourly_%v.label %v (hourly)\n", key, value.name)
				}
				fmt.Println("")
			}
		}
	}else{
		//display data
		for m_type, m_data := range modulesData { // type name + module data
			switch m_type {
			case "CO2":
				fmt.Println("multigraph netatmo_co2")
				for key, value := range m_data {
					fmt.Printf("co2_%v.value %v\n", key, value.data)
				}
				fmt.Println("")
			case "Noise":
				fmt.Println("multigraph netatmo_noise")
				for key, value := range m_data {
					fmt.Printf("noise_%v.value %v\n", key, value.data)
				}
				fmt.Println("")
			case "Pressure":
				fmt.Println("multigraph netatmo_pressure")
				for key, value := range m_data {
					fmt.Printf("pressure_%v.value %00.2f\n", key, value.data * 0.75006375541921)
				}
				fmt.Println("")
			case "Humidity":
				fmt.Println("multigraph netatmo_humidity")
				for key, value := range m_data {
					fmt.Printf("hum_%v.value %v\n", key, value.data)
				}
				fmt.Println("")
			case "Temperature":
				fmt.Println("multigraph netatmo_temp")
				for key, value := range m_data {
					fmt.Printf("temp_%v.value %v\n", key, value.data)
				}
				fmt.Println("")
			case "WindStrength":
				fmt.Println("multigraph netatmo_wind")
				for key, value := range m_data {
					fmt.Printf("wind_speed_%v.value %v\n", key, value.data)
					fmt.Printf("wind_gust_%v.value %v\n", key, value.data)
				}
				fmt.Println("")
			case "Rain1Day":
				fmt.Println("multigraph netatmo_rain")
				for key, value := range m_data {
					fmt.Printf("rain_daily_%v.value %v\n", key, value.data)
					fmt.Printf("rain_hourly_%v.value %v\n", key, value.data)
				}
				fmt.Println("")
			}
		}

	}
}
