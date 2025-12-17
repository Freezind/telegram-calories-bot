# LLM 测试使用指南

## 📋 测试类型：半自动化交互测试

由于 Telegram Bot API 的限制，测试采用**半自动化**方式：
- ✅ **自动化部分**：LLM 评判、报告生成、提示词存档
- 🤝 **手动部分**：你需要在 Telegram 中发送消息和点击按钮

这种方式**不需要部署到云端**，本地开发即可测试！

---

## 🚀 快速开始

### 1. 启动本地服务

```bash
# 终端 1: 启动 bot
go run cmd/unified/main.go

# 终端 2: 启动 Mini App 前端
cd web && npm run dev

# 终端 3: 创建 HTTPS tunnel (Mini App 需要)
ngrok http 5173
# 或
cloudflared tunnel --url http://localhost:5173
```

### 2. 配置测试环境

创建 `.env.test` 文件：

```bash
# Bot token（和 .env 一样）
TELEGRAM_BOT_TOKEN=你的bot_token

# 你的 Telegram 用户 ID（发消息给 @userinfobot 获取）
TELEGRAM_TEST_CHAT_ID=123456789

# Tunnel 的 HTTPS URL
MINIAPP_URL=https://abc123.ngrok.io

# Gemini API key
GEMINI_API_KEY=你的gemini_key

# 测试图片
TEST_FOOD_IMAGE_PATH=tests/fixtures/food0.jpg
```

### 3. 运行测试

```bash
./test-llm.sh
```

---

## 📱 测试流程（5 个场景）

### Scenario 1: /start 命令

```
📱 MANUAL ACTION REQUIRED:
   Please open Telegram and send '/start' to your bot now.
   Waiting for bot response...
   Press ENTER after you've sent the message:
```

**你需要做的：**
1. 在 Telegram 中给你的 bot 发送 `/start`
2. 等待 bot 回复
3. 按 ENTER 继续

**测试会验证：** Bot 回复了欢迎消息和使用说明

---

### Scenario 2: /estimate + 图片上传

```
📱 MANUAL ACTION REQUIRED:
   1. Open Telegram and send '/estimate' to your bot
   2. When bot asks for image, send a food photo
   3. You can use the test image: tests/fixtures/food0.jpg
   Press ENTER after bot responds with estimate:
```

**你需要做的：**
1. 给 bot 发送 `/estimate`
2. 当 bot 要求图片时，发送一张食物照片
3. 等待 bot 返回估算结果
4. 按 ENTER 继续

**测试会验证：** Bot 返回了包含食物列表、卡路里和置信度的结构化回复

---

### Scenario 3: Re-estimate 按钮保留测试

```
📱 MANUAL ACTION REQUIRED:
   1. Look at your previous estimate message in Telegram
   2. Click the 'Re-estimate' button
   3. Observe that:
      - The previous estimate message is NOT deleted
      - Bot sends a NEW message asking for another image
   Press ENTER after clicking Re-estimate:
```

**你需要做的：**
1. 在 Telegram 中找到刚才的估算消息
2. 点击 "Re-estimate" 按钮
3. **观察**：之前的估算消息是否还在（没被删除）
4. 按 ENTER 继续

**测试会验证：** Bot 发送了新的提示消息（说明按钮有效）

---

### Scenario 4: Cancel 按钮保留测试

```
📱 MANUAL ACTION REQUIRED:
   1. Send '/estimate' to the bot again
   2. Send a food image
   3. When bot shows the estimate, click the 'Cancel' button
   4. Observe that:
      - The estimate message is NOT deleted
      - Bot sends a cancellation confirmation message
   Press ENTER after clicking Cancel:
```

**你需要做的：**
1. 再次发送 `/estimate` 和食物图片
2. 等待估算结果
3. 点击 "Cancel" 按钮
4. **观察**：估算消息是否还在（没被删除）
5. 按 ENTER 继续

**测试会验证：** Bot 发送了取消确认消息

---

### Scenario 5: Mini App 页面加载

这个场景是**全自动**的，使用 Playwright 访问 Mini App URL。

**测试会验证：** Mini App 能正常加载，页面包含 "Calorie Log"、"Add New Log" 或 "No logs yet" 等预期文本

---

## 📊 测试结果

测试完成后，你会得到：

### 1. 控制台输出

```
[S1] Testing /start command...
  ✅ PASS - Bot responds with welcome message...

[S2] Testing /estimate + image upload...
  ✅ PASS - Response contains foods list, calories...

...

Test Summary
========================================
Total scenarios: 5
Passed: 5
Failed: 0
Duration: 2m15s

Result: PASS (all scenarios passed)
```

### 2. 详细报告

生成在 `reports/004-test-report.md`，包含：
- 每个场景的执行步骤
- 捕获的证据（JSON 格式）
- LLM 评判结果和理由
- 时间戳

### 3. LLM 提示词存档

所有 LLM 评判提示词会追加到 `prompts.md`

---

## ⚠️ 重要说明

### Bot API 限制

Bot API **只能**让 bot 发送消息，**不能**模拟用户发送消息或点击按钮。因此：

❌ **不可能做到：**
- 自动让用户给 bot 发消息
- 自动点击 inline 按钮

✅ **我们的解决方案：**
- 用户手动发送消息/点击按钮
- 测试自动捕获 bot 响应
- LLM 自动评判响应是否正确

这种**半自动化**方式在保持简单的同时，仍然提供了有价值的自动化验证。

### 本地测试 vs 云端部署

**本地测试**（推荐）：
- Bot: 本地运行，使用 long polling（不需要 webhook，不需要 tunnel）
- Mini App: 本地运行 + ngrok/cloudflare tunnel 提供 HTTPS URL

**云端测试**：
- 部署到 Railway/Render 后也可以测试
- 配置一样，只是 MINIAPP_URL 换成云端地址

---

## 🔧 故障排除

### 问题：找不到 bot 响应

**解决：**
- 确认 bot 正在运行（`go run cmd/unified/main.go`）
- 确认你确实给 bot 发送了消息
- 等待几秒让 bot 处理

### 问题：Mini App 测试失败

**解决：**
- 确认 tunnel 正在运行
- 确认 tunnel URL 是 HTTPS 开头
- 在浏览器中打开 URL 验证能访问
- 确认前端正在运行（`npm run dev`）

### 问题：Playwright 相关错误

**解决：**
- 首次运行会自动安装 Chromium 浏览器
- 如果失败，手动运行：
  ```bash
  go run github.com/playwright-community/playwright-go/cmd/playwright install chromium
  ```

---

## 💡 提示

1. **保持 Telegram 打开**：测试期间保持 Telegram 应用打开，方便快速发送消息

2. **准备好食物图片**：可以直接在 Telegram 中发送 `tests/fixtures/food0.jpg` 或任何食物照片

3. **按自己的节奏**：测试会等你按 ENTER，不用着急，可以仔细检查每一步

4. **查看详细日志**：如果测试失败，查看 `reports/004-test-report.md` 了解详情

---

## ✅ 成功标准

测试通过需要满足：

- S1: Bot 回复了欢迎消息
- S2: Bot 返回了包含食物、卡路里、置信度的估算
- S3: 点击 Re-estimate 后，Bot 发送了新提示，之前的消息没被删除
- S4: 点击 Cancel 后，Bot 发送了确认消息，估算消息没被删除
- S5: Mini App 页面加载成功并显示预期 UI 元素

---

**祝测试顺利！** 🎉
