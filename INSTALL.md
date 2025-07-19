# Install (Linux & Windows)

## 1. MariaDB
### Install MariaDB
```
sudo dnf install mariadb-server
```

### Start server
```
sudo systemctl start mariadb
sudo systemctl enable mariadb
```

### Create user
```
# 注意禁止root远程登陆 root/rF111222k
sudo mysql_secure_installation
sudo mysql -u root -p
CREATE OR REPLACE USER tgame@localhost IDENTIFIED BY 'tc_3469Uk';
```

### Create db
```
CREATE DATABASE testdb;
GRANT ALL ON testdb.* to 'tgame'@'localhost';
CREATE DATABASE webdb;
GRANT ALL ON webdb.* to 'tgame'@'localhost';
```

## 2. Import to mariadb
`sh mariadb-import.sh`


## 3. RabbitMQ
### Install RabbitMQ
```
# sudo dnf install rabbitmq-server 
# Fedora install :
sudo dnf install ~/download/rabbitmq-server-3.13.2-1.el8.noarch.rpm
# https://www.rabbitmq.com/install-rpm.html#downloads
# ps: Fedora这样装： sudo dnf install rabbitmq-server ，有时启动不了。
```

### Enable Plugins
```
# 这样通过 js 可以访问 web_stomp这样
sudo rabbitmq-plugins enable rabbitmq_management rabbitmq_web_stomp rabbitmq_web_stomp_examples
# add to /etc/rabbitmq/rabbitmq.conf:  web_stomp.tcp.port = 15674
# Ps： 还可以通过web访问：http://localhost:15672
```

### Start server
```
sudo systemctl start rabbitmq-server
sudo systemctl enable rabbitmq-server
```

### Add user
```
sudo rabbitmqctl add_user userdw01 pq328hu7
# sudo rabbitmqctl change_password userdw01 pq328hu7
sudo rabbitmqctl set_user_tags userdw01 administrator
sudo rabbitmqctl  set_permissions -p / userdw01 '.*' '.*' '.*'
```

### 开放端口访问
```
# 获取firewalld当前使用的区域名字
firewallZone=`sudo firewall-cmd --list-all | grep active | cut -d\( -f1`
echo $firewallZone
sudo firewall-cmd --permanent --zone=$firewallZone --add-port=15674/tcp
sudo firewall-cmd --reload
```

## 4. Run
bee run

# End Install.
