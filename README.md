### 通过Nexus api 下载maven仓库或者npm仓库中的制品 然后上传到另一个仓库中

### 1.创建本地配置文件`.config`如下
`uploadRepositoryAuth` `downRepositoryAuth`字段值为 `<账户>:<密码>` base64 加密 以下命令可以获得
```bash
echo -n "admin:12345" | base64 

# YWRtaW46MTIzNDU=
```
```
repositorySyncTask:
  # maven同步
  - taskName: "maven-sync-1"                              # 不能重复
    downRepositoryUrl: "http://172.30.86.136:8081"
    downRepositoryName: "sync-maven-public-1"
    downRepositoryAuth: "YWRtasadasd5dEBuZXh1c0AyMDIz"  # optional
    uploadRepositoryUrl: "http://10.147.235.204:8081"
    uploadRepositoryName: "inner-maven-public-1"
    uploadRepositoryAuth: "YWRtasadasd5dEBuZXh1c0AyMDIz"
    repositoryType: "maven2"
  - taskName: "maven-sync-2"                              # 不能重复
    downRepositoryUrl: "http://172.30.86.136:8081"
    downRepositoryName: "sync-maven-public-2"
    downRepositoryAuth: "YWRtasadasd5dEBuZXh1c0AyMDIz"  # optional
    uploadRepositoryUrl: "http://10.147.235.204:8081"
    uploadRepositoryName: "inner-maven-public-2"
    uploadRepositoryAuth: "YWRtasadasd5dEBuZXh1c0AyMDIz"
    repositoryType: "maven2"
    
  # npm 同步
  - taskName: "npm-sync-1"
    downRepositoryUrl: "http://172.30.86.136:8081"
    downRepositoryName: "sync-npm-public-1"
    uploadRepositoryUrl: "http://10.147.235.204:8081"
    uploadRepositoryName: "inner-npm-public-1"
    uploadRepositoryAuth: "YWRtaW46WasdasBuZXh1c0AyMDIz"
    repositoryType: "npm"
  - taskName: "npm-sync-2"
    downRepositoryUrl: "http://172.30.86.136:8081"
    downRepositoryName: "sync-npm-public-2"
    uploadRepositoryUrl: "http://10.147.235.204:8081"
    uploadRepositoryName: "inner-npm-public-2"
    uploadRepositoryAuth: "YWRtaW46WasdasBuZXh1c0AyMDIz"
    repositoryType: "npm"
    
# 守护模式下任务执行间隔 单位秒
timeStep: 30
# 端口
port: 18090
# sqlite数据库路径
dbPath: "./db/nexus.db"
# 下载目录
downloadPath: "testdownload"

# http监听端口
port: 18090
```
### 2.创建本地文件目录
根据配置中的`downloadPath`和`dbPath`字段创建本地文件目录
```bash
mkdir -p testdownload
touch ./db/nexus.db
```
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
# 检索并下载文件至本地目录
./NexusRepositorySync sd

# 检索仓库文件至本地数据库中
./NexusRepositorySync search

# 下载仓库文件
./NexusRepositorySync download

# 上传仓库文件
./NexusRepositorySync upload
```

