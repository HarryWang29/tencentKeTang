# 项目的由来
哎，腾讯课堂的app太难用了，此工具仅仅只是为了将视频下载到本地，使用第三方播放器使用

# 注意事项
* 请自行下载安装ffmpeg与ffprobe
* ffmpeg使用gpu加速，请自行查找资料
* ~~未测试windows使用情况~~

# 使用帮助
1. 请自行下载安装ffmpeg与ffprobe
1. 请确保ffmpeg与ffprobe在同一目录
2. 将ffmpeg安装目录填写到config.yaml中
3. ffmpeg使用gpu加速相关，请自行查找资料
4. 确认文件下载路径
5. 目前已支持微信扫码登录/cookie登录
   1. 通过cookie登录
      1. 浏览器登录后，f12 --> NetWork
      2. 查找 `https://ke.qq.com/cgi-bin/identity/info` 接口请求
      3. 复制cookie到配置文件
   1. 微信扫码登录
      1. 启动程序后输入：`login -type 3`
      2. 手机扫码确认登录
      3. 手动关闭二维码
6. 执行命令，启动程序(Windows现在可以双击启动程序，而不依托cmd)
   ```shell
      tencentKeTang
   ```
7. 执行命令，可查看说明
   ```shell
      help
   ```
8. 可以通过cid直接下载，也可通过cid+tid列出目录后，进行选择下载，如下图所示
   ```shell
      tree -c 123  #获取123中所有视频
      tree -c 123 -t 456 #获取123中的456term
   
      d -c 123  #下载123中所有视频
      d -t 1    #下载tree列目录中的索引
   ```
   ![image](https://user-images.githubusercontent.com/8288067/121004497-585c6d80-c7c1-11eb-9f3c-c7b51785baf2.png)


# TODO List
- [X] 整理日志
- [X] 可通过终端选择要下载的文件
- [X] 显示下载进度
- [X] 优化进度条
- [X] 支持微信扫码登录
- [ ] 打包ffmpeg/ffprobe
- [ ] 支持qq扫码登录
- [ ] 支持qq帐号密码登录

# 感谢
- 感谢腾讯课堂给我们的优质内容，不过app真的不好用。。。
