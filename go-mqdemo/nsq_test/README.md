## nsq 测试demo
## nsq 安装
1. http://nsq.io/deployment/installing.html 下载对应的版本
2. 集群启动
	nsq_start.sh
	
	#服务启动
	#注意更改一下 --data-path 所指定的数据存放路径，否则会无法运行。
	echo '删除日志文件'
	rm -f nsqlookupd.log
	rm -f nsqd1.log
	rm -f nsqd2.log
	rm -f nsqadmin.log
	
	echo '启动nsq服务'
	nohup nsqlookupd >nsqlookupd.log 2>&1&
	
	echo '启动nsqd服务'
	nohup nsqd --lookupd-tcp-address=0.0.0.0:4160 -tcp-address="0.0.0.0:4150"  --data-path=~/nsqd1  >nsqd1.log 2>&1&
	nohup nsqd --lookupd-tcp-address=0.0.0.0:4160 -tcp-address="0.0.0.0:4152" -http-address="0.0.0.0:4153" --data-path=~/nsqd2 >nsqd2.log 2>&1&
	
	echo '启动nsqdadmin服务'
	nohup nsqadmin --lookupd-http-address=0.0.0.0:4161 >nsqadmin.log 2>&1&
3. 关闭
	nsq_shutdown.sh
	#nsq_shutdown.sh
	#服务停止
	ps -ef | grep nsq| grep -v grep | awk '{print $2}' | xargs kill -2 
4. 查看
	运行后，访问本机:4171端口，就能够通过web页面进行查看：
5. 概念
    Topic ：一个topic就是程序发布消息的一个逻辑键，当程序第一次发布消息时就会创建topic。

    Channels ：channel与消费者相关，是消费者之间的负载均衡，channel在某种意义上来说是一个“队列”。每当一个发布者发送一条消息到一个topic，消息会被复制到所有消费者连接的channel上，消费者通过这个特殊的channel读取消息，实际上，在消费者第一次订阅时就会创建channel。Channel会将消息进行排列，如果没有消费者读取消息，消息首先会在内存中排队，当量太大时就会被保存到磁盘中。

    Messages：消息构成了我们数据流的中坚力量，消费者可以选择结束消息，表明它们正在被正常处理，或者重新将他们排队待到后面再进行处理。每个消息包含传递尝试的次数，当消息传递超过一定的阀值次数时，我们应该放弃这些消息，或者作为额外消息进行处理。

    nsqd：nsqd 是一个守护进程，负责接收，排队，投递消息给客户端。它可以独立运行，不过通常它是由 nsqlookupd 实例所在集群配置的（它在这能声明 topics 和 channels，以便大家能找到）。

    nsqlookupd：nsqlookupd 是守护进程负责管理拓扑信息。客户端通过查询 nsqlookupd 来发现指定话题（topic）的生产者，并且 nsqd 节点广播话题（topic）和通道（channel）信息。有两个接口：TCP 接口，nsqd 用它来广播。HTTP 接口，客户端用它来发现和管理。

    nsqadmin：nsqadmin 是一套 WEB UI，用来汇集集群的实时统计，并执行不同的管理任务。

6.   常用工具类：
    nsq_to _file：消费指定的话题（topic）/通道（channel），并写到文件中，有选择的滚动和/或压缩文件。

    nsq_to _http：消费指定的话题（topic）/通道（channel）和执行 HTTP requests (GET/POST) 到指定的端点。

    nsq_to _nsq：消费者指定的话题/通道和重发布消息到目的地 nsqd 通过 TCP。
	
注意：nsq消息的到达是无序的