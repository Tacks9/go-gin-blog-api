# 标签模块管理



- 标签列表

```shell
curl --location 'http://127.0.0.1:54001/api/v1/tags?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIxMjM0NTYiLCJleHAiOjE2ODkwNTgwMDUsImlzcyI6ImdvLWdpbi1ibG9nLWFwaSJ9.bdfRvMfjq-rkM7_GbAWJMsc6IeF7YcmM2wBrivEy4Bg'
```

- 新增标签

```shell
curl --location 'http://127.0.0.1:54001/api/v1/tags?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIxMjM0NTYiLCJleHAiOjE2ODkwNTgwMDUsImlzcyI6ImdvLWdpbi1ibG9nLWFwaSJ9.bdfRvMfjq-rkM7_GbAWJMsc6IeF7YcmM2wBrivEy4Bg' \
--form 'name="Golang"' \
--form 'state="1"' \
--form 'created_by="admin"'
```


- 编辑标签

```shell
curl --location --request PUT 'http://127.0.0.1:54001/api/v1/tags/1?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIxMjM0NTYiLCJleHAiOjE2ODkwNTgwMDUsImlzcyI6ImdvLWdpbi1ibG9nLWFwaSJ9.bdfRvMfjq-rkM7_GbAWJMsc6IeF7YcmM2wBrivEy4Bg' \
--form 'name="PHP"' \
--form 'state="1"' \
--form 'modified_by="test"' \
--form 'id="1"'
```

- 删除标签


```shell
curl --location --request DELETE 'http://127.0.0.1:54001/api/v1/tags/1?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIxMjM0NTYiLCJleHAiOjE2ODkwNTgwMDUsImlzcyI6ImdvLWdpbi1ibG9nLWFwaSJ9.bdfRvMfjq-rkM7_GbAWJMsc6IeF7YcmM2wBrivEy4Bg'
```

