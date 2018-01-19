## 测试link
协议json
压力测试

## 调整windows最大连接数
Windows 下单机最大TCP连接数
调整系统参数来调整单机的最大TCP连接数，Windows 下单机的TCP连接数有多个参数共同决定：
以下都是通过修改注册表[HKEY_LOCAL_MACHINE \System \CurrentControlSet \Services \Tcpip \Parameters]
 
**1.最大TCP连接数 **     TcpNumConnections
**2.TCP关闭延迟时间 **   TCPTimedWaitDelay    (30-240)s
**3.最大动态端口数 **  MaxUserPort  (Default = 5000, Max = 65534) TCP客户端和服务器连接时，客户端必须分配一个动态端口，默认情况下这个动态端口的分配范围为 1024-5000 ，也就是说默认情况下，客户端最多可以同时发起3977 Socket 连接
**4.最大TCB 数量   MaxFreeTcbs**
系统为每个TCP 连接分配一个TCP 控制块(TCP control block or TCB)，这个控制块用于缓存TCP连接的一些参数，每个TCB需要分配 0.5 KB的pagepool 和 0.5KB 的Non-pagepool，也就说，每个TCP连接会占用 1KB 的系统内存。
非Server版本，MaxFreeTcbs 的默认值为1000 （64M 以上物理内存）Server 版本，这个的默认值为 2000。也就是说，默认情况下，Server 版本最多同时可以建立并保持2000个TCP 连接。
**5. 最大TCB Hash table 数量   MaxHashTableSize TCB 是通过Hash table 来管理的。**
这个值指明分配 pagepool 内存的数量，也就是说，如果MaxFreeTcbs = 1000 , 则 pagepool 的内存数量为 500KB那么 MaxHashTableSize 应大于 500 才行。这个数量越大，则Hash table 的冗余度就越高，每次分配和查找 TCP  连接用时就越少。这个值必须是2的幂，且最大为65536.
 
IBM WebSphere Voice Server 在windows server 2003 下的典型配置
MaxUserPort = 65534 (Decimal)
MaxHashTableSize = 65536 (Decimal)
MaxFreeTcbs = 16000 (Decimal)
这里我们可以看到 MaxHashTableSize 被配置为比MaxFreeTcbs 大4倍，这样可以大大增加TCP建立的速度。 