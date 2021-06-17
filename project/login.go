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

func (a *api) showImg(path string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", path).Start()
	case "windows":
		err = exec.Command("cmd", "/C", "start", "", path).Start()
	case "darwin":
		err = exec.Command("open", path).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return errors.Wrap(err, "exec.Command")
	}
	return nil
}

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
	err = a.showImg("./qrcode.jpg")
	if err != nil {
		return errors.Wrap(err, "showImg")
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

func (a *api) QQQRLogin() (nickName string, err error) {
	//获取基础cookie
	cookie, err := a.keTang.XLogin()
	if err != nil {
		return "", errors.Wrap(err, "QQQRLogin")
	}
	a.AddCookie(cookie)
	//获取二维码
	cookie, img, err := a.keTang.PtQrShow()
	if err != nil {
		return "", errors.Wrap(err, "PtQrShow")
	}

	//图片保存起来
	f, err := os.OpenFile("./qrcode.jpg", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return "", errors.Wrap(err, "os.Open")
	}
	_, err = f.Write(img)
	if err != nil {
		f.Close()
		return "", errors.Wrap(err, "f.write")
	}
	f.Close()
	//将cookie保存起来
	a.AddCookie(cookie)

	//将二维码展示
	err = a.showImg("./qrcode.jpg")
	if err != nil {
		return "", errors.Wrap(err, "showImg")
	}

	//计算qtQrToken
	sig := a.getCookieByKey("qrsig")
	e := 0
	for i := 0; i < len(sig); i++ {
		e += (e << 5) + int(sig[i])
	}
	e = e & 2147483647
	//轮询查询登录状态
	for true {
		time.Sleep(3 * time.Second)
		resp, err := a.keTang.PtQrLogin(int64(e), a.getCookieByKey("pt_login_sig"), "", "")
		if err != nil {
			return "", errors.Wrap(err, "keTang.PtQrLogin")
		}
		switch resp.Code {
		case "0":
			//手机授权成功，请求回调连接
			return resp.NickName, nil
		case "65":
			//需要刷新二维码
			return "", errors.New(resp.Msg + "，请关闭二维码图片，并重新输入登录指令")
		case "66":
			//二维码未失效，继续等待
		case "67":
			//手机已扫码，等待手机确认
		}
	}
	return "", nil
}
