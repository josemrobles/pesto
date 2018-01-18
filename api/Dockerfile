FROM golang:1.9

WORKDIR /opt/go/pesto/api

RUN apt-get update &&\
    apt-get install -y supervisor ;\
    mkdir -p /var/log/supervisor

COPY supervisor.conf /etc/supervisor/conf.d/supervisor.conf
ADD . /opt/go/pesto/api
CMD "./entrypoint.sh"
