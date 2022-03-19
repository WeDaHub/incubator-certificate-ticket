FROM actors315/go

MAINTAINER actors315 <actors315@gmail.com>

COPY . /usr/local/src/qcloud-tools/

# 安装 qcloud-tools
RUN . ~/.bashrc && cd /usr/local/src/qcloud-tools \
    && make clean && make build \
    && mv /usr/local/src/qcloud-tools/bin/* /usr/local/qcloud-tools/ \
    && make clean \
# 配置
    && mv /usr/local/src/qcloud-tools/config/config.simple.yaml /usr/local/qcloud-tools/config/config.yaml \
    && mv /usr/local/src/qcloud-tools/config/issue-template.tpl /usr/local/qcloud-tools/config/issue-template.tpl \
    && mv /usr/local/src/qcloud-tools/Dockerstart /start \
    && chmod +x /start

EXPOSE 80

WORKDIR /usr/local/qcloud-tools/

CMD ["/start"]