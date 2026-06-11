# 前端架构详细设计 (`web/default/`)

## 1. 概述
前端是一个现代化的 React 19 应用程序，使用 TypeScript、Rsbuild 和基于组件的架构。支持双主题（default/classic），嵌入式 SPA 部署。

## 2. 技术栈详解
| 类别 | 库/工具 | 版本/说明 |
|---|---|---|
| **包管理器** | Bun | 首选包管理器和脚本运行器 |
| **框架** | React 19 + TypeScript | 最新 React 特性支持 |
| **构建工具** | Rsbuild | 快速构建，支持现代前端特性 |
| **路由** | TanStack Router | 基于文件的路由系统 |
| **数据获取** | TanStack React Query + Axios | 缓存、失效、重新获取 |
| **状态管理** | Zustand | 轻量级，持久化到 localStorage |
| **UI 组件** | Base UI + Tailwind CSS v4 | 无样式组件库 + 原子化 CSS |
| **表单** | React Hook Form + Zod | 类型安全的表单处理 |
| **表格** | TanStack React Table | 高性能表格组件 |
| **图表** | VChart (@visactor) | 数据可视化 |
| **国际化** | i18next + react-i18next | 多语言支持（en, zh, fr, ru, ja, vi） |

## 3. 目录结构详解
```text
web/default/src/
├── routes/                    # 基于文件的路由 (TanStack Router)
│   ├── __root.tsx             # 根路由（设置检查、主题提供者）
│   ├── index.tsx              # 首页
│   ├── setup/                 # 初始设置向导
│   ├── (auth)/                # 认证路由（登录、注册、重置、OAuth）
│   └── _authenticated/        # 受保护路由布局
│       ├── dashboard/         # 仪表盘
│       ├── keys/              # API 密钥管理
│       ├── channels/          # 渠道管理（管理员）
│       ├── users/             # 用户管理（管理员）
│       ├── models/            # 模型管理（管理员）
│       ├── playground/        # 聊天 Playground
│       ├── wallet/            # 钱包与充值
│       ├── profile/           # 用户资料
│       ├── usage-logs/        # 使用日志
│       ├── subscriptions/     # 订阅管理
│       ├── redemption-codes/  # 兑换码管理
│       └── system-settings/   # 系统设置（管理员）
├── stores/                    # Zustand 状态管理
│   ├── auth-store.ts          # 用户认证状态（持久化到 localStorage）
│   ├── notification-store.ts  # 通知状态
│   └── system-config-store.ts # 系统配置状态
├── features/                  # 功能模块
│   ├── wallet/                # 钱包功能（充值、支付、计费、推荐）
│   ├── models/                # 模型管理功能
│   ├── chat/                  # 聊天功能
│   ├── about/                 # 关于页面
│   └── subscriptions/         # 订阅功能
├── components/                # 共享 UI 组件
├── lib/                       # 共享工具函数
├── hooks/                     # 自定义 React hooks
├── i18n/                      # 前端国际化 (i18next)
│   └── locales/               # 翻译文件 (en, zh, fr, ru, ja, vi)
├── styles/                    # 全局样式 (Tailwind 主题)
└── assets/                    # 静态资源和图标
```

## 4. 关键设计模式

### 4.1 路由系统
- **文件路由**: TanStack Router 基于文件系统自动生成路由
- **鉴权保护**: 使用 `_authenticated` 布局路由，`beforeLoad` 守卫检查 `useAuthStore` 中的用户状态
- **嵌套路由**: 支持复杂的嵌套布局和路由

### 4.2 状态管理
- **Zustand**: 轻量级状态管理，支持持久化到 localStorage
- **Auth Store**: 在页面加载时从 localStorage 恢复用户数据，避免不必要的 API 调用
- **System Config Store**: 系统配置状态，支持实时更新

### 4.3 API 通信
- **Axios 实例**: 带有认证头和错误处理的拦截器
- **React Query**: 处理缓存、失效和重新获取
- **类型安全**: 所有 API 调用都有 TypeScript 类型定义

### 4.4 组件架构
- **功能模块化**: 基于功能的组织（`features/<name>/`）
- **共享组件**: `components/` 目录包含可重用的 UI 组件
- **工具函数**: `lib/` 目录包含通用工具函数
- **自定义 Hooks**: `hooks/` 目录包含可重用的业务逻辑

## 5. 构建与开发
```bash
# 安装依赖
bun install

# 开发服务器
bun run dev

# 生产构建
bun run build

# 国际化工具
bun run i18n:sync
```

## 6. 主题系统
- **默认主题**: `web/default/` - React 19, Rsbuild, Base UI, Tailwind CSS v4
- **经典主题**: `web/classic/` - React 18, Vite, Semi Design
- **主题切换**: 通过 `setting/system_setting/theme.go` 配置
- **嵌入式部署**: Go 后端通过 `//go:embed` 指令将前端构建产物打包进二进制文件

---

## 关联文件列表

### 前端核心文件
- `web/default/src/routes/__root.tsx` - 根路由
- `web/default/src/routes/index.tsx` - 首页
- `web/default/src/routes/setup/` - 设置向导
- `web/default/src/routes/(auth)/` - 认证路由
- `web/default/src/routes/_authenticated/` - 受保护路由

### 状态管理
- `web/default/src/stores/auth-store.ts` - 认证状态
- `web/default/src/stores/notification-store.ts` - 通知状态
- `web/default/src/stores/system-config-store.ts` - 系统配置状态

### 功能模块
- `web/default/src/features/wallet/` - 钱包功能
- `web/default/src/features/models/` - 模型管理
- `web/default/src/features/chat/` - 聊天功能
- `web/default/src/features/about/` - 关于页面
- `web/default/src/features/subscriptions/` - 订阅功能

### 国际化
- `web/default/src/i18n/` - 国际化配置
- `web/default/src/i18n/locales/` - 翻译文件

### 构建配置
- `web/default/package.json` - 依赖配置
- `web/default/rsbuild.config.ts` - Rsbuild 配置
- `web/default/tsconfig.json` - TypeScript 配置

### 后端集成
- `router/web-router.go` - 前端路由服务
- `common/embed-file-system.go` - 嵌入式文件系统
- `main.go` - 应用入口（嵌入前端）
