# YunDoBridge

YunDoBridge is a bridging service designed to connect embedded devices with AI voice services. It utilizes Cloudflare Calls as a middleware to help embedded devices like ESP32 easily access OpenAI's real-time voice service and Google Gemini's Live API.

[中文文档](README_zh.md)

### Features

- Embedded device support (ESP32)
- Network optimization via Cloudflare Calls
- Multiple AI service integration:
  - OpenAI real-time voice service
  - Google Gemini Live API
- Optimized for YunDo open-source hardware

### Architecture

```
User Device (ESP32/YunDo) <-> CF Calls <-> YunDoBridge <-> AI Services (OpenAI/Gemini)
```

### Requirements

- Go 1.16+
- Environment variables:
  ```
  OPENAI_API_KEY=your_openai_key
  OPENAI_MODEL_ENDPOINT=your_openai_endpoint
  CALLS_BASE_URL=your_cf_calls_url
  CALLS_APP_ID=your_cf_app_id
  CALLS_APP_TOKEN=your_cf_app_token
  ```
- Configuration file (`config.json`) with:
  ```json
  {
      "openai_api_key": "your_openai_key",
      "openai_model_endpoint": "your_openai_endpoint",
      "calls_base_url": "your_cf_calls_url",
      "calls_app_id": "your_cf_app_id",
      "calls_app_token": "your_cf_app_token",
      "port": "8080"
  }
  ```

### Quick Start

1. Clone the repository:
```bash
git clone https://github.com/hx23840/YunDoBridge.git
cd YunDoBridge
```

2. Create and configure your `config.json` file:
```bash
cp config.json.example config.json
```
Then edit `config.json` with your settings:
```json
{
    "openai_api_key": "your_openai_key",
    "openai_model_endpoint": "your_openai_endpoint",
    "calls_base_url": "your_cf_calls_url",
    "calls_app_id": "your_cf_app_id",
    "calls_app_token": "your_cf_app_token",
    "port": "8080"
}
```

3. Build and run:
```bash
go build -o YunDoBridge cmd/server/main.go
./YunDoBridge -config config.json
```

### Hardware Support

YunDo is an ESP32-S3-based open-source hardware platform specifically designed for AI applications. For hardware details and assembly instructions, please visit: [YunDo Hardware Project](https://github.com/hx23840/YunDo)

Main features:
- ESP32-S3 microcontroller
- Optimized audio circuits
- Built-in microphone and speaker
- Low power consumption design
- Complete reference implementation

### API Reference

#### WebRTC Endpoint

```
POST /endpoint
```

This endpoint is used for establishing WebRTC connections, supporting:
- Automatic SDP negotiation
- Bidirectional audio streaming
- Real-time AI responses

### License

This project is licensed under the [GNU General Public License v3.0](LICENSE).

### Related Projects

- [YunDo Hardware](https://github.com/hx23840/YunDo) - Open source ESP32-based AI dialogue system
- ESP32 WebRTC Client Implementation

### Issues

If you encounter any problems, please feel free to create an issue.