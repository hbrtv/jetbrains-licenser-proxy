FROM centos:7

ENV TZ=Asia/Shanghai

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY bin/jetbrains-licenser /opt/bin/jetbrains-licenser

COPY jetbrains-licenser-proxy /opt/bin/jetbrains-licenser-proxy

EXPOSE 80

VOLUME /opt/log

ENTRYPOINT ["/opt/bin/jetbrains-licenser-proxy"]
