
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

在 promtail/config.yml 写2个 job ，如下示例  promtail/config.yml 开启2个 job 
```
clients:
  - url: http://localhost:3100/loki/api/v1/push

positions:
  filename: /tmp/positions.yaml

scrape_configs:
  - job_name: my-java-app  #服务1
    static_configs:
      - targets:
          - localhost
        labels:
          job: java_ytwl_admin
          __path__: /var/log/xcx_admin/*.log  #服务1的日志
  - job_name: xcx-wx-test  #服务2
    static_configs:
      - targets:
          - localhost
        labels:
          job: go_xcx_wx_test
          __path__: /var/log/xcx_api/test/TestWXmini.log  服务2的日志文件

```
启动容器（挂载 2个日志文件，对应上面 config.yml 中2个服务的 __PATH__ 位置）

```
docker run -itd  --name promtail \
-v /home/promtail/config.yml:/etc/promtail/config.yml \
-v /var/lib/docker/containers/1bffbff5d78962a5f79c9f42032117e96b7a87f54167d41c52432b74b0a8fc40:/var/log/xcx_admin \  
-v /www/wwwlogs/go/TestWXmini.log:/var/log/xcx_api/test/TestWXmini.log \
grafana/promtail:3.0.0 --config.file=/etc/promtail/config.yml
```
查看容器日志信息
```
docker logs -f promtail
```
出现下面内容则成功
![image](https://github.com/ituserxxx/Loki/assets/66945660/25d66e89-137f-40b2-9fd4-62f84942bdd1)

