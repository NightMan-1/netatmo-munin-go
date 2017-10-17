#!/bin/bash

if [ ! -f /usr/sbin/munin-node ]; then
    echo "File $FILE does not exist."
    exit
fi

echo -n "Compilling..."
go build netatmo-munin-plugin.go
if [ ! -f ./netatmo-munin-plugin ]; then
    echo "error"
    exit
else
    echo "ok"
fi
strip ./netatmo-munin-plugin

echo "Configuration ..."
cp netatmo-munin-plugin /usr/share/munin/plugins
ln -s -f /usr/share/munin/plugins/netatmo-munin-plugin /etc/munin/plugins/netatmo-munin-plugin

if [ ! -f /etc/munin/netatmo.cfg ]; then
    echo "#go https://dev.netatmo.com" >/etc/munin/netatmo.cfg
    echo -e "\nClientID = \nClientSecret = \nUsername = \nPassword = \n" >>/etc/munin/netatmo.cfg
    nano /etc/munin/netatmo.cfg
fi

echo All done! Now run "service munin-node restart" and restart munin
