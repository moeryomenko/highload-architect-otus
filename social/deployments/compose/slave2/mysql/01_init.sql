CREATE USER 'haproxy_user'@'%' IDENTIFIED BY ''; -- create user for haproxy health check.

SET GLOBAL rpl_semi_sync_slave_enabled = 1;
CHANGE MASTER TO MASTER_HOST='mysql-master',
  MASTER_USER='repl',
  MASTER_PASSWORD='pass',
  GET_MASTER_PUBLIC_KEY=1,
  MASTER_AUTO_POSITION=1;

START SLAVE;
