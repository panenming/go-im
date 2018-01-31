## centos 运行 tidb
	1.下载 wget http://download.pingcap.org/tidb-latest-linux-amd64.tar.gz

	2. 启动tidb
	mv tidb-latest-linux-amd64.tar.gz  /usr/local/
	cd /usr/local/
	tar xf tidb-latest-linux-amd64.tar.gz 
	ln -s tidb-latest-linux-amd64 tidb
	cd /usr/local/tidb
	./bin/pd-server --data-dir=pd --log-file=pd.log &
	./bin/tikv-server --pd="127.0.0.1:2379" --data-dir=tikv --log-file=tikv.log &
	./bin/tidb-server --store=tikv --path="127.0.0.1:2379" --log-file=tidb.log &
	
	3.  yum install -y mysql-community-client.x86_64 
	
 	4. 使用： mysql -h 127.0.0.1 -P 4000 -u root -D test
	5.功能测试
	grant all privileges on *.* to 'oracle'@'%' identified by 'oracle' with grant option;
	flush  privileges;



	创建用户
	CREATE USER 'ggs'@'%' IDENTIFIED BY 'oracle';
	set password for ggs=password('oracle');
	
	GRANT ALL PRIVILEGES ON *.* TO 'ggs'@'%' with grant option;
	flush  privileges;
	
	
	root用户赋予外网访问权限
	grant all privileges on *.* to 'root'@'localhost' identified by 'oracle' with grant option;
	grant all privileges on *.* to 'root'@'%' identified by 'oracle' with grant option;
	grant all privileges on *.* to 'root'@'127.0.0.1' identified by 'oracle' with grant option;
	flush  privileges;
	
## 使用gorm验证tidb可行性

