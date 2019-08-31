package main

import (
	"fmt"
	"os"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "send name and pass and url to /load use post")
	})

	router.POST("/load", Load)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router.Run(":"+port)
}

func Load(c *gin.Context) {
		filename, err := get_pdf(c.PostForm("url"), c.PostForm("filename"))
		if err != nil {
			c.String(http.StatusInternalServerError, "server error")
		}
		err = send_pdf(c.PostForm("token"), filename)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error")
		}
		err = del_file(filename)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error")
		}
		c.String(http.StatusOK, filename+"is done")
}

func get_pdf(URL,name string) (string, error){
	client := &http.Client{}

    // API Explorer | API情報サイト
    // https://dropbox.github.io/dropbox-api-v2-explorer/#files_download
    url := URL
    filename := name+".pdf"
    
    // Create Request | リクエスト作成
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
			return "", err
    }

    // Send Request | 送信
    res, err := client.Do(req)
    if err != nil {
				os.Exit(1)
				return "", err
    }

    // logging Status | ステータスをロギング
    fmt.Println("status:", res.Status)

    // Read Body(=File's Binary) | レスポンス（ファイルのバイナリ）を取得
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
				os.Exit(1)
				return "", err
    }

    // Create Blank File | 空ファイルを作成
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
    if err != nil {
        return "", err
    }

    // Finally Close File | deferでメソッドの最後にファイルを閉じる予約を行う
    defer func() {
				file.Close()
    }()

    // Fill File with Body(=File's Binary) | 空ファイルにバイナリを書き込む
	file.Write(body)

	return filename, nil
}

func send_pdf(token, filename string) (err error){
	//150MB以下のみ
	
	type Send_json struct {
		Path string `json:"path"`
		Mode string `json:"mode"`
	}

	client := &http.Client{}
	url := "https://content.dropboxapi.com/2/files/upload"
	//param := `{ "path": "/"+filename, "mode": "add", "autorename": true, "mute": false, "strict_conflict": false}`
	param, err := json.Marshal(Send_json{Path:"/"+filename, Mode:"add"})
	if err != nil {
			return
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
			return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(file))
	if err != nil {
			return
	}

	// Setup Header | ヘッダを登録
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Dropbox-API-Arg", string(param))
	req.Header.Set("Content-Type", "application/octet-stream")

	res, err := client.Do(req)
	fmt.Println("status:", res.Status)
	return err
}

func del_file(filename string) (err error){
	err = os.Remove(filename)
	return
}