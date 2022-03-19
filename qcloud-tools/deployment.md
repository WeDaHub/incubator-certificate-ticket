## 准备

需要有一个 MySQL 数据库，初始化 sql 见 `config/init.sql`，根据实际需要调整。

## 主机部署

### 配置项

config 目录下,重建 config.simple.yaml 为 config.yaml

修改数据库对应配置

```
db:
  host: "localhost"
  port: 3306
  database: "qcloud-tools"
  user: "db_qcloud"
  password: "58117aec3b3252a97be0"

```

### 编译 && 运行

```
make cert-monitor

./bin/cert-monitor
```

## Docker 部署

### 配置项

通过环境变量的方式注入，支持如下变量

```
MYSQL_HOST=dev.local.mysql
MYSQL_HOST=3306
MYSQL_DATABASE=qcloud-tools
MYSQL_USER=db_qcloud
MYSQL_PASSWORD=58117aec3b3252a97be0
HTTP_PORT=80
```

### 生成镜像
```
docker build -t qcloud-tools .
```

### 启动

```
docker run --name my-qcloud-tools --env-file=./.env.dev -p 80:80 -d qcloud-tools:latest
```


## 推荐使用腾讯云 CloudBase

[![](https://main.qcloudimg.com/raw/67f5a389f1ac6f3b4d04c7256438e44f.svg)](https://console.cloud.tencent.com/tcb/env/index?action=CreateAndDeployCloudBaseProject&appUrl=https://github.com/actors315/incubator-certificate-ticket&workDir=qcloud-tools)
