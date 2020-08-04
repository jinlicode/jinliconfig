package Template

func TemplateCaddyHttps() string {
	TemplateCaddyFileHttps := `
server {
    listen                  443 ssl http2;
    server_name             www.example.com;
    set                     $base /var/www/www.example.com;
    root                    $base;

    # 这里配上之前我们生成的自签名证书，否则会报错
    ssl_certificate /etc/ssl/default.crt; # managed by Certbot
    ssl_certificate_key /etc/ssl/default.crt; # managed by Certbot

    # security
    include                 /etc/nginx/jinli_nginx_base_config/security.conf;

    # index.php
    index                   index.php;

    # rewrite

    # additional config
    include /etc/nginx/jinli_nginx_base_config/general.conf;

    # handle .php
    location ~ \.php$ {
        fastcgi_pass                  php:9000;
        include /etc/nginx/jinli_nginx_base_config/php_fastcgi.conf;
    }
}

# HTTP redirect
server {
    if ($host = www.example.com) {
        return 301 https://$host$request_uri;
    }

    listen      80;
    server_name www.example.com;
    root                    $base;
}
`

	return TemplateCaddyFileHttps

}

func TemplateCaddyHttp() string {
	TemplateCaddyFileHttp := `
server {
    listen                  80;
    server_name             www.example.com;
    set                     $base /var/www/www.example.com;
    root                    $base;

    # security
    include                 /etc/nginx/jinli_nginx_base_config/security.conf;

    # index.php
    index                   index.php;

    # rewrite

    # additional config
    include /etc/nginx/jinli_nginx_base_config/general.conf;

    # handle .php
    location ~ \.php$ {
        fastcgi_pass                  php:9000;
        include /etc/nginx/jinli_nginx_base_config/php_fastcgi.conf;
    }
}
`

	return TemplateCaddyFileHttp

}

func TemplateCaddyRewriteThinkphp() string {
	RewriteThinkphp := `
	location / {
	try_files $uri $uri/ /index.php$uri;
    }
	`
	return RewriteThinkphp
}

func TemplateCaddyRewriteDiscuz() string {
	RewriteDiscuz := `
	location /{
		rewrite ^([^\.]*)/topic-(.+)\.html$ $1/portal.php?mod=topic&topic=$2 last;
		rewrite ^([^\.]*)/article-([0-9]+)-([0-9]+)\.html$ $1/portal.php?mod=view&aid=$2&page=$3 last;
		rewrite ^([^\.]*)/forum-(\w+)-([0-9]+)\.html$ $1/forum.php?mod=forumdisplay&fid=$2&page=$3 last;
		rewrite ^([^\.]*)/thread-([0-9]+)-([0-9]+)-([0-9]+)\.html$ $1/forum.php?mod=viewthread&tid=$2&extra=page%3D$4&page=$3 last;
		rewrite ^([^\.]*)/group-([0-9]+)-([0-9]+)\.html$ $1/forum.php?mod=group&fid=$2&page=$3 last;
		rewrite ^([^\.]*)/space-(username|uid)-(.+)\.html$ $1/home.php?mod=space&$2=$3 last;
		rewrite ^([^\.]*)/blog-([0-9]+)-([0-9]+)\.html$ $1/home.php?mod=space&uid=$2&do=blog&id=$3 last;
		rewrite ^([^\.]*)/(fid|tid)-([0-9]+)\.html$ $1/index.php?action=$2&value=$3 last;
		rewrite ^([^\.]*)/([a-z]+[a-z0-9_]*)-([a-z0-9_\-]+)\.html$ $1/plugin.php?id=$2:$3 last;
		if (!-e $request_filename) {    return 404;}
	}
	`
	return RewriteDiscuz
}
