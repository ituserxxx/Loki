
grafana 官网示例版本搭建请查看：https://github.com/ituserxxx/grafana_loki_log_install_demo.git

loki  日志系统搭建(以 go 程序为例)

### 启动 go 服务

编译可执行文件
```
go build -o testgo main.go
```

启动 go 服务
```
nohup ./testgo &
```

启动 loki 日志服务系统
```
docker compose -f docker-compose.yaml up
```

启动 Promtail(一个代理，它将本地日志的内容发送到私有的Grafana Loki实例)

在同一台服务器则需要加入 loki 示例的网络中,且需要把配置文件 (promtail/conifg.yml)里面的 url 写成容器名称
```
docker run  --name promtail_testgo \
--network=loki_loki \                   # 
-v $PWD/promtail/config.yml:/etc/promtail/config.yml  \
-v $PWD:/var/log   \
grafana/promtail:2.9.2
```

在不同服务器上面,则需要在配置文件 (promtail/conifg.yml)里面的 url 写成 Loki 实例 Ip : port,然后启动
```
docker run  --name promtail_testgo \
-v $PWD/promtail/config.yml:/etc/promtail/config.yml  \
-v $PWD:/var/log   \
grafana/promtail:2.9.2
```

