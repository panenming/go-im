go语言开发的im服务<br>
进行中...<br>
## 遇到的liteide不提示代码的问题，解决方案:
1. 从gocode下载地址下载gocode-master ，并解压 https://github.com/nsf/gocode<br>

2. go build gocode.go autocompletecontext.go autocompletefile.go client.go config.go cursorcontext.go decl.go declcache.go formatters.go os_windows.go package.go ripper.go rpc.go scope.go server.go utils.go type_alias_build_hack_18.go package_bin.go package_text.go<br>

3. 使用上一步生成的gocode.exe替换liteide bin目录下的gocode.exe文件<br>

原因是新版liteide的gocode.exe文件太老了，需要自己编译生成一个替换一下<br>

## 封装go网络层代码link
## link demo代码 go-tcpdemo/link_test

## 使用nsq做消息分发队列

## 使用minio做非结构化对象存储，包括图片，文件等等
![minio负载均衡思路](https://raw.githubusercontent.com/panenming/go-im/master/minioloadbalance.jpg)
	
	minio 启动不同server，minio下连接相同disk，nginx代理minio server实现负载均衡
	
## 实现minio在线播放器demo（minio官方）

## 经过disruptor 高并发洗礼（控制多线程锁、合理的数据结构、缓存行填充、内存屏障）

## go orm  http://jinzhu.github.io/gorm  中文文档 http://gorm.book.jasperxu.com/

## go inject github.com/facebookgo/inject

## minitor 
	github.com/divan/expvarmon(用于通过监控服务器端口来实现应用监控)
	使用方式 expvarmon -ports="1234" 
	
## 分布式id生成 snowflake 可行性

## 在线画板（考虑在群聊中使用，需要进行room分割）

## TiDB 数据库接入 替换mysql  elasticsearch 代替 mongodb  redis 可以使用

## android ios ui  https://gitee.com/jpush/aurora-imui 前端通信框架选 https://github.com/Tencent/mars

## 协议借鉴参考 github.com/Terry-Mao/goim 中的顶层传输协议 proto,二层协议待定

## 完成长连接初级阶段demo connector 。 go-im/main.go 启动 server， go-im/connector/client/client.go 启动测试client

## web 框架 go get -u github.com/kataras/iris 

