# sbp

Simple Beautiful Package 高频使用的 Go 代码工具包合集

`sbp` is a universal utility collection for Go, it complements offerings such as Boost, Better std, Cloud tools.

## Table of Contents

- [Introduction](#Introduction)
- [Catalogs](#Catalogs)
- [Releases](#Releases)
- [How To Use](#How-To-Use)
- [License](#License)

## Introduction

`sbp` is a universal utility collection for Go, it complements offerings such as Boost, Better std, Cloud tools. It is migrated from the [![Go Reference](https://pkg.go.dev/badge/github.com/bytedance/gopkg.svg)](https://pkg.go.dev/github.com/bytedance/gopkg) `gopkg` 主要对其相关实战项目使用体验提升, 最初用 gopkg 感觉是里面代码水平层次不齐, 实战用起来很别扭, 逐渐想对其设计理念继续雕琢雕琢删繁就简.

We depend on the same code(this repo) in our production environment.

## Catalogs

* [cache](https://github.com/wangzhione/sbp/tree/master/cache): Caching Mechanism
* [cloud](https://github.com/wangzhione/sbp/tree/master/cloud): Cloud Computing Design Patterns
* [collection](https://github.com/wangzhione/sbp/tree/master/collection): Data Structures
* [lang](https://github.com/wangzhione/sbp/tree/master/lang): Enhanced Standard Libraries
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
