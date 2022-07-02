package qzone

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/requests"
)

// GenderGTK 生成GTK
func GenderGTK(sKey string) string {
	hash := 5381
	for _, s := range sKey {
		us, _ := strconv.Atoi(fmt.Sprintf("%d", s))
		hash += (hash << 5) + us
	}
	return fmt.Sprintf("%d", hash&0x7fffffff)
}

type Manager struct {
	r          *requests.Request
	QQ         string
	Gtk        string
	Gtk2       string
	PSkey      string
	Skey       string
	Uin        string
	qzoneToken string
}

func NewQzoneManager(qq int64, cookie OPQBot.Cookie) Manager {
	var m Manager
	m.Skey = cookie.Skey
	m.PSkey = cookie.PSkey.Qzone
	m.Gtk = GenderGTK(m.Skey)
	m.Gtk2 = GenderGTK(m.PSkey)
	m.Uin = "o" + strconv.FormatInt(qq, 10)
	r := requests.Requests()

	c := &http.Cookie{
		Name:  "pt2gguin",
		Value: "o" + strconv.FormatInt(qq, 10),
	}
	r.SetCookie(c)
	c = &http.Cookie{
		Name:  "uin",
		Value: m.Uin,
	}
	r.SetCookie(c)
	c = &http.Cookie{
		Name:  "skey",
		Value: m.Skey,
	}
	r.SetCookie(c)
	c = &http.Cookie{
		Name:  "p_skey",
		Value: m.PSkey,
	}
	r.SetCookie(c)
	c = &http.Cookie{
		Name:  "p_uin",
		Value: m.Uin,
	}
	r.SetCookie(c)
	r.Header.Set("user-agent", "User-Agent: Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	m.r = r
	for _, v := range strings.Split(cookie.Cookies, ";") {
		if v2 := strings.Split(v, "="); len(v2) == 2 {
			c = &http.Cookie{
				Name:  v2[0],
				Value: v2[1],
			}
			r.SetCookie(c)
		}

	}
	m.QQ = strconv.FormatInt(qq, 10)
	return m
}
func (m *Manager) GetQzoneToken() (string, error) {
	if m.qzoneToken == "" {
		return "", errors.New("请先刷新QzoneToken！")
	}
	return m.qzoneToken, nil

}
func (m *Manager) RefreshToken() error {
	res, err := m.r.Get("https://h5.qzone.qq.com/feeds/inpcqq?uin=" + m.QQ + "&qqver=5749&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return err
	}
	r, err := regexp.Compile(`window.g_qzonetoken.*try{return "(.*?)";} catch\(e\)`)
	if err != nil {
		return err
	}
	result := r.FindStringSubmatch(res.Text())
	if len(result) == 2 {
		m.qzoneToken = result[1]
		return nil
	}
	return errors.New("获取qzonetoken失败 ")
}
func (m *Manager) GetShuoShuoList() (ShuoshuoList, error) {
	m.r.Header.Set("referer", "https://user.qzone.qq.com/"+m.QQ)
	c := &http.Cookie{
		Name:  "Loading",
		Value: "Yes",
	}
	m.r.SetCookie(c)
	var ss ShuoshuoList
	res, err := m.r.Get("https://user.qzone.qq.com/proxy/domain/ic2.qzone.qq.com/cgi-bin/feeds/feeds3_html_more?format=json&begintime=" + strconv.FormatInt(time.Now().Unix(), 10) + "&count=20&uin=" + m.QQ + "&g_tk=" + m.Gtk2)
	if err != nil {
		return ss, err
	}
	vm := goja.New()
	v, err := vm.RunString("c=" + res.Text() + ";JSON.stringify(c);")
	if err != nil {
		return ss, err
	}
	err = json.Unmarshal([]byte(v.String()), &ss)
	if err != nil {
		return ss, err
	}
	return ss, nil
}
func (m *Manager) SendShuoShuo(Content string) (SendShuoShuoResult, error) {
	token, err := m.GetQzoneToken()
	var result SendShuoShuoResult
	if err != nil {
		return result, err
	}
	m.r.Header.Set("referer", "https://user.qzone.qq.com/"+m.QQ)
	m.r.Header.Set("Origin", "https://user.qzone.qq.com/")
	log.Println(m.r.Header)
	res, err := m.r.Post("https://user.qzone.qq.com/proxy/domain/taotao.qzone.qq.com/cgi-bin/emotion_cgi_publish_v6?g_tk="+m.Gtk2+"&qzonetoken="+token+"&uin="+m.QQ, requests.Datas{
		"syn_tweet_verson": "1",
		"paramstr":         "1",
		"who":              "1",
		"con":              Content,
		"feedversion":      "1",
		"ver":              "1",
		"ugc_right":        "1",
		"to_sign":          "0",
		"hostuin":          m.QQ,
		"code_version":     "1",
		"format":           "json",
		"qzreferrer":       "https://user.qzone.qq.com/" + m.QQ,
	})
	if err != nil {
		return result, err
	}
	err = res.Json(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (m *Manager) UploadPic(picBase64 string) (UploadPicResult, error) {
	token, err := m.GetQzoneToken()
	var result UploadPicResult
	if err != nil {
		return result, err
	}
	m.r.Header.Set("referer", "https://user.qzone.qq.com/"+m.QQ)
	m.r.Header.Set("Origin", "https://user.qzone.qq.com/")
	//log.Println(m.r.Header)
	res, err := m.r.Post("https://up.qzone.qq.com/cgi-bin/upload/cgi_upload_image?g_tk="+m.Gtk2+"&qzonetoken="+token+"&uin="+m.QQ, requests.Datas{
		"filename":       "filename",
		"zzpanelkey":     "",
		"uploadtype":     "1",
		"albumtype":      "7",
		"exttype":        "0",
		"skey":           m.Skey,
		"zzpaneluin":     m.QQ,
		"p_uin":          m.QQ,
		"uin":            m.QQ,
		"p_skey":         m.PSkey,
		"output_type":    "json",
		"qzonetoken":     token,
		"refer":          "shuoshuo",
		"charset":        "utf-8",
		"output_charset": "utf-8",
		"upload_hd":      "1",
		"hd_width":       "2048",
		"hd_height":      "10000",
		"hd_quality":     "96",
		"backUrls":       "http://upbak.photo.qzone.qq.com/cgi-bin/upload/cgi_upload_image,http://119.147.64.75/cgi-bin/upload/cgi_upload_image",
		"url":            "https://up.qzone.qq.com/cgi-bin/upload/cgi_upload_image?g_tk=" + m.Gtk2,
		"base64":         "1",
		"picfile":        picBase64,
	})
	if err != nil {
		return result, err
	}
	r, _ := regexp.Compile(`_Callback\((.*)\)`)
	r1 := r.FindStringSubmatch(res.Text())
	if len(r1) != 2 {
		return result, errors.New("提取失败! ")
	}
	err = json.Unmarshal([]byte(r1[1]), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (m *Manager) SendShuoShuoWithBase64Pic(Content string, pics []string) (SendShuoShuoResult, error) {
	var result SendShuoShuoResult

	//获取qzone token
	token, err := m.GetQzoneToken()
	if err != nil {
		return SendShuoShuoResult{}, err
	}

	//包装request
	m.r.Header.Set("referer", "https://user.qzone.qq.com/"+m.QQ)
	m.r.Header.Set("Origin", "https://user.qzone.qq.com/")
	postDatas := requests.Datas{
		"syn_tweet_verson": "1",
		"paramstr":         "1",
		"who":              "1",
		"con":              Content,
		"feedversion":      "1",
		"ver":              "1",
		"ugc_right":        "1",
		"to_sign":          "0",
		"hostuin":          m.QQ,
		"code_version":     "1",
		"format":           "json",
		"qzreferrer":       "https://user.qzone.qq.com/" + m.QQ,
	}

	//上传图片
	richvals := make([]string, 0)
	pic_bos := make([]string, 0)

	if len(pics) != 0 {
		//挨个上传图片
		for _, pic := range pics {
			uploadPicResult, err := m.UploadPic(pic)
			if err != nil {
				return SendShuoShuoResult{}, err
			}

			//检出此图片的picbo和richval
			picbo, richva, err := GetPicBoAndRichVal(uploadPicResult)
			if err != nil {
				return SendShuoShuoResult{}, err
			}

			richvals = append(richvals, richva)
			pic_bos = append(pic_bos, picbo)
		}
		//打包进postDatas
		finalRichVal := strings.Join(richvals, "\t")
		finalPicBo := strings.Join(pic_bos, ",")

		postDatas["pic_bo"] = finalPicBo
		postDatas["richtype"] = "1"
		postDatas["richval"] = finalRichVal
	}

	res, err := m.r.Post("https://user.qzone.qq.com/proxy/domain/taotao.qzone.qq.com/cgi-bin/emotion_cgi_publish_v6?g_tk="+m.Gtk2+"&qzonetoken="+token+"&uin="+m.QQ, postDatas)

	if err != nil {
		return result, err
	}
	err = res.Json(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (m *Manager) SendShuoShuoWithLocalPic(content string, pics []string) (SendShuoShuoResult, error) {
	if len(pics) != 0 {
		//挨个上传图片
		base64pics := make([]string, len(pics))
		for i, pic := range pics {
			picBase64, err := GetBase64(pic)
			if err != nil {
				return SendShuoShuoResult{}, err
			}
			base64pics[i] = picBase64
		}
		return m.SendShuoShuoWithBase64Pic(content, base64pics)
	}
	return SendShuoShuoResult{}, errors.New("Hi give me your pics!!!")
}
func (m *Manager) Like(unikey, curkey, appid string) error {
	token, err := m.GetQzoneToken()
	if err != nil {
		return err
	}
	m.r.Header.Set("referer", "https://user.qzone.qq.com/"+m.QQ)
	m.r.Header.Set("Origin", "https://user.qzone.qq.com/")
	//log.Println(m.r.Header)
	res, err := m.r.Post("https://user.qzone.qq.com/proxy/domain/w.qzone.qq.com/cgi-bin/likes/internal_dolike_app?g_tk="+m.Gtk2+"&qzonetoken="+token+"&uin="+m.QQ, requests.Datas{
		"opuin":      m.QQ,
		"unikey":     unikey,
		"curkey":     curkey,
		"appid":      appid,
		"opr_type":   "like",
		"format":     "json",
		"qzreferrer": "https://user.qzone.qq.com/" + m.QQ,
	})
	if err != nil {
		return err
	}
	//err = res.Json(&result)
	//if err != nil {
	//	return result, err
	//}
	log.Println(res.Text())
	return nil
}
func GetPicBoAndRichVal(data UploadPicResult) (PicBo, RichVal string, err error) {
	if data.Ret != 0 {
		err = errors.New("返回错误")
		return
	}
	if v := strings.Split(data.Data.URL, "&bo="); len(v) >= 2 {
		PicBo = v[len(v)-1]
	} else {
		err = errors.New("bo数据错误")
		return
	}
	RichVal = fmt.Sprintf(",%s,%s,%s,%d,%d,%d,,%d,%d", data.Data.Albumid, data.Data.Lloc, data.Data.Sloc, data.Data.Type, data.Data.Height, data.Data.Width, data.Data.Height, data.Data.Width)
	return
}

func GetBase64(path string) (string, error) {
	srcByte, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	res := base64.StdEncoding.EncodeToString(srcByte)

	return res, nil
}
