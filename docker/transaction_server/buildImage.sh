
# 編譯後上傳到 dockerhub liuleo=dockerhub上面的帳號, 請自行帶入自己的
# docker build -t 192.168.100.20:9001/acp-rd/transaction_server:v1 -f ./Dockerfile  ../../../
# docker build -t liuleo/transaction_server:v1 . --no-cache

# 支援 linux / mac / 樹莓派
#docker buildx build --push -t 192.168.100.20:9001/acp-rd/transaction_server:v1 --platform linux/amd64,linux/arm64,linux/arm/v7 -f ./Dockerfile  ../../../

# 交叉編譯 支援 linux / mac
docker buildx build --push -t close0818/transaction_server:v1 --platform linux/amd64,linux/arm64 -f ./Dockerfile  ../../../
#docker buildx build --push -t close0818/transaction_server:v2 --platform linux/amd64,linux/arm64 -f ./Dockerfile  ../../../


#docker buildx build -t liuleo/transaction_server:v1 --platform linux/amd64,linux/arm64 -f ./Dockerfile  ../../../

# 推上去
#docker push 192.168.100.20:9001/acp-rd/transaction_server:v1

# 編譯後上傳到自己筆電的倉庫
#docker build -t test:v6 -f ./Dockerfile  ../../../
