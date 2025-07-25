# go-pico-tool

> Some easy-to-use go-dev components.

For the Chinese version, see [中文说明](./README.cn.md)

## Structure Overview

For most components, the following structure is used:
- config.go: Configuration file
    - Usually contains a **Config struct** and a **default config**
- helper.go: Helper functions (if present)
- tool.go: Main tool implementation
- tool_test.go: Test file
- _readme.md: Usage documentation

## Components List

- ID Generator: [id generator usage](./pkg/id_generator/_readme.en.md)
- Request ID Tool: [request id tool usage](./pkg/gin_pkg/request_id/_readme.en.md)
- Http Decoder: [http decoder usage](./pkg/gin_pkg/http_decoder/_readme.en.md)
- Task DAG Flow: [task dag flow usage](./pkg/task_dagflow/_readme.en.md)