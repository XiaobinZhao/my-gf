go web framework base on goFame.

- https://goframe.org/

# 准备工作

## go语言安装

window版本的安装，自行下载exe进行安装，然后设置环境变量即可。以下以linux安装为主。本次安装版本使用：`go1.16.13`

### apt-get安装

`apt-get install golang`

Debian9 安装的golang版本默认为`golang-1.7`

### 手动安装

1. 到官网`https://golang.google.cn/dl/`下载安装包`wget https://golang.google.cn/dl/go1.16.13.linux-amd64.tar.gz`

2. 创建目录`mkdir -p /usr/local/lib`

3. 解压：`tar -xzf go1.16.13.linux-amd64.tar.gz -C /usr/local/lib`

4. 设置环境变量 编辑 `~/.profile`, 增加以下内容

   ```shell
   # golang
   GOPATH=/root/go-workspace/
   GOROOT=/usr/local/lib/go
   export GO111MODULE=on
   export GOPROXY=https://goproxy.io
   export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
   ```

5. 刷新环境变量。 `source ~/.profile`

6. 安装完成，测试。`go version`

## 为什么选择GoFrame

[Golang框架选型比较: goframe, beego, iris和gin - GoFrame社区投稿 - GoFrame官网 - 类似PHP-Laravel, Java-SpringBoot的Go企业级开发框架](https://goframe.org/pages/viewpage.action?pageId=3673375)

总结：

1. 团队需要一个统一的技术框架，而不是东拼西凑的一堆单仓库包。
2. 我们只需要维护一个框架的版本，而不是维护数十个单仓库包的版本。
3. 框架的组件必须模块化、低耦合设计，保证内部组件也可以独立引用。
4. 核心组件严格禁止单仓库包设计，并且必须由框架统一维护。

走过这么多弯路之后，我们决心建立一套成体系的`Golang`开发框架。除了要求团队能够快速学习，维护成本低，并且我们最主要的诉求，是核心组件不能是半成品，框架必须是上过大规模生产验证的，稳定和成熟的。随着我们重新对行业中流行的技术框架做了技术评估，包括上面说的那些框架。原本的初衷是想将内部的各个轮子统一做一个成体系的框架，在开源项目中找一些有价值的参考。

后来找到了`goframe`，仔细评估和学习了框架设计，发现框架设计思想和我们的经验总结如出一辙！

这里不得不提一件尴尬的事情。其实最开始转`Golang`之前（2019年中旬）也做过一些调研，那时`goframe`版本还不高，并且我们负责评估的团队成员有一种先入为主的思想，看到模块和文档这么多，感觉应该挺复杂，性能应该不高，于是没怎么看就PASS。后来选择一些看起来简单的开源轮子自己做了些二次封装。

这次经过一段时间的仔细调研和源码学习，得出一个结论，`goframe`框架的框架架构、模块化和工程化设计思想非常棒，执行效率很高，模块不仅丰富，而且质量之高，令人惊叹至极！相比较我们之前写的那些半成品轮子，简直就是小巫见大巫。**团队踩了一年多的坑，才发现团队确实需要一个统一的技术框架而不是一堆不成体系的轮子，其实人家早已给了一条明光大道，并且一直在前面默默努力。**

经过团队内部的调研和讨论，我们决定使用`goframe`逐步重构我们的业务项目。由于`goframe`是模块化设计的，因此我们也可以对一些模块做必要的替换。重构过程比较顺利，基础技术框架的重构并不会对业务逻辑造成什么影响，反而通过`goframe`的工程化思想和很棒的开发工具链，在统一技术框架后，极大地提高了项目的开发和维护效率，使得团队可以专心于业务开发，部门也陆续有了更多的产出。目前我们已经有大部门业务项目转向了`goframe`，平台每日流量千万级别。 

## Go Frame 目录结构

`GoFrame`业务项目基本目录结构如下：

```shell
/
├── api
├── internal
│   ├── cmd
│   ├── consts
│   ├── controller
│   ├── model
│   │   └── entity
│   └── service
│       └── internal
│           ├── dao
│           └── do
├── manifest
├── resource
├── utility
├── go.mod
└── main.go 
```

工程目录采用了通用化的设计，实际项目中可以根据项目需要适当增减模板给定的目录。例如，没有`i18n`及`template`需求的场景，直接删除对应目录即可。

| 目录/文件名称   | 说明     | 描述                                                         |
| :-------------- | :------- | :----------------------------------------------------------- |
| `api`           | 接口定义 | 对外提供服务的输入/输出数据结构定义。考虑到版本管理需要，往往以`apiv1/apiv2...`存在。 |
| `internal`      | 内部逻辑 | 业务逻辑存放目录。通过`Golang internal`特性对外部隐藏可见性（导入路径包含`internal`关键字的包，只允许`internal`的父级目录及父级目录的子包导入，其它包无法导入）。 |
| ` - cmd`        | 入口指令 | 命令行管理目录。可以管理维护多个命令行。                     |
| ` - consts`     | 常量定义 | 项目所有常量定义。                                           |
| ` - controller` | 接口处理 | 接收/解析用户输入参数的入口/接口层。                         |
| ` - model`      | 结构模型 | 数据结构管理模块，管理数据实体对象，以及输入与输出数据结构定义。 |
| `  - entity`    | 数据模型 | 数据模型是模型与数据集合的一对一关系，由工具维护，用户不能修改。 |
| ` - service`    | 逻辑封装 | 业务逻辑封装管理，特定的业务逻辑实现和封装。                 |
| `  - dao`       | 数据访问 | 数据访问对象，这是一层抽象对象，用于和底层数据库交互，仅包含最基础的 `CURD` 方法 |
| `  - do`        | 领域对象 | 用于`dao`数据操作中业务模型与实例模型转换，由工具维护，用户不能修改。 |
| `manifest`      | 交付清单 | 包含程序编译、部署、运行、配置的文件。常见内容如下：         |
| ` - config`     | 配置管理 | 配置文件存放目录。                                           |
| ` - docker`     | 镜像文件 | `Docker`镜像相关依赖文件，脚本文件等等。                     |
| ` - deploy`     | 部署文件 | 部署相关的文件。默认提供了`Kubernetes`集群化部署的`Yaml`模板，通过`kustomize`管理。 |
| `resource`      | 静态资源 | 静态资源文件。这些文件往往可以通过 资源打包/镜像编译 的形式注入到发布文件中。 |
| `go.mod`        | 依赖管理 | 使用`Go Module`包管理的依赖描述文件。                        |
| `main.go`       | 入口文件 | 程序入口文件。                                               |

### 业务接口 - `api`

业务接口包含两部分：接口定义（`api`）+接口实现（`controller`）。

`api`包的职责类似于三层架构设计中的`UI`表示层，负责接收并响应客户端的输入与输出，包括对输入参数的过滤、转换、校验，对输出数据结构的维护，并调用 `service` 实现业务逻辑处理。

### 逻辑封装 - `service`

`service`包的职责类似于三层架构设计中的`BLL`业务逻辑层，负责具体业务逻辑的实现以及封装。

### 数据访问 - `dao`

`dao`包的职责类似于三层架构中的`DAL`数据访问层，数据访问层负责所有的数据访问收口。

### 结构模型 - `model`

`model`包的职责类似于三层架构中的`Model`模型定义层。模型定义代码层中仅包含全局公开的数据结构定义，往往不包含方法定义。

这里需要注意的是，这里的`model`不仅负责维护数据实体对象（`entity`）结构定义，也包括所有的输入/输出数据结构定义，被`api/dao/service`共同引用。这样做的好处除了可以统一管理公开的数据结构定义，也可以充分对同一业务领域的数据结构进行复用，减少代码冗余。

![img](https://gfcdn.johng.cn/download/attachments/30740166/image2022-1-18_0-47-31.png?version=1&modificationDate=1642437918159&api=v2)

三层架构设计与框架代码分层映射关系

## 请求分层流转

![img](https://gfcdn.johng.cn/download/attachments/30740166/image2022-1-18_10-38-49.png?version=1&modificationDate=1642473395894&api=v2)

### cmd

`cmd`层负责引导程序启动，显著的工作是初始化逻辑、注册路由对象、启动`server`监听、阻塞运行程序直至`server`退出。

### api

上层`server`服务接收客户端请求，转换为`api`中定义的`Req`接收对象、执行请求参数到`Req`对象属性的类型转换、执行`Req`对象中绑定的基础校验并转交`Req`请求对象给`controller`层。

### controller

`controller`层负责接收`Req`请求对象后做一些业务逻辑校验，随后调用一个或多个`service`实现业务逻辑，将执行结构封装为约定的`Res`数据结构对象返回。

### model

`model`层中管理了所有的业务模型，`service`资源的`Input/Output`输入输出数据结构都由`model`层来维护。

### service

`service`层的业务逻辑需要通过调用`dao`来实现数据的操作，调用`dao`时需要传递`do`数据结构对象，用于传递查询条件、输入数据。`dao`执行完毕后通过`Entity`数据模型将数据结果返回给`service`层。

### dao

`dao`层通过框架的`ORM`抽象层组件与底层真实的数据库交互。

  

## Go Frame常见问题解答

### 1. 框架是否支持常见的`MVC`开发模式

**当然！**

作为一款模块化设计的基础开发框架，`GoFrame`不会局限代码设计模式，并且框架提供了非常强大的模板引擎核心组件，可快速用于`MVC`模式中常见的模板渲染开发。相比较`MVC`开发模式，在复杂业务场景中，我们更推荐使大家用三层架构设计模式。

### 2. 如何清晰界定和管理`service`和`controller`的分层职责

`controller`层处理`Req/Res`外部接口请求。负责接收、校验请求参数，并调用**一个或多个** `service`来实现业务逻辑处理，根据返回数据结构组装数据再返回。

`service`层处理`Input/Output`内部方法调用。负责内部**可复用**的业务逻辑封装，封装的方法粒度往往比较细。

因此， **禁止** 从`controller`层直接透传`Req`对象给`service`，也禁止`service`直接返回`Res`数据结构对象，因为`service`服务的主体与`controller`完全不同。当您错误地使用`service`方法处理特定的`Req`对象的时候，该方法也就与对于的外部接口耦合，仅为外部接口服务，难以复用。这种场景下`service`替代了`controller`的作用，造成了本末倒置。

### 3. 如何清晰界定和管理`service`和`dao`的分层职责

这是一个很经典的问题。

**痛点：**

常见的，开发者把数据相关的业务逻辑实现封装到了`dao`代码层中，而`service`代码层只是简单的`dao`调用，这么做的话会使得原本负责维护数据的`dao`层代码越来越繁重，反而业务逻辑`service`层代码显得比较轻。开发者存在困惑，我写的业务逻辑代码到底应该放到`dao`还是`service`中？

业务逻辑其实绝大部分时候都是对数据的`CURD`处理，这样做会使得几乎所有的业务逻辑会逐步沉淀在`dao`层中，业务逻辑的改变其实会频繁对`dao`层的代码产生修改。例如：数据查询在初期的时候可能只是简单的逻辑，目前代码放到`dao`好像也没问题，但是查询需求增加或变化变得复杂之后，那么必定会继续维护修改原有的`dao`代码，同时`service`代码也可能同时做更新。原本仅限于`service`层的业务逻辑代码职责与`dao`层代码职责模糊不清、耦合较重，原本只需要修改`service`代码的需求变成了同时修改`service`+`dao`，使得项目中后期的开发维护成本大大增加。

**建议：**

我们的建议:`dao`层的代码应该尽量保证通用性，并且大部分场景下不需要增加额外方法，只需要使用一些通用的链式操作方法拼凑即可满足。业务逻辑、包括看似只是简单的数据操作的逻辑都应当封装到`service`中，`service`中包含多个业务模块，每个模块独自管理自己的`dao`对象，`service`与`service`之间通过相互调用方法来实现数据通信而不是随意去调用其他`service`模块的`dao`对象。

# 开始开发

## gf-cli安装

1. 二进制安装

   从[Releases · gogf/gf-cli (github.com)](https://github.com/gogf/gf-cli/releases)下载二进制，然后修改名字为`gf.exe`；放置到`%GOROOT%/bin`或者`%GOPATH%/bin`

2. 手动安装 [gf-cli/README.MD](https://github.com/gogf/gf-cli/blob/master/README.MD)

   `go install github.com/gogf/gf-cli/v2/gf@master`

3. 校验

   `gf -v`

   得到如下输出证明OK：

   ```shell
   GoFrame CLI Tool v2.0.0-rc, https://goframe.org
   GoFrame Version: cannot find go.mod
   CLI Installed At: E:\go-workspace\bin\gf.exe
   CLI Built Detail:
     Go Version:  go1.17.6
     GF Version:  v2.0.0-beta
     Git Commit:  2022-01-24 11:02:57 37d535f61e975dff1765f98e3b505aeccfec338e
     Build Time:  2022-01-24 03:01:42
   ```

## 数据库表设计

在数据库中进行表结构设计，包括字段、长度、主键、描述等。

## gf cli自动生成dao/do/entity

1. 配置好配置文件，重点

   ```toml
   # GF-CLI工具配置
   [gfcli]
       # 自定义DAO生成配置(默认是读取database配置)
       [[gfcli.gen.dao]]
           link   = "mysql:root:mysql@tcp(192.168.212.117:3306)/myapp"
           descriptionTag =   true
           noModelComment =   true
   ```

2.  使用gf  cli命令

   `gf gen dao`，成功时可以得到如下类似的输出：

   ```shell
   > gf gen dao
   generated: internal\service/internal/dao\desktop.go
   generated: internal\service/internal/dao\internal\desktop.go
   generated: internal\service/internal/dao\user.go
   generated: internal\service/internal/dao\internal\user.go
   generated: internal\service/internal/do\desktop.go
   generated: internal\service/internal/do\user.go
   generated: internal\model/entity\desktop.go
   generated: internal\model/entity\user.go
   done!
   ```

   可以查看具体文件是否生成。

3. 编写main.py 和 路由

