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
   $ ./openapi-cli-0.2.0-darwin-amd64 lint ./testdata/swagger2-wrong.json
   2023-05-29T11:14:51.950+0800    info    pkg/lint.go:47  violation       {"result": {"valid":false,"path":"/policy/classifySample/getList","method":"post","startLine":828,"endLine":833,"description":"Operation parameters are unique and non-repeating.","howToFix":"Make sure that all the operation parameters are unique and non-repeating, don't duplicate names, don'tre-use parameter names in the same operation."}}
   ```
2. upgrade <filename>，升级swagger2到openAPI3（生成为一个时间戳后缀的json文件）。
   ```
   $ ./openapi-cli-0.2.0-darwin-amd64 upgrade ./testdata/swagger2-wrong.json 
   2023-05-29T11:32:25.231+0800    info    pkg/openapi2.go:70      api upgrade     {"file": "./testdata/swagger2-wrong.json"}
   2023-05-29T11:32:25.353+0800    info    pkg/openapi2.go:148     invalid operation has been deleted      {"operation": {"valid":false,"path":"/policy/classifySample/getList","method":"post","startLine":1,"endLine":1,"description":"Operation parameters are unique and non-repeating.","howToFix":"Make sure that all the operation parameters are unique and non-repeating, don't duplicate names, don'tre-use parameter names in the same operation."}}
   2023-05-29T11:32:25.364+0800    info    pkg/openapi2.go:100     api upgrade successfully        {"upgraded version": "3.0.3", "duration": "131.405333ms"}
   $ ll testdata
   -rw-r--r--@ 1 eric  staff   280K May 29 10:57 swagger2-wrong.json
   -rwxr-xr-x@ 1 eric  staff   340K May 29 11:33 swagger2-wrong363774000.json
   ```

ref:
https://github.com/urfave/cli
https://github.com/LucyBot-Inc/api-spec-converter
https://github.com/getkin/kin-openapi
