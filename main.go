package main

import (
	"crawler/tencentKeTang/config"
	"crawler/tencentKeTang/internal/httplib"
	"crawler/tencentKeTang/model"
	"crawler/tencentKeTang/project"
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

func doTodo(todo *model.TaskInfo, token *model.TokenResp) error {
	url := fmt.Sprintf("https://playvideo.qcloud.com/getplayinfo/v2/1258712167/%s?t=%s&sign=%s&us=%s&exper=0",
		todo.Video.Vid, token.Result.T, token.Result.Sign, token.Result.Us)
	body, err := httplib.Get(nil, url).Bytes()
	if err != nil {
		return err
	}
	taskResp := &model.TaskResp{}
	err = json.Unmarshal(body, taskResp)
	if err != nil {
		return err
	}
	vodUrl := taskResp.VideoInfo.TranscodeList[len(taskResp.VideoInfo.TranscodeList)-1].URL
	i := strings.LastIndex(vodUrl, "/")
	vodUrl = vodUrl[:i+1] + "voddrm.token.dWluPTEwNjEyNjUwNjI7c2tleT1AZktsQkN6RGFGO3Bza2V5PVkqdXF4MDZ5cUtwWVN0YzNsM2R2WEFrVlp1QjJ0UkNtLTVEem5IVlp1VWtfO3Bsc2tleT0wMDA0MDAwMGVhMTk4YzIwYWM2MjYwNjllMmYxMmM2YTNiMzFjMTIyZTkyNWFjM2RmNjQ5YjFkYzM5ODM1YTBkOTkyZjZiNzVjYTBkYjg4YmFmOTBlNjA2O2V4dD07dWlkX3R5cGU9MDt1aWRfb3JpZ2luX3VpZF90eXBlPTA7Y2lkPTMxMzI4MTU7dGVybV9pZD0xMDMyNTYxOTQ7dm9kX3R5cGU9MA==." + vodUrl[i+1:]
	fmt.Println(vodUrl)
	cmd := exec.Command("sh", "-c", fmt.Sprintf(`ffmpeg -i "%s" -c:v h264_videotoolbox "./tmp/%s"`, vodUrl, todo.Video.Name))
	ret, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(ret))
	return nil
}

func main() {
	//加载入参
	taskUrl := ""
	flag.StringVar(&taskUrl, "u", "", "初始任务链接")
	flag.Parse()
	if taskUrl == "" {
		panic("url is empty")
	}
	//加载配置文件
	c := config.Load("./config.yaml")
	if c.Http.Cookie == "" {
		panic("cookie is empty")
	}
	//执行任务
	if err := project.New(c).Do(taskUrl); err != nil {
		panic(err)
	}
	return
}
