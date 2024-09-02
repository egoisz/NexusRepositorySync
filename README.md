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
### 2.创建本地文件目录
根据配置中的`downloadPath`和`dbPath`字段创建本地文件目录

### 3.守护模式启动(执行目录下需存在上述配置文件)
```bash
./NexusRepositorySync
```
### 4.访问健康检查接口
```bash
curl 127.0.0.1:18090/health
```
### 5.使用子命令进行单次下载或者上传
```bash
# 下载仓库文件
./NexusRepositorySync download

# 上传仓库文件
./NexusRepositorySync upload
```

