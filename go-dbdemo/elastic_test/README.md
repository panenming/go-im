## centos 安装 es
	centos 7 yum -ivh xx.rpm
	运行：sudo systemctl start elasticsearch.service
	
	需要修改配置文件/etc/elasticsearch/elasticsearch.yml
	中network.host，把后面改为0.0.0.0或者虚拟机ip地址 允许其他ip访问
	
## github.com/olivere/elastic 仅支持 http访问模式

## 坑 
	搜了一下应该是5.x后对排序，聚合这些操作用单独的数据结构(fielddata)缓存到内存里了，需要单独开启，官方解释在此fielddata

	简单来说就是在聚合前执行如下操作
	
	PUT megacorp/_mapping/employee/
	{
	  "properties": {
	    "interests": { 
	      "type":     "text",
	      "fielddata": true
	    }
	  }
	}
	
	
	[]string 不能直接转化为 []interface