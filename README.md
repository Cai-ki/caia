# Caia - 基于Actor模型的Go语言高并发TCP服务器框架

[![Go Version](https://img.shields.io/badge/go-1.26+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Benchmark](https://img.shields.io/badge/Benchmark-30k%20conn%2FGB-important?style=flat&logo=speedtest)]()
[![Benchmark](https://img.shields.io/badge/Benchmark-50k%20QPS-important?style=flat&logo=speedtest)]()

Caia是专为海量TCP长连接场景设计的轻量级并发服务器框架，采用Golang原生channel实现Actor模型核心特性。

**核心特性** 🚀
- 🌳 树形协程监督体系 - 父子协程自动级联管理
- 📨 智能消息路由 - 支持同步/异步双模消息投递
- 🛡️ 生产级容错 - 内置panic恢复与连接隔离机制
- 🔧 轻量化设计 - 核心模块<1000LOC 
- 📊 性能标杆(i5-12450H)
  - 🖥️ 资源效率：30K+ 连接/GB内存（1KB消息）
  - ⚡ 处理能力：单节点QPS 50k+（单核/10K连接/1KB消息）

**典型场景** 💡
- 🎮 实时游戏服务器（MMO/棋牌类）
- 📡 物联网设备接入网关
- 💬 千万级IM消息推送中台

**快速开始** 🎯

**设计哲学** 🧠

<!-- ```mermaid

``` -->
- 🔋 物理核尊重：单核性能优先，拒绝无意义抽象

- 🩺 故障即特征：崩溃信息自动诊断，隔离取代全局恢复

- 🔍 透明式并发：开发者仅关注业务Actor逻辑