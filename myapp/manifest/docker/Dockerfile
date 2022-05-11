# gf docker -t my-goframe
FROM loads/alpine:3.8

LABEL maintainer="xiaobin.zhao@cstack.io"

###############################################################################
#                                INSTALLATION
###############################################################################

# 设置在容器内执行时当前的目录
ENV WORKDIR /opt/my-goFrame

# 添加应用可执行文件，并设置执行权限
COPY ./bin/linux_amd64/mygf-app   $WORKDIR/mygf-app
RUN chmod +x $WORKDIR/mygf-app

# TODO: 可能是gf build的问题，造成资源没有被pack,需要继续拷贝
COPY ./resource/i18n   $WORKDIR/resource/i18n
COPY ./resource/public   $WORKDIR/resource/public

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./mygf-app
