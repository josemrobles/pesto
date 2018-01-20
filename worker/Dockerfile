FROM golang:1.9

WORKDIR /opt/go/pesto/worker

RUN apt-get update &&\
    apt-get install -y supervisor ;\
    mkdir -p /var/log/supervisor

COPY supervisor.conf /etc/supervisor/conf.d/supervisor.conf
ADD . /opt/go/pesto/worker
CMD "./entrypoint.sh"
