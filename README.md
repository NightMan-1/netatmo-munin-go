# NetAtmo plug-in for Munin v3.0

[![Build Status](https://travis-ci.org/NightMan-1/netatmo-munin-go.svg?branch=master)](https://travis-ci.org/NightMan-1/netatmo-munin-go)
[![GitHub license](https://img.shields.io/github/license/NightMan-1/netatmo-munin-go.svg)](https://github.com/NightMan-1/netatmo-munin-go/blob/master/LICENSE.md)

## Dependencies
NetAtmo Weather Station  
GoLang 1.x  
Munin-node 2.x

## Installation

Open https://dev.netatmo.com and create first application, then:

~~~sh
go get -u github.com/exzz/netatmo-api-go
go get -u github.com/BurntSushi/toml
git clone https://github.com/NightMan-1/netatmo-munin-go
cd netatmo-munin-go
chmod +x install.sh
./install.sh
#service munin-node restart
~~~

## Sample data

![Temperature](sample/temp.png)
![CO2](sample/co2.png)
![Humidity](sample/humidity.png)
![Pressure](sample/pressure.png)
![Rain](sample/rain.png)



## Credits
Copyright © 2017 [Sergey Gurinovich](mailto:sergey@fsky.info).
