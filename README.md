## A Simple Go-Web Project

一个简单的前后端分离的Go Web项目，实现注册、登陆的基本功能。前端使用基本的HTML；后端使用基于Golang的Gin框架、captcha库、gorm库、go-redis等。项目使用MVC架构，将控制器、模型及视图分开，有助于理解整个项目。



***项目主要实现功能：***

**注册：**

- 使用【用户名、密码、邮箱和手机号】注册；
- 用户名、邮箱和密码验证唯一；
- 密码包含复杂度验证，至少包含【大写字母、小写字母、数字和特殊字符./*】，并使用MD5加密存储；

**登陆：**

- 用户名、密码输入正确可直接进入首页；
- 密码输错5次，登陆页面刷新出验证码；
- 对于不同的浏览器，用sessionID + Redis 记录错误登陆次数；
- 登陆成功，使用jwt生成token，并保存在session中，在home page进行jwt中间件鉴权。



***使用方法：***

- 在MySQL创建Database - “myWeb”，在`dao/mysql.go/InitMySQL()`中修改MySQL密码；
- 启动一个无密码的Redis；

- 运行项目，默认端口地址：`localhost:8080`，注册并登陆。



***效果展示：***

1. **首页**

<img src="https://github.com/M3stark/GoWeb/blob/main/source/index.png" alt="image-20220905013355568" style="zoom:50%;" />



2. **注册**

<img src="https://github.com/M3stark/GoWeb/blob/main/source/register.png" alt="image-20220905013557132" style="zoom:50%;" />

用户名必须唯一：

<img src="https://github.com/M3stark/GoWeb/blob/main/source/register_1.png" alt="image-20220905013653294" style="zoom:50%;" />

密码缺少数字：

<img src="https://github.com/M3stark/GoWeb/blob/main/source/register_2.png" alt="image-20220905013723641" style="zoom:50%;" />

3. **登陆**

<img src="https://github.com/M3stark/GoWeb/blob/main/source/login.png" alt="image-20220905015824081" style="zoom:50%;" />

不同浏览器登陆失败，将不同浏览器的session ID作为key，保存错误登陆次数：

![image-20220905020130559](https://github.com/M3stark/GoWeb/blob/main/source/login_1.png)

登陆错5次，加载验证码：

<img src="https://github.com/M3stark/GoWeb/blob/main/source/login_2.png" alt="image-20220905020210304" style="zoom:50%;" />

登陆成功，返回token：

<img src="https://github.com/M3stark/GoWeb/blob/main/source/login_3.png" alt="image-20220905020255126" style="zoom:50%;" />

4. **home page**

登陆成功后，访问home page，解析浏览器缓存的token，得到用户名，打印到home page，实现jwt鉴权：

<img src="https://github.com/M3stark/GoWeb/blob/main/source/home.png" alt="image-20220905020425558" style="zoom:50%;" />
