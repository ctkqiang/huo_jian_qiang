# 火尖枪- 专业登录页面凭证暴力破解测试工具

## 法律风险与免责声明（务必阅读）
```

【深求科技教育用途软件使用免责声明与法律风险声明】

一、法律依据：

本软件适用下列中华人民共和国法律法规及司法解释，用户必须全面遵守：

1. 《中华人民共和国网络安全法》（2017 年 6 月 1 日实施）
2. 《中华人民共和国刑法》（2015 年修订，2020 年修订版）
3. 《个人信息保护法》（2021 年 11 月 1 日起施行）
4. 《计算机信息网络国际联网安全保护管理办法》（工信部令第 36 号，2000 年）
5. 《网络安全等级保护条例》（2019 年 9 月 1 日起施行）
6. 最高人民法院《审理危害计算机信息系统安全刑事案件具体应用法律若干问题的解释》（法释〔2011〕7 号）

二、用户责任与行为规范：

1. 用户仅限将本软件用于合法的教育研究、技术学习及环境测试，严禁用于任何非法入侵、破坏、数据窃取、恶意攻击等违法犯罪行为。
2. 违反《刑法》第二百八十五条——非法侵入计算机信息系统的，处三年以下有期徒刑或者拘役，并处或者单处罚金；情节严重的，处三年以上七年以下有期徒刑。
3. 违反《刑法》第二百八十六条——制作、传播计算机病毒等破坏性程序，处三年以下有期徒刑、拘役或者罚金；情节严重的，处三年以上七年以下有期徒刑。
4. 违反《刑法》第二百八十七条——非法获取、出售或者提供个人信息，情节严重者可判处五年以下有期徒刑或拘役，并处罚金；情节特别严重的，处五年以上有期徒刑。
5. 违反《个人信息保护法》第六十七条，未采取必要措施导致个人信息泄露，视情节严重，最高可处以 500 万元人民币罚款。
6. 违反《网络安全法》第四十一条规定，擅自提供网络产品、服务存在安全隐患，责任单位依法承担法律责任。

三、免责声明：

1. 本软件仅供合法教育、技术研究及测试使用，开发者不承担任何因用户使用本软件导致的直接或间接损失，包括但不限于数据丢失、系统损坏、经济损失及法律责任。
2. 用户使用本软件即视为已充分理解本声明及法律风险，自愿承担所有责任。
3. 使用本软件请确保已获得目标系统的合法授权，严禁对无授权系统进行任何形式的攻击或测试。

四、风险提示：

1. 网络安全系全社会共同责任，任何非法攻击行为不仅违背法律，更破坏互联网生态环境。
2. 任何违法行为一经发现，将依法追究刑事责任，警方和司法机关可通过技术手段追踪定位违法者。
3. 任何软件滥用后果自负，开发者保留追究相关侵权责任的权利。

五、争议解决：

因使用本软件发生的任何争议，均适用中华人民共和国法律，由软件开发者所在地人民法院管辖。

六、特别声明：

请用户务必慎重下载及使用网络软件，避免因轻率行为导致不可挽回的法律后果。

网络安全不是儿戏，技术学习需守规矩，愿你我共同维护绿色互联网环境。

```

## 1. 项目概述

火尖枪是一款专为安全测试和性能验证设计的高性能HTTP批量请求工具。该工具能够自动化执行大规模的用户名/密码组合测试，支持高度可配置的请求参数和多线程并发处理。

### 1.1 核心功能和使用场景

#### 核心功能
- **批量凭证测试**：从指定的用户文件和密码文件中读取组合，自动生成并发送HTTP请求
- **多协议支持**：完整支持HTTP GET和POST方法
- **并发控制**：可配置的线程池实现高效并发请求
- **智能延迟**：支持请求间隔设置，避免触发服务器防护机制
- **详细日志**：实时记录请求状态、响应代码和响应内容
- **自动文件管理**：支持特殊标记(*U*, *P*)自动下载和生成测试文件

#### 典型应用场景
- **Web应用安全测试**：自动化检测弱密码和默认凭证
- **API压力测试**：模拟真实用户行为进行接口负载测试
- **渗透测试辅助**：在授权范围内进行安全漏洞验证
- **功能回归测试**：验证登录系统的稳定性和正确性

### 1.2 技术特点和优势

| 特性 | 优势说明 |
|------|----------|
| **高性能并发** | 基于Go协程的轻量级并发模型，支持数千并发请求 |
| **智能错误处理** | 自动识别429限流、超时、连接中断等异常状态 |
| **灵活配置** | 命令行参数驱动，无需修改代码即可适应不同测试场景 |
| **可扩展架构** | 模块化设计，易于添加新的协议支持和功能扩展 |
| **详细监控** | 实时进度显示和响应数据分析 |

## 2. 使用指南

### 2.1 完整操作流程

以下示例展示了针对WordPress登录页面的批量测试：

```bash
go run main.go \
  --Host "https://xxxx.xxxxxxxx" \
  -m "POST" \
  -u "*U*" \
  -p "*P*" \
  -a "username={user}&password={pass}" \
  -d 1 \
  -t 10
```

**执行流程详解：**

1. **初始化阶段**
   - 解析命令行参数和配置文件
   - 自动下载密码字典文件(`rockyou.txt`)
   - 生成默认用户列表文件(`users.txt`)
   - 显示配置摘要信息

2. **文件处理阶段**
```bash
   19:39:10 【火尖枪】: [信息] ┌─ 链接:         https://xxxxxxxxxxxx/xxxxx
   19:39:10 【火尖枪】: [信息] ├─ 用户文件:      downloads/data/users.txt
   19:39:10 【火尖枪】: [信息] ├─ 密码文件:      downloads/data/rockyou.txt
   19:39:10 【火尖枪】: [信息] ├─ 请求体:        username={user}&password={pass}
   19:39:10 【火尖枪】: [信息] ├─ 线程数:        10
   19:39:10 【火尖枪】: [信息] ├─ 请求方法:      POST
   19:39:10 【火尖枪】: [信息] └─ 延迟:          1 秒
   19:36:58 【火尖枪】: [信息] 状态码：200
   19:36:58 【火尖枪】: [信息] 响应长度：3160 字节
   19:36:58 【火尖枪】: [信息]     响应体预览：
   19:36:58 【火尖枪】: [信息]    ┌──────────────────────────────────────
  19:36:58 【火尖枪】: [信息]    │ username=admin&password=morgan
  19:36:58 【火尖枪】: [信息]    └──────────────────────────────────────
  19:36:58 【火尖枪】: [信息] -> 请求成功: 状态码=200, 响应长度=3160
```

3. **请求执行阶段**
   - 启动10个并发工作协程
   - 按1秒间隔发送请求
   - 实时监控和记录响应状态

4. **结果展示阶段**
   ```
   19:39:36 【火尖枪】: [信息] -> 请求成功: 状态码=429, 响应长度=1168
   19:39:36 【火尖枪】: [警告] 收到429状态码, 请求被限制, 建议增加延迟后重试
   19:39:36 【火尖枪】: [信息] 请求成功！
   19:39:36 【火尖枪】: [信息] 状态码：200
   19:39:36 【火尖枪】: [信息] 响应长度：3160 字节
   ```

### 2.2 参数说明

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `--Host` | string | 是 | 无 | 目标地址，支持HTTP URL(`http://...`)或MySQL DSN(`mysql://...`) |
| `-m` | string | 是 | GET | 请求方法：GET/POST（MySQL模式下忽略） |
| `-u` | string | 是 | 无 | 用户名字段文件，支持特殊标记`*U*`自动生成 |
| `-p` | string | 是 | 无 | 密码字段文件，支持特殊标记`*P*`自动下载 |
| `-a` | string | 是 | 无 | 请求体模板，支持{user}和{pass}占位符（MySQL模式可填任意值） |
| `-d` | int | 否 | 0 | 请求间隔时间（秒），用于避免触发限流 |
| `-t` | int | 否 | 1 | 并发线程数，控制同时执行的请求/连接数量 |
| `--help` | flag | 否 | 无 | 显示帮助信息 |

### 2.3 请求体模板示例

```bash
# WordPress样式登录
-a "log={user}&pwd={pass}&wp-submit=登录&redirect_to=..."

# 基础认证样式
-a "username={user}&password={pass}&submit=Login"

# JSON格式请求
-a '{"username":"{user}","password":"{pass}"}'
```

### 2.4 文件格式规范

**用户文件示例** (`downloads/data/users.txt`):
```
admin
administrator
root
user
test
```

**密码文件示例** (`downloads/data/rockyou.txt`):
```
password
123456
admin123
qwerty
letmein
```

### 2.5 MySQL数据库连接测试

火尖枪支持对MySQL数据库进行用户名和密码的批量测试。当检测到`--Host`参数包含`mysql`关键字时，系统会自动切换到MySQL测试模式。

#### 2.5.1 命令语法示例

```bash
go run main.go \
  --Host "mysql://root:123456@tcp(127.0.0.1:3306)/" \
  -u "*U*" \
  -p "*P*" \
  -a "dummy" \
  -d 1 \
  -t 10
```

**参数解释：**
- `--Host`: MySQL连接字符串，支持标准DSN格式
- `-u`: 用户名字典文件路径，支持`*U*`自动生成默认文件
- `-p`: 密码字典文件路径，支持`*P*`自动下载rockyou.txt
- `-a`: 请求体参数（MySQL模式下可填写任意值，但必须提供）
- `-d`: 请求间隔时间（秒），避免触发数据库连接限制
- `-t`: 并发线程数，控制同时建立的数据库连接数量

#### 2.5.2 MySQL连接字符串格式

火尖枪支持多种MySQL连接字符串格式：

**标准DSN格式：**
```
mysql://用户名:密码@tcp(主机:端口)/数据库名
```

**示例：**
- `mysql://root:password@tcp(localhost:3306)/test`
- `mysql://admin:admin123@tcp(192.168.1.100:3306)/production`
- `mysql://user:@tcp(127.0.0.1:3306)/` (空密码)
- `mysql://:password@tcp(127.0.0.1:3306)/` (空用户名)

**简化格式（自动补全）：**
- `mysql://127.0.0.1:3306` → 自动转换为 `mysql://root:@tcp(127.0.0.1:3306)/`
- `mysql://localhost` → 自动转换为 `mysql://root:@tcp(localhost:3306)/`

#### 2.5.3 参数说明参考

MySQL测试使用与HTTP测试相同的命令行参数，详细参数说明请参考[2.2 参数说明](#22-参数说明)章节。

**MySQL测试特别注意事项：**
- `--Host`参数必须使用MySQL DSN格式（`mysql://...`）
- `-m`参数在MySQL模式下被忽略，但必须提供
- `-a`参数在MySQL模式下可填写任意值，但必须提供
- 建议根据数据库性能合理设置`-t`（线程数）和`-d`（延迟）参数

#### 2.5.4 安全注意事项

1. **合法授权**：仅对拥有合法授权的数据库进行测试，严禁未经授权访问他人数据库
2. **生产环境**：避免直接在生产数据库上测试，建议使用测试环境或隔离的数据库实例
3. **连接限制**：合理设置`-d`延迟参数，避免对数据库造成过大压力或触发防火墙规则
4. **敏感信息**：MySQL连接字符串中包含用户名和密码，避免在日志或公共场合泄露
5. **网络隔离**：确保测试环境与生产环境网络隔离，防止误操作影响生产服务

#### 2.5.5 环境配置示例

**测试环境配置：**
```bash
# 本地MySQL测试
export MYSQL_HOST="mysql://root:@tcp(localhost:3306)/testdb"
export USERS_FILE="downloads/data/users.txt"
export PASSWORDS_FILE="downloads/data/rockyou.txt"

go run main.go \
  --Host "$MYSQL_HOST" \
  -u "$USERS_FILE" \
  -p "$PASSWORDS_FILE" \
  -a "dummy" \
  -d 2 \
  -t 5
```

**Docker环境测试：**
```bash
# 连接Docker中的MySQL容器
go run main.go \
  --Host "mysql://root:testpassword@tcp(172.17.0.2:3306)/appdb" \
  -u "*U*" \
  -p "*P*" \
  -a "dummy" \
  -d 1 \
  -t 8
```

**云数据库测试：**
```bash
# AWS RDS MySQL测试
go run main.go \
  --Host "mysql://admin:CloudPass123@tcp(database-1.xxx.us-east-1.rds.amazonaws.com:3306)/application" \
  -u "custom_users.txt" \
  -p "custom_passwords.txt" \
  -a "dummy" \
  -d 3 \
  -t 4
```

#### 2.5.6 故障排除指南

| 问题现象 | 可能原因 | 解决方案 |
|----------|----------|----------|
| **连接被拒绝** | 1. 数据库服务未运行<br>2. 防火墙阻止连接<br>3. 主机地址错误 | 1. 检查MySQL服务状态<br>2. 验证防火墙规则(3306端口)<br>3. 确认主机IP和端口正确 |
| **认证失败** | 1. 用户名或密码错误<br>2. 用户没有远程连接权限<br>3. 密码策略限制 | 1. 使用正确的凭证组合<br>2. 授予用户远程访问权限<br>3. 检查密码复杂度要求 |
| **连接超时** | 1. 网络延迟过高<br>2. 数据库负载过大<br>3. 连接池耗尽 | 1. 增加`-d`延迟参数<br>2. 减少`-t`线程数<br>3. 检查数据库性能 |
| **协议不支持** | 1. 使用了错误的连接格式<br>2. MySQL版本不兼容 | 1. 使用正确的DSN格式<br>2. 确认MySQL版本支持 |
| **文件不存在** | 1. 字典文件路径错误<br>2. 自动下载失败 | 1. 检查文件路径是否正确<br>2. 手动下载字典文件 |

#### 2.5.7 敏感参数警告

**重要安全警告**：

1. **连接字符串泄露风险**：MySQL DSN包含数据库凭证，严禁提交到版本控制系统或公开分享
2. **字典文件内容**：密码字典文件可能包含常见密码，妥善保管防止被恶意利用
3. **日志记录**：工具会记录连接尝试信息，测试完成后建议清理日志文件
4. **环境变量**：使用环境变量存储敏感参数，避免在命令行历史中暴露

#### 2.5.8 环境变量迁移说明

为提升安全性，建议将敏感参数迁移到环境变量中：

**传统方式（不安全）：**
```bash
go run main.go --Host "mysql://root:MyPassword123@tcp(127.0.0.1:3306)/" ...
```

**环境变量方式（推荐）：**
```bash
export DB_DSN="mysql://root:MyPassword123@tcp(127.0.0.1:3306)/"
go run main.go --Host "$DB_DSN" ...
```

**高级配置（.env文件）：**
```bash
# 创建.env文件
echo 'DB_DSN="mysql://root:MyPassword123@tcp(127.0.0.1:3306)/"' > .env
echo 'USERS_FILE="downloads/data/users.txt"' >> .env
echo 'PASSWORDS_FILE="downloads/data/rockyou.txt"' >> .env

# 使用环境变量
export $(cat .env | xargs)
go run main.go --Host "$DB_DSN" -u "$USERS_FILE" -p "$PASSWORDS_FILE" ...
```

#### 2.5.9 连接验证方法

在开始批量测试前，建议先验证数据库连接是否正常：

**手动验证连接：**
```bash
# 使用单个已知正确的凭证测试
go run main.go \    
  -Host "mysql://db.example.com:3306" \
  -u "*U*" \
  -p "*P*" \
  -a "db_specific_test" \
  -t 8
```

**使用测试脚本验证：**
```bash
# 创建测试脚本 test_mysql.sh
#!/bin/bash
DSN="$1"
USER="$2"
PASS="$3"
echo "测试连接: $USER@$DSN"
# 使用工具验证连接
```

#### 2.5.10 线程数推荐

根据不同的测试场景，推荐以下线程数配置：

| 测试场景 | 推荐线程数 | 推荐延迟 | 说明 |
|----------|------------|----------|------|
| **本地数据库** | 10-20 | 0-1秒 | 低延迟环境可适度提高并发 |
| **远程数据库** | 5-10 | 2-5秒 | 网络延迟较高，减少并发 |
| **云数据库** | 8-15 | 1-3秒 | 云服务弹性较好，中等并发 |
| **敏感生产环境** | 2-5 | 5-10秒 | 最小化对生产系统的影响 |
| **性能压测** | 20-50 | 0秒 | 仅在授权下进行压力测试 |

> **性能提示**：过高的线程数可能导致数据库连接池耗尽或触发安全防护，建议从低线程数开始逐步增加。

## 3. 技术原理

### 3.1 核心算法和工作原理
![原理](./doc/huojianqiang_sequence.png)

#### 3.1.1 多线程调度算法

```go
// 伪代码：并发请求调度核心逻辑
func processFiles(cfg *Config) error {
    users := readLines(cfg.UsersFile)    // 读取用户列表
    passwords := readLines(cfg.PasswordsFile) // 读取密码列表
    
    semaphore := make(chan struct{}, cfg.Threads) // 创建信号量控制并发数
    
    for _, user := range users {
        for _, password := range passwords {
            wg.Add(1)
            semaphore <- struct{}{} // 获取信号量
            
            go func(u, p string) {
                defer wg.Done()
                defer func() { <-semaphore }() // 释放信号量
                
                if cfg.Delay > 0 {
                    time.Sleep(time.Duration(cfg.Delay) * time.Second)
                }
                
                processRequest(cfg, u, p) // 执行单个请求
            }(user, password)
        }
    }
    
    wg.Wait() // 等待所有协程完成
    return nil
}
```

#### 3.1.2 限流检测与处理机制

当检测到HTTP 429状态码（请求过多）时，系统会自动记录警告并提供优化建议：

```
19:39:36 【火尖枪】: [警告] 收到429状态码, 请求被限制, 建议增加延迟后重试
```

**限流处理策略：**
1. **立即降频**：当前线程暂停执行
2. **智能建议**：推荐增加`-d`参数值
3. **状态记录**：标记该请求组合为受限状态
4. **继续执行**：不中断整体测试流程

#### 3.1.3 请求构建与发送

```go
func processRequest(cfg *Config, user, password string) {
    // 1. 模板变量替换
    requestBody := strings.ReplaceAll(cfg.RequestBody, "{user}", user)
    requestBody = strings.ReplaceAll(requestBody, "{pass}", password)
    
    // 2. 根据方法选择发送方式
    switch strings.ToUpper(cfg.Method) {
    case "POST":
        response, statusCode, err = http.PostRequest(cfg.Url, requestBody, 30)
    case "GET":
        response, statusCode, err = http.GetRequest(cfg.Url, requestBody, 30)
    }
    
    // 3. 结果分析与记录
    if statusCode == 200 {
        logger.Infof("✅ 请求成功！用户=%s, 密码=%s", user, password)
    } else if statusCode == 429 {
        logger.Warnf("⚠️ 请求受限！建议增加延迟参数")
    }
}
```

### 3.2 关键数据结构

```go
// 配置结构体
type Config struct {
    Url           string // 目标URL
    UsersFile     string // 用户文件路径
    PasswordsFile string // 密码文件路径
    RequestBody   string // 请求体模板
    Delay         int    // 请求间隔(秒)
    Threads       int    // 并发线程数
    Method        string // 请求方法
}

// 单个请求结果
type RequestResult struct {
    Username    string    // 测试用户名
    Password    string    // 测试密码
    StatusCode  int       // HTTP状态码
    Response    string    // 响应内容
    Success     bool      // 是否成功
    Timestamp   time.Time // 请求时间
}
```

## 4. 系统架构

### 4.1 整体设计思路

huō jiàn qiáng采用**模块化管道架构**，各组件通过清晰定义的接口进行通信。系统设计遵循单一职责原则，每个模块专注于特定功能领域。

### 4.2 核心组件架构

```
┌─────────────────────────────────────────────────────────────┐
│                     命令行接口层 (CLI)                        │
│  • 参数解析与验证                                          │
│  • 配置文件管理                                            │
│  • 用户交互反馈                                            │
└───────────────────────┬─────────────────────────────────────┘
                        │
┌───────────────────────▼─────────────────────────────────────┐
│                    业务逻辑层 (Core)                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ 文件处理器  │  │ 请求调度器  │  │ 结果分析器  │        │
│  │ • 文件读取  │  │ • 并发控制  │  │ • 状态解析  │        │
│  │ • 内容解析  │  │ • 速率限制  │  │ • 数据统计  │        │
│  │ • 格式验证  │  │ • 错误重试  │  │ • 报告生成  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└───────────────────────┬─────────────────────────────────────┘
                        │
┌───────────────────────▼─────────────────────────────────────┐
│                    基础设施层 (Infrastructure)                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ HTTP客户端  │  │ 文件下载器  │  │ 日志系统    │        │
│  │ • 请求发送  │  │ • 远程获取  │  │ • 分级日志  │        │
│  │ • 响应处理  │  │ • 本地缓存  │  │ • 格式美化  │        │
│  │ • 超时控制  │  │ • 校验验证  │  │ • 文件输出  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

### 4.3 数据流与组件交互

1. **配置阶段**：CLI层解析参数 → 验证配置有效性 → 初始化各组件
2. **准备阶段**：文件下载器获取字典 → 文件处理器验证格式 → 生成测试矩阵
3. **执行阶段**：调度器分配任务 → HTTP客户端发送请求 → 记录响应结果
4. **分析阶段**：结果分析器处理数据 → 生成统计报告 → 输出测试摘要

### 4.4 异常处理机制

系统实现了分层的异常处理策略：

- **网络层异常**：自动重试机制，可配置重试次数
- **应用层异常**：状态码分析，提供针对性建议
- **系统层异常**：资源释放保障，避免内存泄漏
- **用户层异常**：友好错误提示，指导正确使用

## 5. 性能优化建议

### 5.1 并发参数调优

| 场景类型 | 推荐线程数 | 推荐延迟 | 说明 |
|----------|------------|----------|------|
| 本地测试 | 10-20 | 0-1秒 | 低延迟环境可提高并发 |
| 远程测试 | 5-10 | 2-5秒 | 避免触发目标服务器防护 |
| 云服务测试 | 15-30 | 1-3秒 | 利用云服务弹性扩展能力 |
| 敏感系统 | 3-5 | 5-10秒 | 对生产系统的最小影响 |

### 5.2 内存使用优化

- **流式文件读取**：避免一次性加载大文件到内存
- **结果分页存储**：定期将结果写入磁盘释放内存
- **连接池复用**：重用HTTP连接减少开销
- **协程数量控制**：根据系统资源动态调整

## 6. 安全与合规

### 6.1 使用注意事项

1. **合法授权**：仅在被明确授权的系统上进行测试
2. **风险告知**：测试前评估对目标系统的影响
3. **数据保护**：妥善保管测试过程中产生的敏感数据
4. **速率控制**：遵守目标系统的服务条款和使用限制

### 6.2 合规性建议

- 在测试前获取书面授权
- 明确测试范围和时限
- 记录所有测试活动以备审计
- 测试完成后清理测试数据

---

##  开源项目赞助计划

### 用捐赠助力发展

感谢您使用本项目！您的支持是开源持续发展的核心动力。  
每一份捐赠都将直接用于：  
- 服务器与基础设施维护（魔法城堡的维修费哟~）  
- 新功能开发与版本迭代（魔法技能树要升级哒~）  
- 文档优化与社区建设（魔法图书馆要扩建呀~）

点滴支持皆能汇聚成海，让我们共同打造更强大的开源工具！  
（小仙子们在向你比心哟~）

---

### 全球捐赠通道

#### 国内用户

<div align="center" style="margin: 40px 0">

<div align="center">
<table>
<tr>
<td align="center" width="300">
<img src="https://github.com/ctkqiang/ctkqiang/blob/main/assets/IMG_9863.jpg?raw=true" width="200" />
<br />
<strong>支付宝</strong>（小企鹅在收金币哟~）
</td>
<td align="center" width="300">
<img src="https://github.com/ctkqiang/ctkqiang/blob/main/assets/IMG_9859.JPG?raw=true" width="200" />
<br />
<strong>微信支付</strong>（小绿龙在收金币哟~）
</td>
</tr>
</table>
</div>
</div>

#### 国际用户

<div align="center" style="margin: 40px 0">
  <a href="https://qr.alipay.com/fkx19369scgxdrkv8mxso92" target="_blank">
    <img src="https://img.shields.io/badge/Alipay-全球支付-00A1E9?style=flat-square&logo=alipay&logoColor=white&labelColor=008CD7">
  </a>
  
  <a href="https://ko-fi.com/F1F5VCZJU" target="_blank">
    <img src="https://img.shields.io/badge/Ko--fi-买杯咖啡-FF5E5B?style=flat-square&logo=ko-fi&logoColor=white">
  </a>
  
  <a href="https://www.paypal.com/paypalme/ctkqiang" target="_blank">
    <img src="https://img.shields.io/badge/PayPal-安全支付-00457C?style=flat-square&logo=paypal&logoColor=white">
  </a>
  
  <a href="https://donate.stripe.com/00gg2nefu6TK1LqeUY" target="_blank">
    <img src="https://img.shields.io/badge/Stripe-企业级支付-626CD9?style=flat-square&logo=stripe&logoColor=white">
  </a>
</div>

---

### 开发者社交图谱

#### 技术交流

<div align="center" style="margin: 20px 0">
  <a href="https://github.com/ctkqiang" target="_blank">
    <img src="https://img.shields.io/badge/GitHub-开源仓库-181717?style=for-the-badge&logo=github">
  </a>
  
  <a href="https://stackoverflow.com/users/10758321/%e9%92%9f%e6%99%ba%e5%bc%ba" target="_blank">
    <img src="https://img.shields.io/badge/Stack_Overflow-技术问答-F58025?style=for-the-badge&logo=stackoverflow">
  </a>
  
  <a href="https://www.linkedin.com/in/ctkqiang/" target="_blank">
    <img src="https://img.shields.io/badge/LinkedIn-职业网络-0A66C2?style=for-the-badge&logo=linkedin">
  </a>
</div>

#### 社交互动

<div align="center" style="margin: 20px 0">
  <a href="https://www.instagram.com/ctkqiang" target="_blank">
    <img src="https://img.shields.io/badge/Instagram-生活瞬间-E4405F?style=for-the-badge&logo=instagram">
  </a>
  
  <a href="https://twitch.tv/ctkqiang" target="_blank">
    <img src="https://img.shields.io/badge/Twitch-技术直播-9146FF?style=for-the-badge&logo=twitch">
  </a>
  
  <a href="https://github.com/ctkqiang/ctkqiang/blob/main/assets/IMG_9245.JPG?raw=true" target="_blank">
    <img src="https://img.shields.io/badge/微信公众号-钟智强-07C160?style=for-the-badge&logo=wechat">
  </a>

<div align="center" style="margin: 20px 0">
  <img src="https://raw.githubusercontent.com/ctkqiang/ctkqiang/refs/heads/main/assets/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E7%99%BD%E8%89%B2%E7%89%88.png" alt="扫码赞助-白色版" width="300"/>
  <br />
  <em>微信扫码，一键赞助（小精灵们在等你哟~）</em>
</div>
</div>

---

致极客与未来的你

> "世界由代码驱动，安全靠你我守护。"

无论你是网络安全研究员、CTF 挑战者、自由极客，还是热爱数学与工程的探索者，这个项目都向你敞开怀抱。
欢迎你 fork、魔改、重构、注入灵感，它是工具，也是信仰。

但请铭记心底：

> **技术本无善恶，使用才有底线。**
> 愿你在这个开源项目中，找到属于自己的价值与乐趣。
