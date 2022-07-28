# 项目的由来
哎，腾讯课堂的app太难用了，此工具仅仅只是为了将视频下载到本地，使用第三方播放器使用

# v0.2.9更新
* 打脸来的太快，这两天实现了终端多行进度条……

# v0.2.8更新
* 适配了腾讯课堂最新接口
* 增加了当前执行任务的提示
* 由于目前去除了进度条，并且下载ts和合并生成视频是异步过程，在网络良好的情况下，ts的下载一定是领先合并很多的，所以麻烦大家观察下正在合并的视频是否还在增长大小
* Windows使用中，建议把视频下载目录和应用程序放在同一个盘符下
* 建议使用`-c:v copy`命令

# v0.2.6更新
* 很遗憾最后一个版本去除了进度条，而且应该不会再加入回来，使用者需要自行判断视频是否下载完成了，一般可通过cpu/gpu使用率来判断
* 此版本采用了多线程下载分片视频，但是合并视频是使用的单线程，在一定程度上加快了速度
* 经过研究发现，下载后的视频清晰度不足，主要原因是码率的问题，此版本已修复

# 后续计划
* 此项目不会再有大型功能更新，最多是对一些小bug进行修复
* 周围有些朋友找到我，帮忙从这个地方下下视频，那个地方下下视频，写了很多小脚本，后续会开一个新项目，集合成新的工具，欢迎大家使用
* 最近重新看了这个项目代码……写的很糟糕……

# 注意事项
* 请自行下载安装ffmpeg与ffprobe
* ffmpeg使用gpu加速，请自行查找资料
* 未测试windows使用情况

# 使用帮助
1. 请自行下载安装ffmpeg与ffprobe
1. 请确保ffmpeg与ffprobe在同一目录
2. 将ffmpeg安装目录填写到config.yaml中
3. ffmpeg使用gpu加速相关，请自行查找资料
4. 确认文件下载路径
5. 目前已支持qq扫码登录/微信扫码登录/cookie登录
   1. 通过cookie登录
      1. 浏览器登录后，f12 --> NetWork
      2. 查找 `https://ke.qq.com/cgi-bin/identity/info` 接口请求
      3. 复制cookie到配置文件
   1. 微信扫码登录
      1. 启动程序后输入：`login -type 3`
      2. 手机扫码确认登录
      3. 手动关闭二维码
   1. qq扫码登录
      1. 启动程序后输入：`login -type 2`
      2. 手机扫码确认登录
      3. 手动关闭二维码
      4. 若出现二维码已失效，需手动关闭二维码图片，并重新输入`login -type 2`
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
      tree -u https://ke.qq.com/course/123?taid=1234 #通过url获取cid
      tree -u https://ke.qq.com/webcourse/index.html#cid=1111&term_id=2222&taid=3333&type=4444&vid=55555    #通过url获取cid/tid
   
      d -c 123  #下载123中所有视频
      d 1    #下载tree列目录中的索引
      d 1 3 5    #下载tree列目录中的索引
   ```
   ![image](https://user-images.githubusercontent.com/8288067/121004497-585c6d80-c7c1-11eb-9f3c-c7b51785baf2.png)

# TODO List
- [X] 整理日志
- [X] 可通过终端选择要下载的文件
- [X] 显示下载进度
- [X] 优化进度条
- [X] 支持微信扫码登录
- [X] 从列表中选择多个视频下载
- [X] 打包ffmpeg/ffprobe(linux没有打包，目前只打包了mac/windows, 
  mac: ffmpeg version 4.3.1, ffprobe version 4.3.1
  windows: ffmpeg version 4.4-essentials_build-www.gyan.dev, ffprobe version 4.4-essentials_build-www.gyan.dev)
- [X] 支持qq扫码登录
- [ ] 支持qq帐号密码登录
- [ ] 增加桌面版界面

# 感谢
- 感谢腾讯课堂给我们的优质内容，不过app真的不好用。。。
