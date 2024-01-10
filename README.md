### 通过Nexus api 下载maven仓库或者npm仓库中的制品 然后上传到另一个仓库中

### 1.创建本地配置文件`.config`如下
`***Auth`字段值为 `<账户>:<密码>` base64 加密 以下命令可以获得
```bash
echo -n "admin:12345" | base64 

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

# 任务执行间隔 单位秒
timeStep: 20

# http监听端口
port: 18090
```

### 2.二进制启动
```bash
./NexusSync
```

### 3.本地文件
程序启动后会在本地启动以下文件
- `download`目录：本地存储远端仓库下载的制品文件
- `nexus.db`：sqlite数据库,存储制品文件的下载上传状态

### 4.访问健康检查接口
```bash
curl 127.0.0.1:18090/health
```