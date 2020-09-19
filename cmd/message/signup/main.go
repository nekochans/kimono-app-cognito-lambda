package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

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

		signupMessageResponse := events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage: "認証コードは {####} です。",
			EmailMessage: "メールアドレスを検証するには、次のリンクをクリックしてください。 " +
				confirmUrl,
			EmailSubject: "アカウント作成 メールアドレスの確認をお願いします。",
		}

		request.Response = signupMessageResponse
	}

	return request, nil
}

func main() {
	lambda.Start(handler)
}
