# Caia - 基于Actor模型的Go语言高并发TCP框架

[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Benchmark](https://img.shields.io/badge/Benchmark-20k%20conn%2FGB-important?style=flat&logo=speedtest)]()

Caia是专为海量TCP长连接场景设计的轻量级并发框架，采用Golang原生channel实现Actor模型核心特性。

**核心特性** 🚀
- 🌳 树形协程监督体系 - 父子协程自动级联管理
- 📨 智能消息路由 - 支持同步/异步双模消息投递
- 🛡️ 生产级容错 - 内置panic恢复与连接隔离机制
- 🔧 轻量化设计 - 核心模块<1000LOC 
  - 🖥️ 资源效率：20,000+连接/GB内存 
  - ⚡ 处理能力：单节点QPS 10k+（2核2G/1KB消息）
  - 📦 极简依赖：仅需 `golang.org/x/sys` 和 `github.com/panjf2000/ants/v2`

**典型场景** 💡
- 🎮 实时游戏服务器（MMO/棋牌类）
- 📡 物联网设备接入网关
- 💬 千万级IM消息推送中台

**快速开始** ⚡

**设计优势** 🎯

<!-- ```mermaid

``` -->

通过层级化Actor设计，实现：
- 🔋 单连接故障隔离
- ♻️ 资源自动回收
- 🚦 流量分级管控