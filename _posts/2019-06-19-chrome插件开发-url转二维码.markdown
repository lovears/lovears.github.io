---
layout: post
title: 【教程】 chrome插件开发
subtitle: chrome插件开发之地址栏url转二维码插件
date: 2019-06-19
author: "wangnem"
header-img: 
tags:
  - 教程
  - 二维码
  - chrome插件
  - qrcode.js
---
有些时候，我们可能在电脑上看到好看的文章或者视频，但是可能会因为什么事，导致不能一直在电脑上看； 这时候我们可能需要把链接发送到移动设备上，这样看起来可能也
更方便。当然可以使用QQ，微信或者其他的传输软件，发送到手机上。但是在有些网络下，电脑可能会限制使用这种社交软件，这时我们要把链接放到手机上可能就
需要自己在手机上一个一个字符敲击出来。 作为一个程序员，而且是很懒的那种程序员，我可不想这样干，太麻烦。  
此时我想到了二维码，我想着将链接地址通过在线网站生成二维码，然后用手机扫一下就可以了。但是这个流程还是太复杂。懒人可不想这么干！于是我想着用chrome的插件，
直接点击一下就将连接生成二维码，这样方便多了。  
有了这个想法，就开始在网上找使用js生成二维码的教程，这个太简单了，一搜一大把，而且使用也特简单，基本一句代码就解决了。然后就是chrome插件开发，这个是之前从
来没有做过的。所以，我想着可能需要学好久才行。在查阅一番开发资料后（因为网络原因，能看的教程就那么些，但是也足够了），我觉得我可能想得太复杂了；其实chrome插
件可能就只需要一个`manifest.json`文件，里面是一些配置，然后将自己的脚本或者html页面配置进去就行了。然后，我就着手开发了。
‘
### 开发准备
- 打开插件管理，启用开发者模式
```
chrome://extensions
```
- `加载已解压的扩展程序` 
这里选中我们的插件工程目录即可  

### 工程目录结构
```
js
  |_qrcode.min.js
  |_jquery.min.js
  |_backgroud.js
qrcode.html
qrcode.js
manifest.json
```  

注意`**manifes.json**`一定要在根目录下，其他的随意

#### manifes.json 配置  
```
{
  "manifest_version": 2,
  "name": "url2qrcode",
  "description": "url转二维码",
  "version": "1.0",
  
  "background": {
      "scripts": ["js/jquery.min.js","js/qrcode.min.js","js/background.js","qrcode.js"],
      "persistent": false
    },

    "browser_action": {
      "default_popup": "qrcode.html"
    },

  "permissions": [
      "tabs"
  ]
}
```  

- 配置说明：  

```
manifest_version: 这个必须是2
name: 插件名字
version: 插件版本
background: 一些后台常驻的JS或页面
browser_action: 浏览器右上角的插件角标，此外还有`page_action`和`app`；这里我只用了browser_action
permissions: 会用到的权限列表，这里我需要用到`tabs`，因为要获取tab的url
```  

其他具体的参数可以查看[官方的文档](https://developer.chrome.com/extensions/getstarted)，此处只用了最简单的一些。

#### 工程文件
- qrcode.html  

```html
<!DOCTYPE html>
<html>
<head>
	<title></title>
</head>
<body>
<div>
<div id="qrcode"></div>
</div>
<script src="js/jquery.min.js"></script>
<script src="js/qrcode.min.js"></script>
<script src="js/background.js"></script>
<script src="qrcode.js"></script>
</body>
</html>

```

- qrcode.js  

```js

// $("#qrcode").qrcode( "http://www.runoob.com"); 
chrome.tabs.onUpdated.addListener(onUpdated);

getCurrentUrl(function(url){
	createQrCode(url); 
})

$("#generate").click(function(){
	getCurrentUrl(function(url){
			createQrCode(url); 
		})
})

function createQrCode(url) {
	$("body").empty();
	new QRCode(document.body, url); 
}
function onUpdated(tabId,changeInfo,tab) {
		createQrCode(tab.url)
}

```  

- background.js： 这个可以和`qrcode.js`合并  

```js

// 获取当前选项卡url
function getCurrentUrl(callback)
{
	chrome.tabs.query({active: true, currentWindow: true}, function(tabs)
	{
		if(callback) callback(tabs.length ? tabs[0].url: null);
	});
}
```  

- `js/qrcode.min.js`使用的是 [https://github.com/davidshimjs/qrcodejs](https://github.com/davidshimjs/qrcodejs)的`qrcode.min.js`
  
  
  
至此，基本开发完成，因为是自己用，也就没有加图标啊什么的，有兴趣的可以自己参阅[官网文档](https://developer.chrome.com/extensions/getstarted)开发

