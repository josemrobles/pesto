#!/bin/bash

go get
go build -o /usr/local/bin/worker
supervisord -c /etc/supervisor/supervisord.conf
sleep 10
supervisorctl restart all
