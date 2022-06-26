该案例可测试 http to dubbo、nacos、jaeger、jwt

注册中心可以用 nacos 替换 zk

配置文件在 samples/dubbogo/simple/bestdo/pixiu/conf.yml

1. cd 到案例总目录
cd samples/dubbogo/simple/

2. 进行环境准备，启动 注册中心(nacos or zk) 和准备对应配置文件
./start.sh prepare bestdo

3. 启动 dubbo server
./start.sh startServer bestdo

4.  启动 pixiu
./start.sh startPixiu bestdo

5.  启动 Client 测试用例
./start.sh startTest bestdo
