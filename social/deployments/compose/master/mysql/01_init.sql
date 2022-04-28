CREATE USER 'haproxy_user'@'%' IDENTIFIED BY ''; -- create user for haproxy health check

SET GLOBAL rpl_semi_sync_master_enabled = 1;
SET GLOBAL rpl_semi_sync_master_wait_for_slave_count = 2;
CREATE USER 'repl'@'%' IDENTIFIED BY 'pass';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%';
