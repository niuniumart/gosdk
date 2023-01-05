## 功能列表
go web sdk 提供以下功能：
1. 自动捕获panic；
2. 自动接入分布式日志；
3. 提供协程级web 请求链路追踪能力；
4. 自动接入普罗米修斯，上报所有接口的请求数、错误数、耗时、系统错误情况；
5. 提供配置中心能力。

## 使用说明
提供了日志功能seelog，如果原来使用的开源seelog，可以使用以下命令全局替换。

sed -i "" "s#github.com/cihub/seelog#github.com/niuniumart/gosdk/seelog#g" `grep "github\.com\/cihub\/seelog"  -rl ./`
