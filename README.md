# Caia - 基于Actor模型的Go语言高并发TCP框架

[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Caia是专为海量TCP长连接场景设计的轻量级并发框架，采用Golang原生channel实现Actor模型核心特性。

**核心特性** 🚀
- 🌳 树形协程监督体系 - 父子协程自动级联管理
- 📨 智能消息路由 - 支持同步/异步双模消息投递
- 🛡️ 生产级容错 - 内置panic恢复与连接隔离机制
- ⚡ 高性能通信 - 单协程消息处理+非阻塞IO调度
- 🔧 轻量化设计 - 无第三方依赖
<!-- ，2000+连接/GB内存 -->

**典型场景** 💡
- 实时消息推送中间件
- 多玩家在线游戏服务器

**快速开始** ⚡
<!-- ```go
// 创建网络Actor
netActor, _ := root.CreateChild("net", 100, func(msg ctypes.Message) {
    // TCP监听逻辑
})

// 启动Actor协程
netActor.Start()

// 发送控制消息
netActor.SendMessage(ctrl.StartListening{Port: 8080})
``` -->

**设计优势** 🎯

<!-- ```mermaid

``` -->

通过层级化Actor设计，实现：
- 单连接故障隔离
- 资源自动回收
- 流量分级管控
- 集群化扩展支持