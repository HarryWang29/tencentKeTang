package project

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func (a *api) WxQRLogin() (err error) {
	//这两个参数没有找到来源
	var rd float64 = 0.26400682992232993
	var r float64 = 0.3950
	//获取二维码
	qrResp, err := a.keTang.MiniAppQrcode(rd, r)
	if err != nil {
		return errors.Wrap(err, "miniAppQrcode")
	}
	qrByte, err := base64.StdEncoding.DecodeString(qrResp.Result.QRCode) //成图片文件并把文件写入到buffer
	if err != nil {
		return errors.Wrap(err, "base64.StdEncoding.DecodeString")
	}
	qrBuf := bytes.NewBuffer(qrByte)
	qrcode, _, err := image.Decode(qrBuf)
	if err != nil {
		return errors.Wrap(err, "image.Decode")
	}
	//图片保存起来
	f, err := os.OpenFile("./qrcode.jpg", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return errors.Wrap(err, "os.Open")
	}
	err = jpeg.Encode(f, qrcode, nil)
	if err != nil {
		f.Close()
		return errors.Wrap(err, "jpeg.Encode")
	}
	f.Close()
	//将cookie保存起来
	a.AddCookie(qrResp.Cookie)

	//将二维码展示
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", "./qrcode.jpg").Start()
	case "windows":
		err = exec.Command("cmd", "/C", "start", "", "./qrcode.jpg").Start()
	case "darwin":
		err = exec.Command("open", "./qrcode.jpg").Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return errors.Wrap(err, "exec.Command")
	}

	//轮询查询登录状态
	var loginState int
	for true {
		time.Sleep(time.Second)
		state, err := a.keTang.LoginState(rd, r)
		if err != nil {
			break
		}
		//state:0 未登录 1 登录成功 3 二维码失效
		if state.Result.State != 0 {
			loginState = state.Result.State
			//将应答cookie保存
			a.AddCookie(state.Cookies)
			break
		}
	}
	if loginState != 1 {
		return ErrExpire
	}
	//获取cookie
	a2Resp, err := a.keTang.A2Login("", r)
	if err != nil {
		return errors.Wrap(err, "keTang.A2Login")
	}
	a.AddCookie(a2Resp.Cookies)

	return nil
}
