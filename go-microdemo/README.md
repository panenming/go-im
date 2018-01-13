构建micro的微服务步骤<br>
一、环境搭建<br>
	1. 获取protoc文件，地址 https://github.com/google/protobuf/releases<br>
	2. go get -u github.com/micro/micro<br>
	3. go get -u github.com/micro/protobuf/protoc-gen-go<br>
	生成的micro.exe,protoc-gen-go.exe在%GOPATH%目录下<br>
二. 写.proto文件<br>
	根据业务需求书写正确的.proto文件<br>
三、生成.go文件<br>
	命令：protoc --go_out=plugins=micro:. message.proto<br>
	
需要注意的是启动的时候最好用cmd启动service.exe，因为liteide在杀死进程的时候不会触发micro的deregistry操作，会导致consul上存储很多无法访问的微服务<br>

启动micro api网关<br>
一、启动命令<br>
	micro api --address=:9000<br>
二、使用http访问微服务<br>
	http : post<br>
	url : http://127.0.0.1:9000/rpc<br>
	http header content-type : application/json<br>
	http body : {"service": "messageMicro","method": "MessageService.CreateMsg","request": {"Id": 123}}<br>
	retrun : {"Id": 32131}<br>
	
	这时表示可以正常访问<br>
	
注意使用的rpc默认走的协议是http，如果使用grpc等等协议需要修改<br>