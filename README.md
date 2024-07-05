
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

启动 loki 日志服务系统 **注意各个组件版本必须保持一致**
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
docker run -itd --name promtail_testgo \
-v $PWD/promtail/config.yml:/etc/promtail/config.yml  \
-v $PWD:/var/log   \
grafana/promtail:2.9.2 --config.file=/etc/promtail/config.yml
```

### 假设同一台服务器上面要监控多个日志文件

先在这台服务器上面启动一个 promtail 容器，然后通过修改 promtail/config.yml 新增 job ，然后再容器外部挂载目录，最后重启容器实现

启动容器
```
docker  run -itd --name promtail -v $PWD/promtail/config.yml:/etc/promtail/config.yml  grafana/promtail:3.0.0 --config.file=/etc/promtail/config.yml
```

如下示例  promtail/config.yml 开启2个 job 
```
clients:
  - url: http://172.16.9.116:3100/loki/api/v1/push

positions:
  filename: /tmp/positions.yaml

scrape_configs:
  - job_name: my-java-app
    static_configs:
      - targets:
          - localhost
        labels:
          job: java_ytwl_admin
          __path__: /var/log/xcx_admin/*.log # 监控服务1的日志
  - job_name: xcx-wx-test
    static_configs:
      - targets:
          - localhost
        labels:
          job: go_xcx_wx_test
          __path__: /var/log/xcx_api/test/*.log  # 监控服务1的日志
```

挂载容器目录文件
```

docker exec -it promtail mkdir -p  /var/log/xcx_admin
docker exec -it promtail mount -t volume /var/log/admin /var/log/xcx_admin


挂载 xcx test api 日志目录

docker exec -it promtail mkdir -p  /var/log/xcx_api/test/
docker exec -it promtail mount -t volume /var/log/a.log /var/log/xcx_api/test/a.log
```
重启容器
```
docker restart promtail
```

