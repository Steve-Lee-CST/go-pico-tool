# go-pico-tool

> 一些简单易用的 Golang 组件

## 基本结构

对于大部分组件，遵循以下结构
- config.go: 配置文件
    - 一般情况下, 会包含 **配置类** 和一个 **默认配置**
- helper.go: 辅助函数（如果文件存在）
- tool.go: 工具本体
- tool_test.go: 测试文件
- _readme.md: 使用说明

## 工具列表

- ID Generator: [id generator 使用说明](./pkg/id_generator/_readme.cn.md)
- Request ID Tool: [request id tool 使用说明](./pkg/gin_pkg/request_id/_readme.cn.md)
- Http Decoder: [http decoder 使用说明](./pkg/gin_pkg/http_decoder/_readme.cn.md)
- Task DAG Flow: [task dag flow 使用说明](./pkg/task_dagflow/_readme.cn.md)
