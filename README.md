# sbp

Simple Beautiful Package 高频使用的 Go 代码工具包合集

`sbp` is a universal utility collection for Go, it complements offerings such as Boost, Better std, Helper tools.

## Table of Contents

- [Introduction](#Introduction)
- [Catalogs](#Catalogs)
- [Releases](#Releases)
- [How To Use](#How-To-Use)
- [License](#License)

## Introduction

`sbp` is a universal utility collection for Go, it complements offerings such as Boost, Better std, Helper tools. It is migrated from the [![Go Reference](https://pkg.go.dev/badge/github.com/bytedance/gopkg.svg)](https://pkg.go.dev/github.com/bytedance/gopkg) `gopkg` 主要对其相关实战项目使用体验提升, 最初用 gopkg 感觉是里面代码水平层次不齐, 实战用起来很别扭, 逐渐想在实战中雕琢耐用.

We depend on the same code(this repo) in our production environment.

## Catalogs

* [localcache](https://github.com/wangzhione/sbp/tree/master/localcache): Caching Mechanism
* [structs](https://github.com/wangzhione/sbp/tree/master/structs): Data Structures or Collection
* [util](https://github.com/wangzhione/sbp/tree/master/util): Utilities Useful across Domains

> 设计者注: 通常 **util** 与业务无关的，可以独立出来，可供其他项目使用通用代码集。方法通常是 public static; **tool** 可以与某些业务有关，通用性限于某几个业务类之间; **helper** 通常与业务相关. 随后是否加 s, 不加 s 看个人喜好了. 

## Releases

`sbp` recommends users to "live-at-head" (update to the latest commit from the main branch as often as possible).
We develop at `develop` branch and will only merge to `master` when `develop` is stable.

## How To Use

You can use `go get -u github.com/wangzhione/sbp@master` to get or update `sbp`.

## License

`sbp` is licensed under the terms of the MIT License. See [LICENSE](LICENSE) for more information.

欢迎喜欢用的朋友, 补充常用 package 代码集, 或者发评论提思路, 主动帮忙添加.

## 扩展阅读

- [Effective Go](https://golang.org/doc/effective_go)
- [Pingcap General advice](https://pingcap.github.io/style-guide/general.html)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

![](god.webp)

## 拓展配置

我本身用的是 Visual Studio Code 简单说一下, 用这个 IDE 开发 Golang 基础配置

在全局 settings.json 加入和 go env 有关配置, 用于控制 go import 和 go test 相关行为

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

# sbp 前传 One Package 

看见公司在用的 gopkg 项目, 有一股啰嗦妈宝氛围扑面而来. 为了方便自己提炼内在细节, 构造当前 one package 项目用于学习和工程使用.

## How To Use

You can use go get -u github.com/wangzhione/sbp@latest to get or update onepkg

## 后记

题破山寺后禅院

常建 〔唐代〕

清晨入古寺，初日照高林。

曲径通幽处，禅房花木深。(曲 一作：竹)

山光悦鸟性，潭影空人心。

万籁此都寂，但余钟磬音。(都寂 一作：俱寂；但余 一作：惟余)

```
                                          0@
                                         @@@0
                                       :0@@@@0
                                      L@@@@@@@0
                                     00@@@@@@@@0;
                                    0@@@@@@@@@@@@G
                                   @@@@@@@@@@@@@@@0
                                  @@@@@@@@@@@@@@@@00
                                .0@@@@@@@@@@@@@@@@@@0
                               10@@@@@@@@@@@@@@@@@@@@@
                              C@@@@@@@@@@@@@@@@@@@@@@@0,

                           0@@@@@@@@@@@@@@@@@@@@@@@@@@@@@0@
                          0@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@0
                        :@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@0
                       L@@@@@@@@@@@0@@@0@@@0@@00@@@@@@@@@@@@00
                      0@@@@@@0@0@00G.            :C0@@00@@@@@@@i
                     0@@@@0@0:                          L@@@0@@0C
                    0@@00@                                  @@0@@0
                   @@@@0              ;0@@@@@00               @0@@0
                  0@@08              0@@@@@@@@@00              @@@00
                10@@@@              8@@@@@@@@@@@0              i0@@@0
               C@@@@@@t             G@@@@@@@@@@0@              0@@@@@@,
              0@@@@@@0@L             0@@@@@@@@0@8             0@@@@@@@@L
             00@@@@@@@@@0t             0@@00@00             00@@@@@@@@@08
            0@@@@@@@@@@@@@@@G                            00@@@@@@@@@@@@@@0
           0@@@@@@@@@@@@@@@@@@@@00i                t0@@@@@@@@@@@@@@@@@@@@@0
         :0@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@00@@@@@@@@@@@@@@@@@@@@@@@@000
        G@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@0
       00@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@01
      0@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@8
     0@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@00
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@00
```