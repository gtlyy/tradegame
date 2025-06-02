#!/bin/bash

# 客户端数据库信息
db_host="localhost"
db_user="tgame"
db_password="tc_3469Uk"
backup_dir="."

# 导出 testdb 数据库
testdb_backup_file="$backup_dir/testdb_backup.sql"
mysqldump -h $db_host -u $db_user -p$db_password testdb > $testdb_backup_file

# 导出 webdb 数据库
webdb_backup_file="$backup_dir/webdb_backup.sql"
mysqldump -h $db_host -u $db_user -p$db_password webdb > $webdb_backup_file

# 将备份文件传输到 VPS
# scp $testdb_backup_file $webdb_backup_file username@vps_ip:/path/on/vps/
