package main

import (
	"fmt"
	"os"
	"github.com/go-ini/ini"
	netatmo "github.com/exzz/netatmo-api-go"
)

type configGlobalStruct struct {
	ClientID, ClientSecret, Username, Password (string)
}
var configGlobal (configGlobalStruct)

func main() {
	//read configuration
	cfg, err := ini.Load("/etc/munin/netatmo.cfg")
	if err != nil {
		fmt.Println("Can not read /etc/munin/netatmo.cfg")
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
	if (len(os.Args) == 1 ||(len(os.Args) == 2 && os.Args[1] == "help")){
		fmt.Println("Netatmo Munin PlugIn 1.0")
		fmt.Println("--------------------")
		fmt.Println("Usage:")
		fmt.Println("\t -co2 - for CO2 info")
		fmt.Println("\t -temp - for Temperature info")
		fmt.Println("\t -hum - for Humidity info")
		fmt.Println("\t -noise - for Noise info")
		fmt.Println("\t -pressure - for Pressure info")
		os.Exit(0)
	}

    // plugin info
	if (len(os.Args) == 3 && os.Args[1] == "-co2" && os.Args[2] == "config"){
		fmt.Println("graph_title Netatmo CO2 level")
		fmt.Println("graph_vlabel level")
		fmt.Println("graph_category netatmo")
		fmt.Println("co2.label level")
		fmt.Println("co2.warning  1000")
		fmt.Println("co2.critical 1500")
		os.Exit(0)
	}
	if (len(os.Args) == 3 && os.Args[1] == "-noise" && os.Args[2] == "config"){
		fmt.Println("graph_title Netatmo Noise level")
		fmt.Println("graph_vlabel level")
		fmt.Println("graph_category netatmo")
		fmt.Println("noise.label level")
		os.Exit(0)
	}
	if (len(os.Args) == 3 && os.Args[1] == "-pressure" && os.Args[2] == "config"){
		fmt.Println("graph_title Netatmo Pressure level")
		fmt.Println("graph_vlabel level")
		fmt.Println("graph_category netatmo")
		fmt.Println("pressure.label level")
		os.Exit(0)
	}
	if (len(os.Args) == 3 && os.Args[1] == "-hum" && os.Args[2] == "config"){
		fmt.Println("graph_title Netatmo Humidity level")
		fmt.Println("graph_vlabel level")
		fmt.Println("graph_category netatmo")
		fmt.Println("hum_indor.label indor")
		fmt.Println("hum_outdor.label outdor")
		os.Exit(0)
	}
	if (len(os.Args) == 3 && os.Args[1] == "-temp" && os.Args[2] == "config"){
		fmt.Println("graph_title Netatmo Temperature level")
		fmt.Println("graph_vlabel level")
		fmt.Println("graph_category netatmo")
		fmt.Println("temp_indor.label indor")
		fmt.Println("temp_outdor.label outdor")
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
	
	//output info
	for _, station := range dc.Stations() {
		for _, module := range station.Modules() {
			if (module.ModuleName == "Indoor"){
				_, data := module.Data()
				for dataType, value := range data {
					//Noise info
					if (len(os.Args) == 2 && os.Args[1] == "-noise" && dataType == "Noise"){
						fmt.Printf("noise.value %v\n", value)
					//CO2 info
					}else if (len(os.Args) == 2 && os.Args[1] == "-co2" && dataType == "CO2"){
						fmt.Printf("co2.value %v\n", value)
					//Pressure info
					}else if (len(os.Args) == 2 && os.Args[1] == "-pressure" && dataType == "Pressure"){
						fmt.Printf("pressure.value %v\n", value)
					//Humidity Indor info
					}else if (len(os.Args) == 2 && os.Args[1] == "-hum" && dataType == "Humidity"){
						fmt.Printf("hum_indor.value %v\n", value)
					//Temperature Indor info
					}else if (len(os.Args) == 2 && os.Args[1] == "-temp" && dataType == "Temperature"){
						fmt.Printf("temp_indor.value %v\n", value)
					} 
				}
			}
			if (module.ModuleName == "Outdoor"){
				_, data := module.Data()
				for dataType, value := range data {
					//Humidity Outdoor info
					if (len(os.Args) == 2 && os.Args[1] == "-hum" && dataType == "Humidity"){
						fmt.Printf("hum_outdor.value %v\n", value)
					//Temperature Outdoor info
					}else if (len(os.Args) == 2 && os.Args[1] == "-temp" && dataType == "Temperature"){
						fmt.Printf("temp_outdor.value %v\n", value)
					} 

				}
			}
		}
	}
    
}

