# 杂项


- 文档

```shell
http://127.0.0.1:54001/swagger/index.html
```

- 测试

```shell
curl --location 'http://127.0.0.1:54001/test'
```

- 项目二维码

```shell
curl --location --request POST 'http://127.0.0.1:54001/qrcode'
```

- 上传图片

```shell
curl --location 'http://127.0.0.1:54001/upload?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIxMjM0NTYiLCJleHAiOjE2ODkxMzQ1NTUsImlzcyI6ImdvLWdpbi1ibG9nLWFwaSJ9.hBtWF4CiXeO3IBwtwbfhNYw16ccZl516cS704FTbE7s' \
--form 'image=@"/xxx.jpg"'
```



