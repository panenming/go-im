## 安装redis
	1.下载(建议在/usr/local下)
	wget http://download.redis.io/releases/redis-4.0.6.tar.gz

	2.解压
	tar -xzvf redis-4.0.6.tar.gz

	3.编译并安装
	cd redis-4.0.6
	make test #如果没有什么错误再往下执行
	make && make install
	
	You need tcl 8.5 or newer in order to run the Redis test 解决方案
	wget http://downloads.sourceforge.net/tcl/tcl8.6.1-src.tar.gz  
	tar xzvf tcl8.6.1-src.tar.gz  -C /usr/local/  
	cd  /usr/local/tcl8.6.1/unix/  
	./configure  
	make && make install  

	4.相关配置
	cp redis.conf  /etc/ #复制配置文件
	vim /etc/redis.conf #修改配置文件
	#daemonize 改为yes 后台运行
	#port 端口号

	5.启动redis
	cd /usr/local/bin/
	redis-server /etc/redis.conf

	6.添加到开机自启
	echo "/usr/local/bin/redis-server /etc/redis.conf" >>/etc/rc.local 
	
## github.com/go-redis/redis 实现了redis的连接池，直接使用即可。


