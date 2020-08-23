package Template

func DockerComposeVersion() string {
	Version := `3.1`
	return Version
}

func DockerComposeNetWorks() string {
	NetWorks := `
jinli_net:
  ipam:
    driver: default
    config:
      - subnet: "10.99.1.0/16"
`
	return NetWorks
}
func DockerComposeNginx() string {
	Nginx := `
    image: registry.cn-beijing.aliyuncs.com/jinlicode/nginx:v1
    container_name: nginx
    ports:
        - "80:80"
        - "443:443"
    volumes:
        - ./config/nginx/:/etc/nginx/conf.d/
        - ./code:/var/www
        - ./log/nginx/:/var/log/nginx
        - ./config/cert/:/etc/letsencrypt/
    restart: always
    environment:
        - TZ=Asia/Shanghai
        - JINLIVER=1.1
    networks:
      jinli_net:
        ipv4_address: 10.99.1.2
`
	return Nginx
}

func DockerComposePhp() string {
	Php := `
    image: registry.cn-beijing.aliyuncs.com/jinlicode/php:latest
    user: 10000:10000
    container_name: www_example_com
    volumes:
        - ./code/www_example_com:/var/www/www_example_com
        - ./config/php/www_example_com/php.ini:/usr/local/etc/php/php.ini
        - ./config/php/www_example_com/php-fpm.conf:/usr/local/etc/php-fpm.conf
        - ./config/php/www_example_com/www.conf:/usr/local/etc/php-fpm.d/www.conf
        - ./log/openrasp/www_example_com:/opt/rasp/logs/alarm
    restart: always
    expose:
        - "9000"
    environment:
        - TZ=Asia/Shanghai
    depends_on:
        - mysql
        - nginx
    networks:
      jinli_net:
        ipv4_address: 10.99.2.2
`
	return Php
}
func DockerComposeMysql() string {
	Mysql := `
    image: registry.cn-beijing.aliyuncs.com/jinlicode/mysql
    restart: always
    container_name: mysql
    volumes: 
      - ./db:/var/lib/mysql
      - ./imput_db:/docker-entrypoint-initdb.d
      - ./config/mysql/my.cnf:/etc/mysql/my.cnf
    expose:
      - "3306"
    environment:
      MYSQL_ROOT_PASSWORD: root
      TZ: Asia/Shanghai
    networks:
      jinli_net:
        ipv4_address: 10.99.3.2
  `
	return Mysql
}
func DockerComposeMemcached() string {
	Memcached := `
  memcached:
    image: registry.cn-beijing.aliyuncs.com/jinlicode/memcached:1.6.6
    restart: always
    container_name: memcached
    expose:
    - "11211"
    environment:
      - TZ=Asia/Shanghai
      - MEMCACHED_CACHE_SIZE=64
    networks:
      jinli_net:
        ipv4_address: 10.99.4.2
  `
	return Memcached
}
func DockerComposeRedis() string {
	Redis := `
  redis:
    image: registry.cn-beijing.aliyuncs.com/jinlicode/redis:5.0.9
    restart: always
    container_name: redis
    expose:
    - "6379"
    environment:
      - TZ=Asia/Shanghai
    networks:
      jinli_net:
        ipv4_address: 10.99.5.2
  `
	return Redis
}
func DockerComposePhpmyadmin() string {
	Phpmyadmin := `
  phpmyadmin:
    image: registry.cn-beijing.aliyuncs.com/jinlicode/phpmyadmin:5.0.2:
    restart: always
    container_name: phpmyadmin
    ports:
      - "8080:80"
    environment:
      - TZ=Asia/Shanghai
    networks:
      jinli_net:
        ipv4_address: 10.99.6.2
  `
	return Phpmyadmin
}
