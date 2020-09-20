package main

import (
	"bytes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nekochans/kimono-app-cognito-lambda/infrastructure"
	"html/template"
	"log"
	"os"
)

var templates *template.Template

func init() {
	t := template.New("signup-template.html")

	templatePath := "bin/message/signup-template.html"
	if infrastructure.IsTestRun() {
		currentDir, _ := os.Getwd()
		templatePath = currentDir + "/signup-template.html"
	}

	templates = template.Must(t.ParseFiles(templatePath))
}

type Message struct {
	ConfirmUrl string
}

func BuildMessage(ms Message) (*bytes.Buffer, error) {
	var bodyBuffer bytes.Buffer

	err := templates.Execute(&bodyBuffer, ms)
	if err != nil {
		return nil, err
	}

	return &bodyBuffer, nil
}

func handler(request events.CognitoEventUserPoolsCustomMessage) (events.CognitoEventUserPoolsCustomMessage, error) {
	targetUserPoolId := os.Getenv("TARGET_USER_POOL_ID")

	// 実行対象のユーザープールからのリクエストではない場合は何もせずにデフォルトのメッセージを返す
	if targetUserPoolId != request.UserPoolID {
		return request, nil
	}

	// サインアップ時に送られる認証メール
	if request.TriggerSource == "CustomMessage_SignUp" || request.TriggerSource == "CustomMessage_ResendCode" {
		frontendUrl := os.Getenv("KIMONO_APP_FRONTEND_URL")
		confirmUrl := frontendUrl + "/accounts/create/confirm?code=" + request.Request.CodeParameter + "&sub=" + request.UserName

		ms := Message{
			ConfirmUrl: confirmUrl,
		}

		body, err := BuildMessage(ms)
		if err != nil {
			// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
			log.Fatalln(err)
		}

		signupMessageResponse := events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body.String(),
			EmailSubject: "アカウント作成 メールアドレスの確認をお願いします。",
		}

		request.Response = signupMessageResponse
	}

	return request, nil
}

func main() {
	lambda.Start(handler)
}
