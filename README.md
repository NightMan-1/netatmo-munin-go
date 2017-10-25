# NetAtmo plug-in for Munin v2.2

## Dependencies
NetAtmo Weather Station  
GoLang 1.x  
Munin-node 2.x

## Installation

Open https://dev.netatmo.com and create first application, then:

~~~sh
go get github.com/NightMan-1/netatmo-api-go
go get github.com/go-ini/ini
git clone https://github.com/NightMan-1/netatmo-munin-go
cd netatmo-munin-go
chmod +x install.sh
./install.sh
#service munin-node restart
~~~

## Credits
Copyright © 2017 [Sergey Gurinovich](mailto:sergey@fsky.info).
