# openapi-cli

openapi-cli 是一个快速检查、升级swagger/openAPI文档的小工具。

## 使用方式：

```
USAGE:
openapi-cli [command] [command options] [arguments...]

COMMANDS:
lint, l     lint swagger/openapi document
upgrade, u  upgrade swagger2 to openapi3
version, v  Show version
help, h     Shows a list of commands or help for one command
```

### 基础命令

1. lint <filename>，检查openAPI中有问题的语法错误；
    ```
   $ ./openapi-cli-v0.2.3-darwin-amd64 lint testdata/swagger2-wrong.json
   2023-06-19T17:15:06.591+0800    info    pkg/lint.go:37  api lint        {"file": "testdata/swagger2-wrong.json"}
   2023-06-19T17:15:06.611+0800    info    pkg/lint.go:48  violation       {"result": {"valid":false,"path":"/policy/classifySample/getList","method":"post","startLine":26,"endLine":31,"description":"Operation parameters are unique and non-repeating.","howToFix":"Make sure that all the operation parameters are unique and non-repeating, don't duplicate names, don'tre-use parameter names in the same operation."}}
   2023-06-19T17:15:06.611+0800    info    pkg/lint.go:52  api lint finished       {"file": "testdata/swagger2-wrong.json"}
   ```
2. upgrade <filename>，升级swagger2到openAPI3（生成为一个时间戳后缀的json文件）。
   ```
   $ ./openapi-cli-v0.2.3-darwin-amd64 upgrade testdata/swagger2-wrong.json
   2023-06-19T17:18:33.333+0800    info    pkg/openapi2.go:65      api upgrade     {"file": "testdata/swagger2-wrong.json"}
   2023-06-19T17:18:33.347+0800    info    pkg/openapi2.go:144     delete invalid operation        {"operation": {"valid":false,"path":"/policy/classifySample/getList","method":"post","startLine":1,"endLine":1,"description":"Operation parameters are unique and non-repeating.","howToFix":"Make sure that all the operation parameters are unique and non-repeating, don't duplicate names, don'tre-use parameter names in the same operation."}}
   2023-06-19T17:18:33.347+0800    info    pkg/openapi2.go:96      api upgrade successfully        {"file": "testdata/swagger2-wrong-1687166313347.json", "version": "3.0.3", "duration": "14.235083ms"}
   ```

ref:
https://github.com/urfave/cli
https://github.com/LucyBot-Inc/api-spec-converter
https://github.com/getkin/kin-openapi
