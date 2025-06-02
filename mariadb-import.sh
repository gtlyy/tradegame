#!/bin/bash

# VPS 数据库信息
db_user="tgame"
db_password="tc_3469Uk"
backup_dir="."

# 导入 testdb 数据库
testdb_backup_file="$backup_dir/testdb_backup.sql"
# mysql -u $db_user -p$db_password -e "CREATE DATABASE testdb;"
mysql -u $db_user -p$db_password testdb < $testdb_backup_file

# 导入 webdb 数据库
webdb_backup_file="$backup_dir/webdb_backup.sql"
# mysql -u $db_user -p$db_password -e "CREATE DATABASE webdb;"
mysql -u $db_user -p$db_password webdb < $webdb_backup_file

