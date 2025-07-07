# Performance indicator collection
Go语言编写的计算机资源指标的采集工具，集成prometheus构建完整的监控系统 
采集图表如下所示:
#### CPU&MEM
![image](https://github.com/user-attachments/assets/88c4eccf-f482-4982-9720-034309d9d254)

## 数据采集解析
### 系统级数据
| 资源指标      | 数据项        | 数据来源| 描述                         |单位 | 告警限度|
|--------------|---------------|-----------|----------------------|:-----:|------|
| **CPU**       | CPU 平均负载     | /proc/loadavg     | 采集最近1、5、15分钟的平均负载   |进程数|<= 核数*0.7|
|               | 用户态           | /proc/stat#user   | 提供1s和15s的平均使用率, 下同   |%| <= 80%|
|               | 用户态低优先级   | /proc/stat#nice   | /                              |%|        /     |
|               | 内核态           | /proc/stat#system | /                             |%|        /     |
|               | 空闲等待         | /proc/stat#idle   | /                             |%|        /     |
|               | 等待IO读写       | /proc/stat#iowait | /                             |%|        /     |
|               | 硬中断           | /proc/stat#irq| /                             |%|        /     |
|               | 软中断           | /proc/stat#softirq| /                             |%|        /     |
|               | 宿主机资源竞争   | /proc/stat#steal| /                             |%|        /     |
| **内存**      | 总内存           | /proc/meminfo#memtotal | 总内存量                  | MB | / |
|               | 已使用的物理内存  | /proc/meminfo#memused | total-free-buffer-cached  | MB | / |
|               | 内存使用率       | /proc/meminfo#memused% | used / total             | % | <= 70% |
|               | 空闲内存         | /proc/meminfo#memfree  | 尚未分配的物理内存         | MB | / |
|               | 可获得内存       | /proc/meminfo#memavail | 包括未分配和可回收         | MB | / |
|               | 元数据缓存       | /proc/meminfo#Buffers  | 文件目录结构、权限信息等    | MB | / |
|               | 文件数据缓存     | /proc/meminfo#Cached   | 文件数据内容的缓存         | MB | / |
|               | 活跃内存         | /proc/meminfo#Active   | 被标记的常被访问的内存     | MB | / |
|               | 不活跃内存       | /proc/meminfo#Inactive | 不经常使用的内存          | MB | / |
|               | 总交换空间       | /proc/meminfo#SwapTotal | 交换空间的总大小          | MB| / |
|               | 交换分区缓存     | /proc/meminfo#SwapCached | 交换空间中的数据被swap in读到缓冲区中占用的大小 | MB | / |
|               | 空闲的交换空间   | /proc/meminfo#SwapFree  | 交换空间空闲的大小         | MB | / |
|               | 共享内存         | /proc/meminfo#Shmem    | 已分配的共享内存的大小     | MB | / |
| **上下文切换**     | 上下文切换次数           | 进程切换上下文的次数                                                 |
| **中断次数**       | 中断发生次数             | 系统发生中断的次数                                                   |
| **队列长度**       | 磁盘I/O请求队列长度      | 磁盘I/O请求的平均队列长度                                            |
| **Swap/IO**        | Swap空间使用情况         | Swap空间的使用情况                                                   |
|                    | 磁盘读写速度             | 磁盘读写数据的速度                                                   |
|                    | 平均服务时间             | 完成磁盘I/O请求的平均时间                                            |
| **网络**           | 网络连接状态             | 当前系统的网络连接状态                                               |
|                    | 网络带宽                 | 网络传输数据的速度                                                   |
|                    | 吞吐量                   | 网络传输数据的吞吐量                                                 |
|                    | 延迟                     | 数据包从发送方到接收方的传输时间                                     |
|                    | 错误和丢包率             | 网络传输中的错误和丢失数据包的比例                                   |
| **文件存储**       | 文件系统使用情况         | 文件系统的使用情况和剩余空间                                         |
| **磁盘性能**       | 磁盘I/O性能              | 磁盘的输入输出性能                                                   |


### /proc/loadavg 文件

```
[root@localhost ~]# cat /proc/loadavg 
0.00 0.01 0.05 1/561 27848
```


| 列| 含义| 值|
|-----|----|---|
|1| 最近1分钟的平均负载| 0.00|
|2| 最近5分钟的平均负载| 0.01|
|3| 最近15分钟的平均负载| 0.05|
|4| 正在运行的进程数 / 总进程数| 1/561|
|5| 最近运行的进程pid| 27848|

### /proc/stat 文件

```
[root@localhost ~]# cat /proc/stat
cpu  77395599 98038 7722431 2581562838 210211 0 53382 0 0 0
cpu0 22174837 30164 2005109 641862586 24046 0 45986 0 0 0
cpu1 37848314 32905 1790830 627385985 29808 0 3365 0 0 0
cpu2 10933156 17247 1977764 653984058 31883 0 1080 0 0 0
cpu3 6439291 17720 1948727 658330207 124473 0 2950 0 0 0
intr 2372730827 38 5 0 0 0 0 0 0 1 5 0 0 6 0 0 0 0 27 0 0 0 0 0 29 0 0 0 263836738 4691392 29 13 150 779 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 2853775527
btime 1738811073
processes 12142158
procs_running 1
procs_blocked 0
softirq 2159465392 6 1239680310 1389451 263840897 4685200 0 1965140 291064155 0 356840233
```


##### [1-5] 单行数据解析
```
cpu  77395599 98038 7722431 2581562838 210211 0 53382 0 0 0
cpu0 22174837 30164 2005109 641862586 24046 0 45986 0 0 0
cpu1 37848314 32905 1790830 627385985 29808 0 3365 0 0 0
cpu2 10933156 17247 1977764 653984058 31883 0 1080 0 0 0
cpu3 6439291 17720 1948727 658330207 124473 0 2950 0 0 0
```

| 列| 列名 |含义| 值(首行)|
|-----|:----:|----|---:|
|1| CPU序列号| / |cpu (总)
|2| user| 用户态非低优先级消耗的CPU总时间| 77395599 |
|3| nice| 用户态低优先级消耗的CPU总时间| 98038|
|4| system| 系统内核消耗的CPU总时间| 7722431|
|5| idle| CPU空闲时间, 不包括磁盘IO等待时间| 2581562838|
|6| iowait| CPU等待磁盘IO的时间| 210211|
|7| irp| 消耗在被硬中断切出到其他任务的时间| 0|
|8| softirq| 消耗在被软中断切出到其他任务的时间| 53382|
|9| steal| 当前虚拟机因宿主机资源竞争被迫等待的时间(宿主机为0)| 0|
|10|guest| 运行虚拟机程序消耗的时间(虚拟机为0)| 0|
|11|guest_nice| 运行低优先级虚拟机程序消耗的时间(虚拟机为0)| 0|

##### [6] 单行数据解析
```
intr 2372730827 38 5 0 0 0 0 0 0 1 5 0 0 6 0 0 0 0 27 0 0 0 0 0 29 0 0 0 263836738 4691392 29 13 150 779 0 0 0*n
```

| 列| 列名 |含义| 值|
|-----|:----:|----|---:|
|1|系统中断次数| 自系统运行以来, 所有类型中断的总和| 2372730827|
|2-n| 中断次数| 不同类型的中断次数统计| / |


##### [7-11] 数据解析
```
ctxt 2853775527
btime 1738811073
processes 12142158
procs_running 1
procs_blocked 0
```

| 名称 |含义| 值|
|----|----|---:|
| ctxt | 上下文交换次数| 2853775527|
| btime | 系统启动时间| 1738811073 |
| processes | 总创建进程数统计| 12142158 |
| procs_running | 正在运行的进程数 | 1 |
| procs_blocked | 被阻塞运行的进程数| 0 |

##### [12] 单行数据解析
```
softirq 2159465392 6 1239680310 1389451 263840897 4685200 0 1965140 291064155 0 356840233
```

| 列| 列名 |含义| 值|
|-----|:----:|----|---:|
|1|系统软中断次数| 自系统运行以来, 所有软中断不同类型的总和| 2372730827|
|2-n| 中断次数| 不同类型的软中断次数统计| / |


### /proc/meminfo 文件
```
[root@localhost ~]# cat /proc/meminfo
MemTotal:        7920380 kB
MemFree:         1142776 kB
MemAvailable:    5719472 kB
Buffers:               0 kB
Cached:          5152920 kB
SwapCached:            0 kB
Active:          3252780 kB
Inactive:        3083752 kB
Active(anon):     706772 kB
Inactive(anon):   923196 kB
Active(file):    2546008 kB
Inactive(file):  2160556 kB
Unevictable:         644 kB
Mlocked:               0 kB
SwapTotal:       8126460 kB
SwapFree:        8124824 kB
Dirty:                64 kB
Writeback:             0 kB
AnonPages:       1184272 kB
Mapped:            88184 kB
Shmem:            446356 kB
Slab:             234856 kB
SReclaimable:     173964 kB
SUnreclaim:        60892 kB
KernelStack:        9088 kB
PageTables:        37560 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:    12086648 kB
Committed_AS:    4608288 kB
VmallocTotal:   34359738367 kB
VmallocUsed:      356308 kB
VmallocChunk:   34358947836 kB
HardwareCorrupted:     0 kB
AnonHugePages:    704512 kB
CmaTotal:              0 kB
CmaFree:               0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:      191776 kB
DirectMap2M:     4962304 kB
DirectMap1G:     3145728 kB
```

| 名称 |含义| 值|
|----|----|---|
| MemTotal | 系统总内存(去除预留位和二进制代码)| 7920380 kB |
| MemFree | 系统尚未分配的物理内存| 1142776 kB |
| MemAvailable | 系统可获得的物理内存, 包括未分配和可以回收的buffer/cache内存| 5719472 kB |
| Buffers | 元数据的缓存(文件目录结构、权限信息)| 0 kB|
| Cached | 文件内容的缓存(高速缓存) |         5152920 kB |
| SwapCached | 交换空间中的数据被swap in读到缓冲区中占用的大小(此时源数据仍保留在交换空间中) | 0 kB |
| Active | 经常使用的缓存Buffers&Cached的大小 | 3252780 kB |
| Inactive | 不经常使用的缓存Buffers&Cached的大小| 3083752 kB|
| SwapTotal | 交换空间的总大小 | 8126460 kB |
| SwapFree | 交换空间空闲的大小 | 8124824 kB |
| Shmem | 已分配的共享内存的大小 | 446356 kB |

<!-- 
Active(anon):     706772 kB
Inactive(anon):   923196 kB
Active(file):    2546008 kB
Inactive(file):  2160556 kB
Unevictable:         644 kB
Mlocked:               0 kB
Dirty:                64 kB
Writeback:             0 kB
AnonPages:       1184272 kB
Mapped:            88184 kB
Slab:             234856 kB
SReclaimable:     173964 kB
SUnreclaim:        60892 kB
KernelStack:        9088 kB
PageTables:        37560 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:    12086648 kB
Committed_AS:    4608288 kB
VmallocTotal:   34359738367 kB
VmallocUsed:      356308 kB
VmallocChunk:   34358947836 kB
HardwareCorrupted:     0 kB
AnonHugePages:    704512 kB
CmaTotal:              0 kB
CmaFree:               0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:      191776 kB
DirectMap2M:     4962304 kB
DirectMap1G:     3145728 kB -->
