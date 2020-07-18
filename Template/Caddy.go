package Template

func TemplateCaddyHttps() string {
	TemplateCaddyFileHttps := `
vip.jinli.plus
{
	tls maniac.cn@gmail.com
	encode gzip
	file_server * browse
	root * /var/www/html/
	php_fastcgi php:9000
	log {
		output file /var/log/access.log {
			roll_size 10gb
			roll_keep 365
			roll_keep_for 720h
		}
	}

}`

	return TemplateCaddyFileHttps

}

func TemplateCaddyHttp() string {
	TemplateCaddyFileHttp := `
http://http://vip.jinli.plus
{
	tls maniac.cn@gmail.com
	encode gzip
	file_server * browse
	root * /var/www/html/
	php_fastcgi php:9000
	log {
		output file /var/log/access.log {
			roll_size 10gb
			roll_keep 365
			roll_keep_for 720h
		}
	}

}`

	return TemplateCaddyFileHttp

}

func TemplateCaddyRewriteThinkphp() string {
	RewriteThinkphp := `
	@key0 {
		not file 
		path_regexp key0 ^(.*)$ 
	}
	rewrite @key0 /index.php?s={re.key0.1}
	`
	return RewriteThinkphp
}

func TemplateCaddyRewriteDiscuz() string {
	RewriteDiscuz := `
	@key0 {
	not file 
	path_regexp key0 ([\.]*)/topic-(.+)\.html 
	}
	rewrite @key0 {re.key0.1}/portal.php?mod=topic&topic={re.key0.2}
	@key1 {
	not file 
	path_regexp key1 ([\.]*)/article-([0-9]+)-([0-9]+)\.html 
	}
	rewrite @key1 {re.key1.1}/portal.php?mod=view&aid={re.key1.2}&page={re.key1.3}
	@key2 {
	not file 
	path_regexp key2 ([\.]*)/forum-(\w+)-([0-9]+)\.html 
	}
	rewrite @key2 {re.key2.1}/forum.php?mod=forumdisplay&fid={re.key2.2}&page={re.key2.3}
	@key3 {
	not file 
	path_regexp key3 ([\.]*)/thread-([0-9]+)-([0-9]+)-([0-9]+)\.html 
	}
	rewrite @key3 {re.key3.1}/forum.php?mod=viewthread&tid={re.key3.2}&extra=page%3D{re.key3.4}&page={re.key3.3}
	@key4 {
	not file 
	path_regexp key4 ([\.]*)/group-([0-9]+)-([0-9]+)\.html 
	}
	rewrite @key4 {re.key4.1}/forum.php?mod=group&fid={re.key4.2}&page={re.key4.3}
	@key5 {
	not file 
	path_regexp key5 ([\.]*)/space-(username|uid)-(.+)\.html 
	}
	rewrite @key5 {re.key5.1}/home.php?mod=space&{re.key5.2}={re.key5.3}
	@key6 {
	not file 
	path_regexp key6 ([\.]*)/blog-([0-9]+)-([0-9]+)\.html 
	}
	rewrite @key6 {re.key6.1}/home.php?mod=space&uid={re.key6.2}&do=blog&id={re.key6.3}
	@key7 {
	not file 
	path_regexp key7 ([\.]*)/(fid|tid)-([0-9]+)\.html 
	}
	rewrite @key7 {re.key7.1}/index.php?action={re.key7.2}&value={re.key7.3}
	@key8 {
	not file 
	path_regexp key8 ([\.]*)/([a-z]+[a-z0-9_]*)-([a-z0-9_\-]+)\.html 
	}
	rewrite @key8 {re.key8.1}/plugin.php?id={re.key8.2}:{re.key8.3}
	`
	return RewriteDiscuz
}
