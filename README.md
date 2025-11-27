# 博客系统部署与使用手册

本手册将指导您如何部署和使用本全栈博客系统。系统由 Go (Gin) 后端和 Vue 3 前端组成。

## 1. 环境要求

在开始之前，请确保您的服务器或本地环境已安装以下软件：

*   **Go**: 版本 1.25 或更高
*   **Node.js**: 版本 16 或更高 (推荐使用 LTS)
*   **MySQL**: 版本 8.0 或更高
*   **Git**: 用于代码版本管理

## 2. 数据库配置

1.  登录 MySQL 数据库。
2.  创建一个新的数据库（例如 `blog_db`）：
    ```sql
    CREATE DATABASE blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    ```
3.  确保您拥有该数据库的访问权限（用户名和密码）。

## 3. 后端部署 (Backend)

### 3.1 配置

1.  进入 `backend` 目录：
    ```bash
    cd backend
    ```
2.  创建 `.env` 配置文件（如果不存在），并填入以下内容：
    ```env
    SERVER_PORT=8080
    DB_USERNAME=root
    DB_PASSWORD=your_password
    DB_HOST=127.0.0.1
    DB_PORT=3306
    DB_NAME=blog_db
    API_SECRET=your_jwt_secret_key
    TOKEN_HOUR_LIFESPAN=24
    ```
    *请将 `your_password` 和 `your_jwt_secret_key` 替换为您的实际密码和密钥。*

### 3.2 运行

**开发模式：**
```bash
go run main.go
```

**生产模式（编译运行）：**
1.  编译二进制文件：
    ```bash
    go build -o blog-backend main.go
    ```
2.  运行二进制文件：
    *   **Windows**: `.\blog-backend.exe`
    *   **Linux/Mac**: `./blog-backend`

后端服务默认将在 `http://localhost:8080` 启动。

## 4. 前端部署 (Frontend)

### 4.1 配置

1.  进入 `frontend` 目录：
    ```bash
    cd frontend
    ```
2.  安装依赖：
    ```bash
    npm install
    ```
3.  配置 API 地址：
    *   打开 `src/api/axios.js`。
    *   确保 `baseURL` 指向您的后端地址（开发环境通常为 `http://localhost:8080/api/v1`）。

### 4.2 运行

**开发模式：**
```bash
npm run dev
```
访问地址通常为 `http://localhost:5173`。

**生产模式（构建）：**
1.  构建生产文件：
    ```bash
    npm run build
    ```
2.  构建完成后，生成的静态文件位于 `dist` 目录。
3.  您可以使用 Nginx、Apache 或任何静态文件服务器来托管 `dist` 目录。

    **Nginx 配置示例：**
    ```nginx
    server {
        listen 80;
        server_name your_domain.com;

        location / {
            root /path/to/your/frontend/dist;
            index index.html;
            try_files $uri $uri/ /index.html;
        }

        location /api {
            proxy_pass http://localhost:8080;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
        
        location /uploads {
            proxy_pass http://localhost:8080;
        }
    }
    ```

## 5. 使用指南

### 5.1 用户注册与登录
*   访问首页，点击右上角的 **Register** 注册新账号。
*   注册成功后，点击 **Login** 进行登录。
*   登录后，您将看到 **Write**（写文章）和 **Profile**（个人资料）入口。

### 5.2 个人资料设置
*   点击导航栏的 **Profile**。
*   您可以查看个人信息并上传/更新头像。

### 5.3 发布文章
*   点击 **Write** 进入文章编辑器。
*   **标题**：输入文章标题。
*   **分类**：选择现有分类或点击 **+ New Category** 创建新分类。
*   **标签**：输入标签，多个标签用逗号分隔（如 `tech, life`）。
*   **内容**：支持 **Markdown** 语法。右侧（或下方）会实时预览渲染效果。
*   点击 **Save Article** 发布。

### 5.4 浏览与互动
*   **首页**：展示文章列表。
    *   右侧侧边栏提供 **分类** 和 **标签云** 筛选功能。
    *   顶部搜索框支持按标题或内容搜索。
*   **文章详情**：
    *   阅读文章内容（Markdown 渲染）。
    *   **点赞**：点击 Like 按钮点赞。
    *   **评论**：在底部发表评论。
    *   **阅读量**：标题下方显示文章阅读次数。

### 5.5 管理文章
*   如果您是文章的作者，在文章详情页会看到 **Edit** 和 **Delete** 按钮，可对文章进行编辑或删除。

## 6. 常见问题

*   **图片无法显示？**
    确保后端 `uploads` 目录存在且有写入权限，并检查 Nginx/服务器是否正确代理了 `/uploads` 路径。
*   **刷新页面 404？**
    前端是单页应用 (SPA)，生产环境部署时需配置服务器（如 Nginx 的 `try_files`）将所有 404 请求重定向到 `index.html`。
