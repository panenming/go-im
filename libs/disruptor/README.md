## 高性能并发框架(异步)
## 简单原理
	关键在于不需要同步锁，使用ringbuffer
	给一个最简单的disruptor，single consumer single producer的实现，不需要同步锁：
	共享的变量：
	int volatile readBarrier;
	Data[] ring = new Data[A_BIG_NUMBER];
	
	Producer 线程:
	int writeCount;
	void produce(Data newItem) {  
		ring[++writeCount % ring.size() ] = newItem;  
		readBarrier++;
	}		
	
	Consumer 线程:
	void run() {  
		int currentRead;  
		while (true) {    
			if (readBarrier > currentRead) {        
				consume (ring [ currentRead % ring.size() ];        				currentRead++;    
			}  
		}
	}
	
	关键点在于：Consumer线程是只读的，所以理论上并不需要锁的参与，只要控制好readBarrier增量的时机，而Consumer线程只要一直轮询这个变量即可。
	
## 核心原理之1－－单线程写

	Linux内核的kfifo队列，也是RingBuffer。其关键特征是“一读一写“，因此可以“完全无锁“，连CAS都不需要。
	
	同样，Disruptor的RingBuffer，至所以可以做到完全无锁，也是因为“单线程写“，这是所有“前提的前提“。离了这个前提条件，没有任何技术可以做到完全无锁。

## 核心原理之2 －－ 内存屏障

	除了上面的“单写“这个前提条件，要正确的实现无锁，还需要另外一个关键技术：内存屏障。
	
	对应到Java语言，就是valotile变量与happen before语义。
	
	下面将对内存屏障与Java的volatile变量之间的关系做一个梳理。
	内存屏障 – Linux的smp_wmb()/smp_rmb()
	
	内存屏障其实就是1条cpu指令，这条cpu指令有以下2个作用：
	（1）阻止指令的重排序。插入1个内存屏障之后，屏障之后的代码，不会被重排序到屏障之前。
	（2）flush store缓存/load缓存
	
	具体来讲，有2种内存屏障：
	store barrier（写屏障）：刷新store缓存。即把store barrier 之前的写操作，也就是store缓存里面的内容，刷新到主存，从而其它cpu可以看到写的值；
	load barrier（读屏障）：失效load缓存。从而使得load barrier之后的读操作，不会读到store缓存里面的旧值，而会直接读到其他cpu更新后的新值。
	
	full barrier: 即同时具有store barrier + load barrier的功能

## 核心原理之3 －－ 伪共享与缓存行填充
	
	“缓存行填充“，简单讲就是不要让2个变量分配到1个cache line上面，这样会造成1个变量修改，整个cache line失效，另一个变量的缓存也失效。
	
	假如有两个变量x, y，一个线程修改x，另一个线程读y，看起来是相互独立的。但如果x, y处在同一个cache line里面，这会导致一个thread对x的修改，影响另一个thread对y的读取性能，也就是2个线程对同1个cache line发生竞争，这也称之为“伪共享“。
	
	典型的，比如1个链表有头部、尾部2个指针，如果这2个指针分配到了同一个cache line上面，当你不断修改头部指针的时候，尾部指针的缓存也受影响。
	
	通常来讲，cache line是64Byte，比如我想让2个long型变量不要分配到同1个cache line上，就可以为其中每个变量填充7个long型。
