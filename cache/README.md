# cache 数据结构设计思路

1. 使用嵌入式数据库
2. 数据目录默认是 ~/.portainer-cli
3. 数据库的命名方式：通过 md5(host) 标识不同 portainer 实例的数据
4. 不同的数据结构建立不同的 bucket，目前有以下类型

- token 
- endpoint
- container
