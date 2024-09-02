### 通过Nexus api 下载maven仓库或者npm仓库中的制品 然后上传到另一个仓库中

### 1.创建本地配置文件`.config`如下
`uploadRepositoryAuth`字段值为 `<账户>:<密码>` base64 加密 以下命令可以获得
```bash
echo -n "admin:12345" | base64 

# YWRtaW46MTIzNDU=
```
```
repositorySyncTask:
  # maven同步
  - downRepositoryUrl: "http://172.30.86.136:8081"
    downRepositoryName: "sync-maven-public"
    uploadRepositoryUrl: "http://10.147.235.204:8081"
    uploadRepositoryName: "inner-maven-public"
    uploadRepositoryAuth: "YWRtasadasd5dEBuZXh1c0AyMDIz"
    repositoryType: "maven2"
    
  # npm 同步
  - downRepositoryUrl: "http://172.30.86.136:8081"
    downRepositoryName: "sync-npm-public"
    uploadRepositoryUrl: "http://10.147.235.204:8081"
    uploadRepositoryName: "inner-npm-public"
    uploadRepositoryAuth: "YWRtaW46WasdasBuZXh1c0AyMDIz"
    repositoryType: "npm"

# 任务执行间隔 单位秒
timeStep: 30
# 端口
port: 18090
# 数据库路径
dbPath: "./db/nexus.db"
# 下载目录
downloadPath: "testdownload"

# http监听端口
port: 18090
```

### 2.二进制启动
```bash
./NexusRepositorySync
```

### 3.本地文件
程序启动后会在本地启动以下文件
- `download`目录：本地存储远端仓库下载的制品文件
- `nexus.db`：sqlite数据库,存储制品文件的下载上传状态

### 4.访问健康检查接口
```bash
curl 127.0.0.1:18090/health
```