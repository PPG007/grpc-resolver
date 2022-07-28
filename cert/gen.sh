#!/usr/bin/zsh

# 首先生成私钥
openssl genrsa -out cakey.pem
# 创建请求文件
openssl req -new -key cakey.pem -out rootCA.csr
# 自签署
openssl x509 -req -in rootCA.csr -signkey cakey.pem -out x509.crt

# 首先生成私钥
openssl genrsa -out prikey.pem
# 创建请求文件
openssl req -new -key prikey.pem -out local.csr
# 签署证书请求
openssl x509 -req -in local.csr -CA x509.crt -CAkey cakey.pem -out localhost.crt -CAcreateserial -extfile extfile.cnf
