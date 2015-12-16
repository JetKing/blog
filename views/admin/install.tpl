<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>安装页面</title>
</head>
<body>
	欢迎安装goblog博客系统
	<form action="/install/start" method="post">
		<label for="dbname">数据库名</label><input type="text" name="dbname" id="dbname" value="test">
		<label for="user">用户名</label><input type="text" name="user" id="user" value="root">
		<label for="passwd">密码</label><input type="password" name="passwd" id="passwd" value="lijun">
		<label for="host">主机</label><input type="text" name="host" id="host" value="127.0.0.1">
		<label for="port">端口</label><input type="text" name="port" id="port" value="3306">
		<input type="submit" value="安装">
	</form>
</body>
</html>