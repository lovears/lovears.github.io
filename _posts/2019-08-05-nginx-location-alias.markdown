---
layout:     post 
title:      "【Nginx】Nginx proxy_pass,root,alias"
subtitle:   "nginx 代理的几种方式： alias,proxy_pass,root等"
date:       2019-08-05
author:     "老回路"
header-img: "img/bg-todo.png"
tags:
    - nginx
    - location
    - proxy_pass
---

nginx的使用，可以说是比较普遍了；在nginx中，有很多的配置，完成的功能都是相似的，但是他们又有一些细微的差异，如果在使用中忽略了这些，可能会导致我们多走很多弯路。
本文就主要梳理一下nginx当中`root`，`alias`的差异，和`proxy_pass`的不同的配置方式，对请求uri的不同处理。

### root 和 alias
假设在`/data/html`的目录结构如下：
```
/data
|__/html
   |__/a/b/c/index.html    --①
   |__/b/c/index.html  --②
```

假设有请求： `http://localhost/a/b/c/index.html`
#### root
- 语法
```nginx
Syntax: root path;
Default: root html;
Context: http, server, location, if in location
```
最终映射到的资源是： path/`[uri]` 目录下
- 示例
如果我们的nginx配置如下：
```nginx
location /a {
    root  /data/html/;
    index  index.html index.htm;
 }
```
则最终映射到的 `/data/html/a/` 资源目录下；
所以我们请求中的，得到的是**①** 的html资源文件

#### alias
- 语法
```nginx
Syntax: alias path;
Default:    —
Context:    location
```
最终映射到的资源是： `[path]` 目录下; 正如`**alias**`单词含义，它是对location中的uri的别名;
- 示例
如果我们的nginx配置如下：
```nginx
location /a {
     alias  /data/html/;
     index  index.html index.htm;
 }
```
则最终映射到的 `/data/html/` 资源目录下；所以我们请求中的，得到的是**②** 的html资源文件。
如果我们要访问`①`中的index.html;那么我们的访问请求响应的要改为：`http://localhost/a/a/b/c/index.html`

### proxy_pass的配置
#### proxy_pass
- 语法
```nginx
Syntax: proxy_pass URL;
Default:    —
Context:    location, if in location, limit_except
```
设置代理服务器的协议（http/https）和地址，还可以指定uri映射location

假设当前我们有请求`http://127.0.0.1/a/b/c/index`;通过以下两种配置来理解proxy_pass作用

- 不带uri
```nginx
 location /a/ {
    proxy_pass  http://127.0.0.1:8003;
}
```
当如此配置配置时，最后的请求是`http://127.0.0.1:8003/a/b/c/index`

- 带uri
```nginx
location /a/ {
    proxy_pass  http://127.0.0.1:8003/;
}
```
当如此配置配置时，最后的请求是`http://127.0.0.1:8003/b/c/index`; 此处需要注意的是location中是`/a/`最后的`/`不可省略，如果省略了会被重定向到`http://127.0.0.1/b/c/index`;其实这种配置方式，就是用 `proxy_pass`中的`/`代替了location中的`/a/`



### 结语
本文只是很基础的讲了一些nginx配置中比较常用但是也比较容易弄混的配置点。其实这些在[官网](http://nginx.org/)都可以找到，此处记录，只是为了以后方便查询，毕竟官网文档大而全，要找这一两个小点还是比较麻烦。 
