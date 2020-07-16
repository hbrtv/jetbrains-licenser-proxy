FROM centos:8

ENV TZ=Asia/Shanghai

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY bin/jetbrains-licenser /opt/bin/jetbrains-licenser

COPY jetbrains-licenser-proxy /opt/bin/jetbrains-licenser-proxy

COPY tmpl/statistics.html /opt/tmpl/statistics.html

EXPOSE 80

VOLUME /opt/log

CMD ["/opt/bin/jetbrains-licenser-proxy"]
