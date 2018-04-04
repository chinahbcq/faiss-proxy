FROM daocloud.io/ubuntu:trusty

RUN apt-get update  

RUN apt-get install -y ca-certificates

ENV TZ=Asia/Shanghai

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN mkdir -p /src/bin

RUN mkdir -p /src/conf

RUN mkdir -p /src/log

COPY faiss-proxy /src/bin

COPY conf /src/conf

WORKDIR /src

ENV PATH "$PATH:/src/bin"

