# 项目的由来
哎，腾讯课堂的app太难用了，此工具仅仅只是为了将视频下载到本地，使用第三方播放器使用

# 注意事项
* 请自行下载安装ffmpeg
* ffmpeg使用gpu加速，请自行查找资料
* ~~未测试windows使用情况~~

# 使用帮助
1. 请自行下载安装ffmpeg
2. 将ffmpeg安装目录填写到config.yaml中
3. ffmpeg使用gpu加速相关，请自行查找资料
4. 确认文件下载路径
5. 目前没有对接腾讯的登录体系，所以需要用户自己找到cookie
   1. 浏览器登录后，f12 --> NetWork
   2. 查找 `https://ke.qq.com/cgi-bin/identity/info` 接口请求
   3. 复制cookie到配置文件
6. 获取目录播放网址，例如 `https://ke.qq.com/webcourse/index.html#course_id=aaa&term_id=bbb&taid=ccc&type=ddd&vid=eee`
7. 执行命令
   ```shell
      tencentKeTang -u "https://ke.qq.com/webcourse/index.html#course_id=aaa&term_id=bbb&taid=ccc&type=ddd&vid=eee"
   ```

# TODO List
- [X] 整理日志
- [X] 可通过终端选择要下载的文件
- [X] 显示下载进度
- [ ] 优化进度条
- [ ] 对接腾讯登录

# 感谢
感谢腾讯课堂给我们的优质内容，不过app真的不好用。。。