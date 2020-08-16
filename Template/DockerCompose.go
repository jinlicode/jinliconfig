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
func DockerComposeNginx() string {
	Nginx := `
  nginx:
    image: jinlicode/nginx:v1
    ports:
        - "80:80"
        - "443:443"
    volumes:
        - ./config/nginx/conf/:/etc/nginx/conf.d/
        - ./code:/var/www/test1.jinli.plus
        - ./log/nginx/:/var/log/nginx
        - ./config/cert/:/etc/letsencrypt/
    restart: always
    environment:
        # - KEYSIZE="4096"
        # - KEY_ALGO="rsa"
        - CONTACT_EMAIL="maniac.cn@gmail.com"
        - TZ=Asia/Shanghai
        #- XDG_DATA_HOME=/root
    networks:
      discuz:
        ipv4_address: 10.99.1.2
	`
	return Nginx
}

func DockerComposePhp() string {
	Php := `
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
        - nginx
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
