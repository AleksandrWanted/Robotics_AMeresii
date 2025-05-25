#!/bin/bash
ip=$1

env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/bin/ameresii_smart_home github.com/AleksandrWanted/AMeresii_SMART_HOME/cmd/server
ssh root@$ip "mkdir /opt/ameresii_smart_home"
ssh root@$ip "mkdir /opt/ameresii_smart_home/bin"
ssh root@$ip "rm -f /opt/ameresii_smart_home/bin/ameresii_smart_home"
scp build/bin/ameresii_smart_home root@$ip:/opt/ameresii_smart_home/bin/ameresii_smart_home
scp deploy/ameresii_smart_home.service root@$ip:/etc/systemd/system/ameresii_smart_home.service
ssh root@$ip "mkdir /opt/ameresii_smart_home/etc"
ssh root@$ip "mkdir /opt/ameresii_smart_home/etc/conf.d"
scp .env root@$ip:/opt/ameresii_smart_home/etc/conf.d/smart_home-env
ssh root@$ip "rm -rf /home/ameresii_smart_home/configs"
ssh root@$ip "mkdir -p /home/ameresii_smart_home/configs"
scp configs/main.yaml root@$ip:/home/ameresii_smart_home/configs/main.yaml
ssh root@$ip "systemctl daemon-reload"
ssh root@$ip "systemctl enable ameresii_smart_home"
ssh root@$ip "systemctl restart ameresii_smart_home"