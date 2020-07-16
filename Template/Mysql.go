package Template

func TemplateMysql() string {
	TemplateMysqlCnf := `
[client]
default-character-set = utf8mb4
[mysql]
default-character-set = utf8mb4
[mysqld]
skip-host-cache
skip-name-resolve
datadir = /var/lib/mysql
external-locking = FALSE
character-set-client-handshake = FALSE
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
binlog-cache-size = 32M
sync-binlog = 1
max_binlog_cache_size = 1G
max_binlog_size = 1G
expire_logs_days = 90
max_connections = 2000
max_user_connections = 2000
max_connect_errors = 100000
open_files_limit = 1024
max_allowed_packet = 8M
wait_timeout = 360
interactive_timeout = 360
slow_query_log = on
log-queries-not-using-indexes
long_query_time = 2
min_examined_row_limit = 5
query_cache_size = 128M
query_cache_limit = 512k
query_cache_min_res_unit = 1k
key_buffer_size = 128M
sort_buffer_size = 16M
join_buffer_size = 8M
read_buffer_size = 16M
read_rnd_buffer_size = 32M
thread_stack = 192K
bulk_insert_buffer_size = 16M
myisam_sort_buffer_size = 128M
myisam_max_sort_file_size = 512M
myisam_repair_threads = 1
key_buffer_size = 4M
transaction_isolation = READ-COMMITTED
tmp_table_size = 16M
max_heap_table_size = 16M
default-storage-engine=InnoDB
innodb_old_blocks_time =1000
innodb_flush_method = O_DIRECT
innodb_buffer_pool_size = 512M
innodb_thread_concurrency = 8
innodb_flush_log_at_trx_commit = 2
innodb_log_buffer_size = 2M
innodb_log_file_size = 64M
innodb_log_files_in_group = 3
innodb_max_dirty_pages_pct = 90
innodb_lock_wait_timeout = 120
innodb_file_per_table = 1
innodb_autoextend_increment = 256
[mysqldump]
quick
max_allowed_packet = 64M
!includedir /etc/mysql/conf.d/
`
	return TemplateMysqlCnf
}
