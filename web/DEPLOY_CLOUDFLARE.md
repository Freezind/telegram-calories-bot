# Cloudflare Pages 部署指南

## 前提条件

- 后端已部署到 Railway: `https://telegram-calories-bot-production.up.railway.app`
- 已安装 Node.js 和 npm

## 步骤 1: 构建前端

```bash
cd web
npm install
npm run build
```

构建产物在 `dist/` 目录。

## 步骤 2: 部署到 Cloudflare Pages

### 方法 1: 通过 Cloudflare Dashboard（推荐）

1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. 进入 **Pages** → **Create a project**
3. 选择 **Connect to Git** 或 **Direct Upload**

#### 使用 Git 连接（推荐）:
- 连接你的 GitHub 仓库
- 选择 `telegram-calories-bot` 仓库
- 配置构建设置：
  - **Build command**: `cd web && npm install && npm run build`
  - **Build output directory**: `web/dist`
  - **Root directory**: `/` (留空或填 `/`)

- 环境变量：
  ```
  VITE_API_BASE_URL=https://telegram-calories-bot-production.up.railway.app
  ```

- 点击 **Save and Deploy**

#### 使用 Direct Upload:
```bash
cd web
npm run build
# 然后在 Cloudflare Dashboard 上传 dist/ 目录
```

### 方法 2: 使用 Wrangler CLI

```bash
# 安装 Wrangler
npm install -g wrangler

# 登录 Cloudflare
wrangler login

# 部署
cd web
npm run build
wrangler pages deploy dist --project-name=telegram-calories-miniapp
```

## 步骤 3: 配置 Railway 后端 CORS

部署完成后，你会得到一个 Cloudflare Pages URL，例如：
```
https://telegram-calories-miniapp.pages.dev
```

在 Railway Dashboard 中设置环境变量：

1. 进入你的 Railway 项目
2. 点击 **Variables** 标签
3. 添加新变量：
   ```
   MINIAPP_URL=https://telegram-calories-miniapp.pages.dev
   ```
4. 点击 **Save**
5. 重新部署

## 步骤 4: 在 BotFather 设置 Mini App URL

1. 打开 Telegram，找到 **@BotFather**
2. 发送 `/mybots`
3. 选择你的 bot
4. 选择 **Bot Settings**
5. 选择 **Menu Button**
6. 选择 **Configure Menu Button**
7. 输入你的 Cloudflare Pages URL:
   ```
   https://telegram-calories-miniapp.pages.dev
   ```
8. 输入按钮文本（例如: "📊 My Logs"）

## 步骤 5: 测试

1. 在 Telegram 中打开你的 bot
2. 点击底部的 Menu Button
3. Mini App 应该能够正常加载并与 Railway 后端通信

## 故障排查

### CORS 错误
如果看到 CORS 错误，确认：
- Railway 环境变量 `MINIAPP_URL` 设置正确
- URL 完全匹配（包括 `https://` 和域名）
- 已重新部署 Railway

### API 连接失败
如果无法连接后端：
- 检查 `web/.env.production` 中的 `VITE_API_BASE_URL`
- 确认 Railway 后端正在运行
- 检查 Railway 后端日志

### Mini App 无法加载
- 确认在 BotFather 中正确设置了 Menu Button URL
- 确认 URL 是完整的 HTTPS URL
- 尝试在浏览器中直接打开 Mini App URL

## 自动部署

Cloudflare Pages 可以配置自动部署：
- 每次推送到 `master` 分支自动触发构建
- Preview deployments for pull requests
- 环境变量在 Cloudflare Dashboard 中管理

## 更新部署

### 更新前端
```bash
git add .
git commit -m "update frontend"
git push
# Cloudflare Pages 会自动构建和部署
```

### 更新后端 API 地址
如果 Railway URL 改变了：
1. 更新 `web/.env.production`
2. 在 Cloudflare Pages Dashboard 更新环境变量
3. 触发重新构建
