# 更新代码地图

分析 Go 项目结构并更新架构文档：

1. 扫描项目结构（cmd, internal, pkg）、包依赖关系、接口定义和结构体

2. 以以下格式生成简洁的代码地图：
   * codemaps/architecture.md - 整体架构、模块依赖与层级划分
   * codemaps/packages.md - 核心包、服务接口与业务逻辑
   * codemaps/models.md - 数据结构体 (Structs) 与领域模型
   * codemaps/api.md - API 路由与处理函数定义

3. 计算与之前版本的差异百分比

4. 如果变更 > 30%，则在更新前请求用户批准

5. 为每个代码地图添加新鲜度时间戳

6. 将报告保存到 .reports/codemap-diff.txt

使用 Go 语言特性（包、接口、结构体）进行分析。专注于高层组件关系，而非具体实现细节。
