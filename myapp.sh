#!/usr/bin/env bash

# resolve links - $0 may be a softlink
PRG="$0"

while [ -h "$PRG" ]; do
  ls=`ls -ld "$PRG"`
  link=`expr "$ls" : '.*-> \(.*\)$'`
  if expr "$link" : '/.*' > /dev/null; then
    PRG="$link"
  else
    PRG=`dirname "$PRG"`/"$link"
  fi
done

# Get standard environment variables
PRGDIR=`dirname "$PRG"`
cd ${PRGDIR}

APP_NAME="myapp"
PID_PATH="/var/run/${APP_NAME}.pid"
LOG_PATH="/var/log/${APP_NAME}"
PYENV_PATH=$(/root/.poetry/bin/poetry debug|grep Path|awk '{print $2}')

function get_setting_key()
{
  APP_CONFIG_CMD="${PYENV_PATH}/bin/dynaconf -i myapp.conf.config.settings list --key "
  ${APP_CONFIG_CMD} $1 |grep -v environment|awk '{print $2}'
}

# Setting the binding host and port for the service
BIND_HOST=`get_setting_key "default.host" | sed $'s/\'//g'`
BIND_PORT=`get_setting_key "default.port"`

if [ -n ${BIND_HOST} ];then
    HOST=${BIND_HOST}
fi
if [ -n ${BIND_PORT} ];then
    PORT=${BIND_PORT}
fi

BIND_ADDR="${HOST}:${PORT}"

if [ ! -d "${LOG_PATH}" ]; then
mkdir -p "${LOG_PATH}"
fi

function stop()
{
    if [ -f "${PID_PATH}" ]; then
        echo "kill process `cat ${PID_PATH}`"
        kill `cat ${PID_PATH}`
        sleep 2
    fi
}

function start()
{
    # 查看逻辑CPU的个数
    NUM_CORES=`cat /proc/cpuinfo| grep "processor"| wc -l`
    # 使用gunicorn管理Uvicorn进程。使用自定义的worker启动Uvicorn。
    nohup ${PYENV_PATH}/bin/gunicorn myapp.main:app \
        -b ${BIND_ADDR} \
        -w $((NUM_CORES + 1)) \
        -k myapp.worker.MyUvicornWorker \
        --timeout 60 \
        --log-level debug \
        --access-logfile ${LOG_PATH}/app.log \
        --error-logfile ${LOG_PATH}/app.log \
        -p ${PID_PATH} &
}

function debug()
{
    # 查看逻辑CPU的个数
    NUM_CORES=`cat /proc/cpuinfo| grep "processor"| wc -l`
    ${PYENV_PATH}/bin/gunicorn myapp.main:app \
        -b ${BIND_ADDR} \
        -w $((NUM_CORES + 1)) \
        -k myapp.worker.MyUvicornWorker \
        --timeout 60 \
        --log-level debug \
        -p ${PID_PATH}
}

function usage()
{
    echo "Usage: ${PRG} ( commands ... )"
    echo "commands:"
    echo "  start        Start ${APP_NAME} service"
    echo "  debug        Start ${APP_NAME} service in terminal"
    echo "  stop         Stop ${APP_NAME} service"
    echo "  restart      Restart ${APP_NAME} service"
    echo "  reload       reload ${APP_NAME} service"
    echo "  status       Show ${APP_NAME} service status"
    echo "  help         Show help"
}

function status()
{
    if [ -f "${PID_PATH}" ]; then
        pid=`cat ${PID_PATH}`  # 获取 pid
        if [ "x${pid}" != "x" ]; then
            c=`ps -ef | grep ${pid} | wc -l`  # 获取 pid 进程数
            if [ "x$c" != "x0" ]; then
                # pid 进程存在
                echo "service ${APP_NAME} is running, pid: ${pid}"
            else
                # pid 进程不存在
                echo "service ${APP_NAME} is stoped"
            fi
        else
            # pid 文件是空的
            echo "service ${APP_NAME} is stoped"
        fi
    else
        # pid 文件不存在
        echo "service ${APP_NAME} is stoped"
    fi
}

function reload()
{
    if [ -f "${PID_PATH}" ]; then
        pid=`cat ${PID_PATH}`
        if [ "x${pid}" != "x" ]; then
            c=`ps -ef | grep ${pid} | wc -l`
            if [ "x$c" != "x0" ]; then
                # pid 进程存在
                echo "reload service, pid: `cat ${PID_PATH}`"
                kill -HUP ${pid}
            else
                # pid 进程不存在
                echo "service ${APP_NAME} is not running"
            fi
        else
            # pid 文件是空的
            echo "service ${APP_NAME} is not running"
        fi
    else
        # pid 文件不存在
        echo "service ${APP_NAME} not running"
    fi
}

case $1 in
    start)
        start
    ;;
    debug)
        debug
    ;;
    stop)
        stop
    ;;
    restart)
        stop
        start
    ;;
    reload)
        reload
    ;;
    status)
        status
    ;;
    *) usage
    ;;
esac

exit 0
