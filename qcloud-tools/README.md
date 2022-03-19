# 腾讯云小工具

用 `letsencrypt` 来生成免费证书，并同步至腾讯云 cdn / ecdn / CLB 。因为 `letsencrypt` 证书有效期只有三个月，所以也提供一个定时器及时更新  

## 部署

见 [deployment.md](deployment.md)

## 目录结构
```
- config #配置文件
- cmd/main.go # 程序入口
- certificate/ # 证书处理核心代码
    - issue.go # 通过 shell 脚本的方式签发 `letsencrypt` 证书
    - issue_history.go # 证书签发历史，还有一些取数方法，一起放在这里了
    - sync.go # 同步证书至 cdn / ecdn
    - task.go # 定时器
- core/
    - config/ # 配置文件解析
    - db/ # db 操作相关代码
    - utils.go # 工具类
- services/ # http 接口相关代码
```

## Todo
- [ ] web 完善
    - [ ] 鉴权功能
    - [ ] 证书下载
    - [ ] 手动更新
- [ ] 监控
- [ ] 单元测试
- [ ] 定时器策略修改，不要扫全表
- [ ] 功能完善，支持更多平台