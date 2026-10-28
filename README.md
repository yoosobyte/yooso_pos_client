# 悦速进销存

基于 **Go + Wails v3 + Vue 3 + SQLite** 的桌面端进销存 / 收银系统，单机运行、数据本地存储，可打包为 Windows / macOS 桌面应用。

<p align="center">
  <img src="docs/screenshots/01-dashboard.png" width="900" alt="工作台">
</p>

## 界面预览

| 收银台 | 商品列表 |
| :---: | :---: |
| [![收银台](docs/screenshots/02-cashier.png)](docs/screenshots/02-cashier.png) | [![商品列表](docs/screenshots/03-goods.png)](docs/screenshots/03-goods.png) |
| **商品分类 / 单位** | **库存总览** |
| [![商品分类](docs/screenshots/04-category.png)](docs/screenshots/04-category.png) | [![库存总览](docs/screenshots/05-stock.png)](docs/screenshots/05-stock.png) |
| **出入库明细** | **采购入库** |
| [![出入库明细](docs/screenshots/06-stock-log.png)](docs/screenshots/06-stock-log.png) | [![采购入库](docs/screenshots/07-purchase.png)](docs/screenshots/07-purchase.png) |
| **采购记录** | **供应商管理** |
| [![采购记录](docs/screenshots/08-purchase-record.png)](docs/screenshots/08-purchase-record.png) | [![供应商管理](docs/screenshots/09-supplier.png)](docs/screenshots/09-supplier.png) |
| **销售订单** | **客户管理** |
| [![销售订单](docs/screenshots/10-sales-record.png)](docs/screenshots/10-sales-record.png) | [![客户管理](docs/screenshots/11-customer.png)](docs/screenshots/11-customer.png) |
| **销售报表** | **采购报表** |
| [![销售报表](docs/screenshots/12-report-sales.png)](docs/screenshots/12-report-sales.png) | [![采购报表](docs/screenshots/13-report-purchase.png)](docs/screenshots/13-report-purchase.png) |

<details>
<summary>更多截图（内置官网入口）</summary>

<p align="center">
  <img src="docs/screenshots/14-website.png" width="900" alt="内置官网">
</p>

</details>

## 快速开始

### 开发模式
```bash
wails3 dev
```

### Win生产构建
```bash
wails3 build
```

### Mac生产构建
```bash
wails3 task darwin:package
```

## 常用命令

 命令  说明 
------------
 `wails3 dev`  开发模式（热重载） 
 `wails3 build`  构建生产版本 
 `wails3 task common:update:build-assets`  更新应用信息/图标（修改配置后执行） 
 `wails3 task common:update:build-assets && wails3 build`  修改配置后完整构建 

## 修改配置后

```bash
# 修改了 build/config.yml（应用名称、版本、版权）或替换了图标后
wails3 task common:update:build-assets && wails3 build
```

## 项目结构

```
.
├── build/
│   ├── config.yml       # 应用配置（名称、版本、版权）
│   ├── appicon.png      # 应用图标（1024x1024）
│   ├── Taskfile.yml     # 通用构建任务
│   ├── darwin/
│   │   └── Taskfile.yml # macOS 构建任务
│   └── windows/
│       └── Taskfile.yml # Windows 构建任务
├── bin/                 # 构建产物（不提交 Git）
├── docs/screenshots/    # README 截图
├── frontend/            # 前端代码（Vue/React）
├── main.go              # Go 入口
├── app.go               # 应用逻辑
├── Taskfile.yml         # 主构建配置
└── wails.json           # Wails 配置（可选）
```

## Git 提交规则

 ✅ 提交（源文件）  ❌ 不提交（自动生成） 
-----------------------------------------
 `build/config.yml`  `build/bin/` 
 `build/Taskfile.yml`  `build/darwin/Info.plist` 
 `build/darwin/Taskfile.yml`  `build/darwin/Info.dev.plist` 
 `build/windows/Taskfile.yml`  `build/darwin/Assets.car` 
 `build/appicon.png`  `build/darwin/*.icns` 
 `frontend/` 源码  `build/windows/*.exe` 
 `*.go` 源码  `build/windows/*.ico` 
 `go.mod` / `go.sum`  `frontend/node_modules/` 
 `Taskfile.yml`  `frontend/dist/` 

## 修改应用信息

编辑 `build/config.yml`：

```yaml
info:
  companyName: "悦速互联"
  productName: "悦速进销存"
  productIdentifier: "com.yooso.jxc"
  description: "悦速进销存 - 专业进销存管理软件"
  copyright: "© 2026-2030 悦速互联. All Rights Reserved."
  comments: "悦速进销存  定制开发: 17337771173(微信同号)"
  version: "1.0.0"
```

修改后执行：
```bash
wails3 task common:update:build-assets && wails3 build
```

## 更换应用图标

1. 准备一张 **1024x1024** 的 PNG 图片
2. 替换 `build/appicon.png`
3. 执行：
```bash
wails3 task common:update:build-assets && wails3 build
```

## 打包发布

### macOS
```bash
wails3 build
# 生成 bin/悦速进销存.app
```

### Windows（交叉编译）
```bash
wails3 build GOOS=windows
# 生成 bin/悦速进销存.exe
```

## 常见问题

### 1. 启动报错 `Info.dev.plist: No such file or directory`
```bash
wails3 task common:update:build-assets
```

### 2. 图标没更新
```bash
# 清除缓存后重新生成
rm -rf bin/ build/darwin/*.icns build/darwin/Assets.car
wails3 task common:update:build-assets
wails3 build
```

### 3. 端口被占用
```bash
# 杀掉占用 9245 端口的进程
lsof -ti :9245  xargs kill -9
wails3 dev
```

### 4. 构建报错 `open windows/wails.exe.manifest: The system cannot find the file specified`
说明 `build/windows/` 下的生成资源（`wails.exe.manifest`、`info.json`、`nsis/` 等）缺失，
常见于新克隆仓库、`git clean -xdf` 或手动删过 `build/` 之后。此时 `wails3 dev` / `wails3 build`
会在生成 `.syso` 这一步直接失败，程序自然起不来。补回资源后再构建：

```bash
wails3 task common:update:build-assets
wails3 build        # 或 wails3 dev
```

## 联系

定制开发: 17337771173（微信同号）

## 许可证

本项目基于 [MIT License](LICENSE) 开源，可自由使用、修改和商用，保留版权声明即可。