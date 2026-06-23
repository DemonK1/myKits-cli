```bash
git remote remove github

# 先确保all远程存在（如果之前有脏配置，先删了重建更干净）
git remote remove all 2>/dev/null

# 新建all远程，先绑GitHub作为fetch源
git remote add all git@github.com:DemonK1/myKits-cli.git

# 给all追加Gitee的推送地址（现在all有两个push目标了）
git remote set-url --add all git@gitee.com:Chosen1uu/myKits-cli.git

# 绑定上游分支
git push --set-upstream all main

# 看远程配置
git remote -v
# 输出
all	git@github.com:DemonK1/myKits-cli.git (fetch)
all	git@github.com:DemonK1/myKits-cli.git (push)
all	git@gitee.com:Chosen1uu/myKits-cli.git (push)

git push all
```



```
https://chat.deepseek.com/share/daep9d0f2do9toha1b
```

```
go get github.com/charmbracelet/bubbletea
```

```
git push all main
```

```
go build -trimpath -ldflags="-s -w" -o app.exe
```



main.go 启动 ==》 主菜单==》照片重命名 ==》 文件夹重命名 ==》读取 Excel ==》目录结构创建 ==》 微信多开==》 数据库启动 

```
mykits-cli/
├── go.mod
├── main.go                   # 主菜单 + 导出功能
├── tools/                    # 你的各个工具独立文件夹
│   ├── photoBatch/
│   │   └── run.go
│   ├── dirBatch/
│   │   └── run.go
│   ├── excelHraderDirs/
│   │   └── run.go
│   ├── wechatMulti/
│   │   └── run.go
│   ├── dbStart/
│   │   └── run.go
└── cmd/                      # 用于导出编译的独立入口（每个工具一个）
    ├── photoCmd/
    │   └── main.go
    ├── dirCmd/
    │   └── main.go
    ...
```

