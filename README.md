### 通过Nexus api 下载maven仓库或者npm仓库中的制品 然后上传到另一个仓库中

### 1.创建本地配置文件`.config`如下
`***Auth`字段值为 `<账户>:<密码>` base64 加密 以下命令可以获得
```bash
echo -n "admin:12345" | base64 

# return
# YWRtaW46MTIzNDU=
```
```
# downloadMaven
downloadMavenRepositoryUrl: "http://172.30.86.46:18081"
downloadMavenRepositoryName: "maven-proxy-148-ali"
downloadMavenRepositoryAuth: "YWRtaW46SHlkZXZAbmV4dXMyMDIz"

#  uploadMaven
uploadMavenRepositoryUrl: "http://172.30.86.46:18081"
uploadMavenRepositoryName: "test-upload"
uploadMavenRepositoryAuth: "YWRtaW46SHlkZXZAbmV4dXMyMDIz"

# downloadNpm
downloadNpmRepositoryUrl: "http://172.30.84.90:8081"
downloadNpmRepositoryName: "npm-local"
downloadNpmRepositoryAuth: "YWRtaW46WnlqY0AyMDIx"

# uploadNpm
uploadNpmRepositoryUrl: "http://172.30.86.46:18081"
uploadNpmRepositoryName: "test-npm-upload"
uploadNpmRepositoryAuth: "YWRtaW46SHlkZXZAbmV4dXMyMDIz"
```