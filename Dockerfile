FROM phusion/baseimage:master

CMD ["/sbin/my_init"]

WORKDIR /root

# 在基础镜像 /etc/my_init.d/ 中的脚本会在容器启动时执行
COPY main /root/main

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
