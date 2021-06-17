package main

import (
	"crawler/tencentKeTang/config"
	"crawler/tencentKeTang/internal/httplib"
	"crawler/tencentKeTang/pcsliner"
	"crawler/tencentKeTang/pcsliner/args"
	"crawler/tencentKeTang/project"
	"fmt"
	"github.com/peterh/liner"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
)

var (
	Version = "Not Found"
)

func main() {
	myApp := config.NewApp()
	//设置日志
	httplib.SetDebug(myApp.Config.App.Debug)

	//启动程序
	app := cli.NewApp()
	app.Name = "TencentKeTang"
	app.Version = Version
	app.Description = "腾讯课堂视频下载工具"
	app.Authors = []*cli.Author{
		{
			Name:  "HarryWRZ",
			Email: "wrz890829@gmail.com",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg() != 0 {
			return fmt.Errorf("未找到命令: %s\n运行命令 %s help 获取帮助\n", c.Args().Get(0), app.Name)
		}

		var (
			line = pcsliner.NewLiner()
		)
		if myApp.Config.Http.Cookie == "" {
			log.Println("配置文件中未找到cookie,可输入`login -h`查看登录命令")
		}

		for {
			var (
				prompt     string
				activeUser string
			)
			if un, ok := c.App.Metadata["user_name"]; ok && un != nil {
				activeUser = un.(string)
			}

			if activeUser != "" {
				prompt = app.Name + ":" + activeUser + "$ "
			} else {
				// TencentKeTang >
				prompt = app.Name + " > "
			}

			commandLine, err := line.State.Prompt(prompt)
			switch err {
			case liner.ErrPromptAborted:
				return nil
			case nil:
				// continue
			default:
				return err
			}

			cmdArgs := args.Parse(commandLine)
			if len(cmdArgs) == 0 {
				continue
			}

			s := []string{os.Args[0]}
			s = append(s, cmdArgs...)

			// 恢复原始终端状态
			// 防止运行命令时程序被结束, 终端出现异常
			line.Pause()
			err = c.App.Run(s)
			if err != nil {
				fmt.Println(err)
			}
			line.Resume()
		}
	}

	app.Commands = []*cli.Command{
		{
			Name:  "login",
			Usage: "登录课堂",
			Description: `
		登录腾讯课堂
		示例：
		tencentKeTang login 123 456                 #qq帐号密码登录，暂未支持
		tencentKeTang login -type 1 -user_name=123  #qq帐号密码登录，暂未支持
		tencentKeTang login -type 2                 #qq扫码登录，暂未支持
		tencentKeTang login -type 3                 #微信扫码登录
`,
			Action: func(context *cli.Context) error {
				myApp.Project.SetCookie("")
				myApp.KeTang.SetCookie(nil)
				nickName := ""
				var err error
				switch context.Int64("type") {
				case 1:
					log.Println("帐号密码登录，暂未支持")
				case 2:
					nickName, err = myApp.Project.QQQRLogin()
					if err != nil {
						return errors.Wrap(err, "QQQrLogin")
					}
				case 3:
					err = myApp.Project.WxQRLogin()
					if err == project.ErrExpire {
						return fmt.Errorf("登录码超时，请重新申请")
					}
					resp, err := myApp.KeTang.Info()
					if err != nil {
						return err
					}
					nickName = resp.NickName
				default:
					log.Println("帐号密码登录，暂未支持")
				}

				log.Printf("登录成功，欢迎%s使用本程序（请自行关闭二维码）", nickName)
				context.App.Metadata["user_name"] = nickName
				return nil
			},
			Flags: []cli.Flag{
				&cli.Int64Flag{
					Name:  "type",
					Usage: "登录方式",
				},
			},
		},
		{
			Name:      "tree",
			Aliases:   []string{"t"},
			Usage:     "列目录",
			UsageText: app.Name + " tree -c 123 -t 456",
			Description: `
		列出指定cid/term下所有章节与文件（暂时直接只支持单个cid）
		
		示例：
		TencentKeTang tree -c 123 -t 456 (获取123中的456term)
		TencentKeTang t -c 123  (获取123中所有视频)
`,
			Action: func(context *cli.Context) error {
				if context.String("cid") == "" {
					return fmt.Errorf("请填写cid")
				}
				//获取章节列表
				list, err := myApp.Project.GetCatalogue(context.String("cid"), context.Int64("tid"))
				if err != nil {
					return err
				}
				for i, l := range list {
					fmt.Printf("[%d] ", i+1)
					for j := 0; j < int(l.Depth); j++ {
						fmt.Printf("\t")
					}
					fmt.Println(l.Name)
				}
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "cid",
					Usage:   "通过cid下载",
					Aliases: []string{"c"},
				},
				&cli.Int64Flag{
					Name:    "tid",
					Usage:   "通过cid中的termID过滤",
					Aliases: []string{"t"},
				},
			},
		},
		{
			Name:      "download",
			Aliases:   []string{"d"},
			Usage:     "下载文件",
			UsageText: app.Name + " download <cid>",
			Description: `
		下载文件默认保存到当前目录的download目录
		可输入cid直接下载全部内容
		若不填写 flag 则通过tree中的序号下载（支持多序号下载）
		
		示例：
		TencentKeTang d -c 123456
		TencentKeTang d 1
		TencentKeTang d 1 2 3 4
`,
			Action: func(context *cli.Context) error {
				var err error
				switch true {
				case context.IsSet("cid"):
					err = myApp.Project.DownLoadByCID(context.String("cid"))
				default:
					for _, s := range context.Args().Slice() {
						var index int64
						index, err = strconv.ParseInt(s, 10, 64)
						if err != nil {
							return errors.Wrap(err, "strconv.ParseInt")
						}
						err = myApp.Project.DownLoadByIndex(index - 1)
						if err != nil {
							break
						}
					}
				}
				if err != nil {
					return errors.Wrap(err, "myApp.Project.DownLoad*")
				}
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "cid",
					Usage:   "通过cid下载",
					Aliases: []string{"c"},
				},
			},
		},
	}
	app.Run(os.Args)
}
