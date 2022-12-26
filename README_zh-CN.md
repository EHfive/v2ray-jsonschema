# V2Ray JSON Schema

利用Go类型反射系统构建 [V2Ray](https://github.com/v2fly/v2ray-core) v4, v5 配置格式的 JSON schema

[English](./README.md)

## 使用方法

#### 在配置中声明schema

在V2Ray配置文件中将 `$schema` 设为指向对应配置格式的 JSON schema 的 URL，一些 IDE 可以通过此 schema 为你提供编辑提示

VSCode 支持此格式, 见此[文档](https://code.visualstudio.com/docs/languages/json#_json-schemas-and-settings)

#### Schemas
- jsonv4 schema (GitHub): https://github.com/EHfive/v2ray-jsonschema/raw/main/v4-config.schema.json
- jsonv5 schema (GitHub): https://github.com/EHfive/v2ray-jsonschema/raw/main/v5-config.schema.json

- jsonv4 schema (jsDelivr CDN): https://cdn.jsdelivr.net/gh/EHfive/v2ray-jsonschema/v4-config.schema.json
- jsonv5 schema (jsDelivr CDN): https://cdn.jsdelivr.net/gh/EHfive/v2ray-jsonschema/v5-config.schema.json

```json
{
  "$schema": "https://github.com/EHfive/v2ray-jsonschema/raw/main/v5-config.schema.json",
  "inbounds": [
    {
      "protocol": "socks",
      "listen": "127.0.0.1",
      "port": 1080
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom"
    }
  ]
}
```

![autosuggestion-demo](./doc/code-suggestion-demo.gif)
