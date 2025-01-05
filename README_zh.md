# YunDoBridge

YunDoBridge 是一个连接嵌入式设备与 AI 语音服务的桥接服务。它使用 Cloudflare Calls 作为中间层，帮助像 ESP32 这样的嵌入式设备轻松接入 OpenAI 的实时语音服务和 Google Gemini 的 Live API。

[English Documentation](README.md)

### 特性

- 支持 ESP32 等嵌入式设备
- 使用 Cloudflare Calls 作为中间层，解决网络连接问题
- 支持多种 AI 服务：
  - OpenAI 实时语音服务
  - Google Gemini Live API
- 专门优化支持 YunDo 开源硬件

### 架构

```
用户设备 (ESP32/YunDo) <-> CF Calls <-> YunDoBridge <-> AI 服务 (OpenAI/Gemini)
```

### 环境要求

- Go 1.16 或更高版本
- 配置文件 (`config.json`)，包含：
  ```json
  {
      "openai_api_key": "你的OpenAI密钥",
      "openai_model_endpoint": "OpenAI模型端点",
      "calls_base_url": "CF_Calls服务地址",
      "calls_app_id": "CF_Calls应用ID",
      "calls_app_token": "CF_Calls令牌",
      "port": "8080"
  }
  ```

### 快速开始

1. 克隆仓库：
```bash
git clone https://github.com/hx23840/YunDoBridge.git
cd YunDoBridge
```

2. 创建并配置 `config.json` 文件：
```bash
cp config.json.example config.json
```
然后编辑 `config.json` 文件，填入你的配置：
```json
{
    "openai_api_key": "你的OpenAI密钥",
    "openai_model_endpoint": "OpenAI模型端点",
    "calls_base_url": "CF_Calls服务地址",
    "calls_app_id": "CF_Calls应用ID",
    "calls_app_token": "CF_Calls令牌",
    "port": "8080"
}
```

3. 构建和运行：
```bash
go build -o YunDoBridge cmd/server/main.go
./YunDoBridge -config config.json
```

### 硬件支持

YunDo 是一个基于 ESP32-S3 的开源硬件平台，专门为 AI 应用设计。硬件详情和组装说明请访问：[YunDo 硬件项目](https://github.com/hx23840/YunDo)

主要特点：
- ESP32-S3 微控制器
- 优化的音频电路设计
- 内置麦克风和扬声器
- 低功耗设计
- 完整的参考实现

### API 接口

#### WebRTC 端点

```
POST /endpoint
```

用于建立 WebRTC 连接的端点，支持：
- 自动 SDP 协商
- 双向音频流传输
- 实时 AI 响应

### 许可证

本项目采用 [GNU General Public License v3.0](LICENSE) 许可证。

### 相关项目

- [YunDo 硬件](https://github.com/hx23840/YunDo) - 基于 ESP32-S3 的开源 AI 对话系统
- ESP32 WebRTC Client 实现

### 问题反馈

如果您在使用过程中遇到任何问题，欢迎通过 Issues 反馈。
