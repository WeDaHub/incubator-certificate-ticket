# 腾讯云小工具

用 `letsencrypt` 来生成免费证书，并同步至腾讯云 cdn / CLB 。因为 `letsencrypt` 证书有效期只有三个月，所以也提供一个定时器及时更新  

## 云函数更新证书

- [更新 ELB 证书](syn-certificate-lb)
- [更新 CDN 证书](syn-certificate-cdn)

## 云托管
[qcloud-tools](qcloud-tools)  

提供定时器监控证书情况，快过期时发起更新 

提供管理 API  

## 静态页面

[pages](pages)

提供前端页面
