package Template

//HTMLIndex 创建成功之后的默认index文件内容
func HTMLIndex() string {
	HTML := `
<!doctype html>
<html>
<head>
	<meta charset="utf-8">
	<title>恭喜，站点创建成功！</title>
	<style>
		.container {
			width: 60%;
			margin: 10% auto 0;
			background-color: #f0f0f0;
			padding: 2% 5%;
			border-radius: 10px
		}

		ul {
			padding-left: 20px;
		}

			ul li {
				line-height: 2.3
			}

		a {
			color: #20a53a
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>恭喜, 站点创建成功！</h1>
		<h3>这是默认index.html，本页面由系统自动生成</h3>
		<ul>
			<li>本页面在根目录下的index.html</li>
			<li>您可以修改、删除或覆盖本页面</li>
		</ul>
	</div>
</body>
</html>
`
	return HTML
}

//HTML404 创建成功之后的默认index文件内容
func HTML404() string {
	HTML := `
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
<title>404</title>
<style>
	body{
		background-color:#444;
		font-size:14px;
	}
	h3{
		font-size:60px;
		color:#eee;
		text-align:center;
		padding-top:30px;
		font-weight:normal;
	}
</style>
</head>

<body>
<h3>404，您请求的文件不存在!</h3>
</body>
</html>
`
	return HTML

}
