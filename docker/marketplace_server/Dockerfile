# build stage
FROM golang:1.21-alpine AS build-env

# 拷貝需要的專案
ADD barketplace_server /src/barketplace_server

# 編譯server執行檔
RUN cd /src/barketplace_server/cmd/marketplace_server && go build 

# final stage
#FROM alpine
FROM ubuntu
WORKDIR /app
# 拷貝編譯好的執行檔案 server
COPY --from=build-env /src/barketplace_server/cmd/marketplace_server/marketplace_server /app/

# 拷貝 server 的 等待db啟動延遲執行腳本
#COPY --from=build-env /src/marketplace_server/linux/wait-for-it.sh /app/
#RUN chmod +x /app/wait-for-it.sh

# 安裝 supervisor
#RUN apt-get update
#RUN apt-get upgrade -y
#RUN apt-get install -y supervisor

RUN mkdir -p /var/run/sshd
#RUN mkdir -p /var/log/supervisor
RUN mkdir -p /app/log

#COPY --from=build-env /src/marketplace_server/linux/supervisor.conf /etc/supervisor/conf.d/supervisord.conf
#COPY --from=build-env /src/marketplace_server/linux/start-container /usr/local/bin/start-container
#RUN chmod +x /usr/local/bin/start-container

#server 開放的port 
EXPOSE 8888


ENTRYPOINT ["./marketplace_server"]