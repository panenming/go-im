## Go语言网络层脚手架 
## 核心
link包的核心是Session，Session的字面含义是会话，就是一次对话过程。每一个连接的生命周期被表达为一个会话过程，这个过程中通讯双方的消息有来有往。

会话过程所用的具体通讯协议通过Codec接口解耦。通过Codec接口可以自定义通讯的IO实现方式，如：TCP、UDP、UNIX套接字、共享内存等，也可以自定义流的实现方式，如：压缩、加密、校验等，也可以实现自定义的协议格式，如：JSON、Gob、XML、Protobuf等。

在实际项目中，通常不会只有一个会话，所以link提供了几种不同的Session管理方式。

Manager是最基础的Session管理方式，它负责创建和管理一组Session。Manager是不与通讯形式关联的，与通讯有关联的Manager叫Server，它的行为比Manager更具体，它负责从net.Listener上接收新连接并创建Session，然后在独立的goroutine中处理来自新连接的消息。

link还提供了Channel用于对Session进行按需分组，Channel用key-value的形式管理Session，Channel的key类型通过代码生成的形式来实现自定义。
