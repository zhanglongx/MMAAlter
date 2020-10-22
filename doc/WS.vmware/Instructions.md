# WS VMware

## 步骤

1. 安装VMWare

安装说明(https://www.tecmint.com/install-vmware-workstation-in-linux/)

2. 增加NAT端口映射

在/etc/vmware/vmnet8/nat/nat.conf下增加以下字段，参照(https://docs.vmware.com/en/VMware-Workstation-Pro/12.0/com.vmware.ws.using.doc/GUID-C2EC7B92-A499-4B47-95B6-0BFDDA28AC34.html)

注意替换以下文本中的IP地址
 
	[incomingtcp]

	# Use these with care - anyone can enter into your VM through these...
	# The format and example are as follows:
	#<external port number> = <VM's IP address>:<VM's port number>
	#8080 = 172.16.3.128:80
	# 60000 = 192.168.64.135:60000
	5600 = 192.168.64.135:5600
	5602 = 192.168.64.135:5602
	40002 = 192.168.64.135:40002
	40005 = 192.168.64.135:40005
	9000 = 192.168.64.135:9000
	3306 = 192.168.64.135:3306
	1901 = 192.168.64.135:1901
	# 111 = 192.168.64.135:111
	4369 = 192.168.64.135:4369
	4370 = 192.168.64.135:4370
	# 786 = 192.168.64.135:768
	12346 = 192.168.64.135:12346
	10013 = 192.168.64.135:10013
	40000 = 192.168.64.135:40000
	7001 = 192.168.64.135:7001
	7002 = 192.168.64.135:7002
	5050 = 192.168.64.135:5050
	
	[incomingudp]
	
	# UDP port forwarding example
	#6000 = 172.16.3.0:6001

3. 创建systemd.unit

参照
(https://www.atrixnet.com/autostart-vmware-virtual-machines-on-boot-in-linux/)

使用以下两个文件：

vmware-autostarts   - 虚拟机启动脚本

vmware-vmx.service  - systemd units

4. （可选，安装MMAAlter可能需要）安装配置transit

	go build github.com/zhanglongx/transit

修改/usr/local/etc/transit.json

	{
	  "IPArray": [
	    "11.11.11.104",
	    "11.11.11.106"
	  ],
	  "ThirdPartyAddr": "11.11.11.104:8001",
	  "IP": "11.11.11.109",
	  "Port": 7001,
	  "Pattern": "(serverip=)'\\d+\\.\\d+\\.\\d+\\.\\d+'",
	  "Replace": "$1'11.11.11.109'"
	}
