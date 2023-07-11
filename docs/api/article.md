# 文章模块管理

- 新增文章

```shell
curl --location 'http://127.0.0.1:54001/api/v1/articles?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIxMjM0NTYiLCJleHAiOjE2ODkwNDY2MTksImlzcyI6ImdvLWdpbi1ibG9nLWFwaSJ9.dz__PHakmdneNa9nPtgwCs4UJYYYN4p0ea83HO51jt4' \
--form 'tag_id="1"' \
--form 'title="文章测试发布"' \
--form 'desc="这是一篇测试文章"' \
--form 'content="测试啊测试啊测试啊TEST"' \
--form 'created_by="admin"' \
--form 'state="1"' \
--form 'cover_image_url="test"'
```


- 编辑文章

```shell
curl --location --request PUT 'http://127.0.0.1:54001/api/v1/articles/1?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIxMjM0NTYiLCJleHAiOjE2ODkwNDY2MTksImlzcyI6ImdvLWdpbi1ibG9nLWFwaSJ9.dz__PHakmdneNa9nPtgwCs4UJYYYN4p0ea83HO51jt4' \
--form 'id="1"' \
--form 'tag_id="1"' \
--form 'title="文章测试发布(更新)"' \
--form 'desc="这是一篇测试文章(更新)"' \
--form 'content="测试啊测试啊测试啊TEST(更新)"' \
--form 'modified_by="admin"' \
--form 'state="1"' \
--form 'cover_image_url="test"'
```

- 删除文章


```shell
curl --location --request DELETE 'http://127.0.0.1:54001/api/v1/articles/2?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIxMjM0NTYiLCJleHAiOjE2ODkwNDY2MTksImlzcyI6ImdvLWdpbi1ibG9nLWFwaSJ9.dz__PHakmdneNa9nPtgwCs4UJYYYN4p0ea83HO51jt4'
```