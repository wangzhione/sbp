# sbp

Simple Beautiful Package 高频使用的 Go 代码工具包合集

`sbp` is a universal utility collection for Go, it complements offerings such as Boost, Better std, Helper tools.

这个库为工程实战服务的, 重度依赖 Go slog 日志库, 纪录线上业务 内部状态. 目前线上几个项目用的还凑合.

## Table of Contents

- [Introduction](#Introduction)
- [Catalogs](#Catalogs)
- [Releases & How To Use](#releases--how-to-use)
- [License](#License)

## Introduction

`sbp` is a universal utility collection for Go, it complements offerings such as Boost, Better std, Helper tools. It is migrated from the 实战项目使用 util or tool 体验提升.

We depend on the same code(this repo) in our production environment.

## Catalogs 粗略介绍

* [chain](https://github.com/wangzhione/sbp/tree/master/chain): trace and log
* [util](https://github.com/wangzhione/sbp/tree/master/util): Utilities Useful across Domains
* [helper](https://github.com/wangzhione/sbp/tree/master/helper): helper redis mysql safego local cahce
* [structs](https://github.com/wangzhione/sbp/tree/master/structs): Data Structures or Collection

> 设计者注: 通常 **util** 与业务无关的，可以独立出来，可供其他项目使用通用代码集。方法通常是 public static; **tool** 可以与某些业务有关，通用性限于某几个业务类之间; **helper** 通常与业务相关. 随后是否加 s, 不加 s 看个人喜好了. 

## Releases & How To Use

You can use `go get -u github.com/wangzhione/sbp@master` to get or update `sbp`.

## License

`sbp` is licensed under the terms of the MIT License. See [LICENSE](LICENSE) for more information.

欢迎喜欢用的朋友, 补充常用 package 代码集, 或者发评论提思路, 主动帮忙添加.

## [可选] 扩展阅读

- [Effective Go](https://golang.org/doc/effective_go)
- [Pingcap General advice](https://pingcap.github.io/style-guide/general.html)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## [可选] Visual Studio Code Golang 

IDE 开发 Golang 基础配置在全局 settings.json 加入和 go env 有关配置, 用于控制 go import 和 go test 相关行为

```JSON
{
    "go.toolsManagement.autoUpdate": true,
    "go.testEnvVars": {
        
    },
    "go.testFlags": [
        "-v", "-count=1"
    ],
    "gopls": {
        "ui.importShortcut": "both",
        "formatting.gofumpt": true,
        "ui.semanticTokens": true,
    },
    "[go]": {
        "editor.codeActionsOnSave": {
            "source.organizeImports": "explicit" // 仅在显式保存时触发
        },
        "editor.formatOnSave": true // 可选：启用自动格式化
    },
    "go.testTimeout": "120s",
}
```

[可选] 本地 .vscode/launch.json 添加相关 F5 启动 main 配置

```json
{
    // 使用 IntelliSense 了解相关属性。 
    // 悬停以查看现有属性的描述。
    // 欲了解更多信息，请访问: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}",

            "args": [
                
            ]
        }
    ]
}
```

# One Package 

long long ago 在公司用 gopkg 项目, 有一股啰嗦妈宝氛围扑面而来. 为了方便自己用起来舒服, 构造 one package 项目用于 Go 软件工程使用.

## How To Use

You can use go get -u github.com/wangzhione/sbp@latest to get or update onepkg

## 后记

题破山寺后禅院
常建 〔唐代〕

清晨入古寺，初日照高林。
曲径通幽处，禅房竹木深。
山光悦鸟性，潭影空人心。
万籁此俱寂，惟余钟磬音。

