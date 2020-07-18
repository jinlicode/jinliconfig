package Template

func DockerComposeNetWorks() string {
	NetWorks := `
networks: 
	discuz:
	  ipam:
		driver: default
		config:
		  - subnet: "10.99.1.0/24"
`
	return NetWorks
}
func DockerComposeCaddy() string {
	Caddy := `
	caddy:
    image: caddy:alpine
    ports:
        - "80:80"
        - "443:443"
    volumes:
        - ./config/caddy/Caddyfile:/etc/caddy/Caddyfile
        - ./config/caddy/cert:/root/caddy
        - ./code:/var/www/html
        - ./log/caddy/:/var/log/
    restart: always
    environment:
        - XDG_DATA_HOME=/root
        - ACME_AGREE=true
        - TZ=Asia/Shanghai
    networks:
      discuz:
        ipv4_address: 10.99.1.2
	`
	return Caddy
}

func DockerComposePhp() string {
	Php := `
  php:
    image: jinlicode/discuz_docker:latest
    user: 10000:10000
    volumes:
        - ./code:/var/www/html
        - ./config/php/php.ini:/usr/local/etc/php/php.ini
        - ./config/php/php-fpm.conf:/usr/local/etc/php-fpm.conf
        - ./config/php/www.conf:/usr/local/etc/php-fpm.d/www.conf
    restart: always
    expose:
        - "9000"
    environment:
        - TZ=Asia/Shanghai
    depends_on:
        - mysql
        - caddy
    networks:
      discuz:
        ipv4_address: 10.99.1.3
`
	return Php
}
func DockerComposeMysql() string {
	Mysql := `
  mysql:
    image: mysql:5.7.30
    restart: always
    container_name: mysql
    volumes: 
      - ./db:/var/lib/mysql
      - ./imput_db:/docker-entrypoint-initdb.d
      - ./config/mysql/my.cnf:/etc/mysql/my.cnf
    expose:
      - "3306"
    environment:
      MYSQL_ROOT_PASSWORD: discuz
      MYSQL_DATABASE: discuz
      TZ: Asia/Shanghai
    networks:
      discuz:
        ipv4_address: 10.99.1.4
  `
	return Mysql
}
func DockerComposeMemcached() string {
	Memcached := `
  memcached:
    image: bitnami/memcached:1.6.6
    restart: always
    container_name: memcached
    environment:
      - MEMCACHED_CACHE_SIZE=64
    networks:
      discuz:
        ipv4_address: 10.99.1.5
  `
	return Memcached
}
func DockerComposeRedis() string {

}
