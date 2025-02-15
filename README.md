# Performance indicator collection
Golang编写的计算机资源指标的采集工具，集成prometheus构建完整的监控系统\n

├── cmd/              # 存放可执行入口文件的主程序目录\n
│   └── myapp/        # 每个子目录对应一个可执行程序（如：`main.go`）\n
│       └── main.go\n
├── internal/         # 私有代码（仅限当前项目内部使用，外部无法导入）\n
│   └── auth/         # 例如：认证模块\n
│       └── auth.go\n
├── pkg/              # 公共库代码（可被其他项目导入）\n
│   └── utils/        # 例如：通用工具函数\n
│       └── helper.go\n
├── api/              # API 定义（如：Protobuf、gRPC、Swagger 等）\n
├── configs/          # 配置文件（如：YAML、TOML、JSON）\n
├── web/              # Web 相关静态资源（如：HTML/CSS/JS）\n
│   ├── static/\n
│   └── templates/\n
├── scripts/          # 辅助脚本（如：部署、构建脚本）\n
├── test/             # 测试相关（如：集成测试、测试数据）\n
├── go.mod            # Go 模块定义文件（依赖管理）\n
├── go.sum            # 依赖校验文件\n
└── README.md         # 项目文档\n