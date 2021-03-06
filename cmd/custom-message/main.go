package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nekochans/kimono-app-cognito-lambda/infrastructure"
	"html/template"
	"log"
	"os"
)

var templates *template.Template

func init() {
	signupTemplatePath := "bin/signup-template.html"
	forgotPasswordTemplatePath := "bin/forgot-password-template.html"
	if infrastructure.IsTestRun() {
		currentDir, _ := os.Getwd()
		signupTemplatePath = currentDir + "/signup-template.html"
		forgotPasswordTemplatePath = currentDir + "/forgot-password-template.html"
	}

	templates = template.Must(template.ParseFiles(signupTemplatePath, forgotPasswordTemplatePath))
}

type SignUpMessage struct {
	ConfirmUrl string
}

type ForgotPasswordMessage struct {
	ConfirmUrl string
}

func BuildSignUpMessage(m SignUpMessage) (*bytes.Buffer, error) {
	var bodyBuffer bytes.Buffer

	err := templates.ExecuteTemplate(&bodyBuffer, "signup-template.html", m)
	if err != nil {
		return nil, err
	}

	return &bodyBuffer, nil
}

func BuildForgotPasswordMessage(m ForgotPasswordMessage) (*bytes.Buffer, error) {
	var bodyBuffer bytes.Buffer

	err := templates.ExecuteTemplate(&bodyBuffer, "forgot-password-template.html", m)
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
		confirmUrl := fmt.Sprintf("%v/accounts/create/confirm?code=%v&userName=%v", frontendUrl, request.Request.CodeParameter, request.UserName)

		m := SignUpMessage{
			ConfirmUrl: confirmUrl,
		}

		body, err := BuildSignUpMessage(m)
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

	// パスワードリセット時に送られる認証メール
	if request.TriggerSource == "CustomMessage_ForgotPassword" {
		frontendUrl := os.Getenv("KIMONO_APP_FRONTEND_URL")
		confirmUrl := fmt.Sprintf("%v/password/reset/confirm?code=%v&userName=%v", frontendUrl, request.Request.CodeParameter, request.UserName)

		m := ForgotPasswordMessage{
			ConfirmUrl: confirmUrl,
		}

		body, err := BuildForgotPasswordMessage(m)
		if err != nil {
			// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
			log.Fatalln(err)
		}

		forgotPasswordMessageResponse := events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body.String(),
			EmailSubject: "パスワードをリセットします。メールの確認をお願いします。",
		}

		request.Response = forgotPasswordMessageResponse
	}

	return request, nil
}

func main() {
	lambda.Start(handler)
}
