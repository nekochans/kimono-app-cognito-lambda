package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	// TriggerSourceが 'CustomMessage_SignUp' の場合はCustomMessageが返却される
	t.Run("Return Signup CustomMessage", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_SignUp",
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		kimonoAppFrontendUrl := os.Getenv("KIMONO_APP_FRONTEND_URL")
		confirmUrl := fmt.Sprintf(
			"%v/accounts/create/confirm?code=%v&userName=%v",
			kimonoAppFrontendUrl,
			createEventParams.CodeParameter,
			createEventParams.UserName,
		)

		m := SignUpMessage{
			ConfirmUrl: confirmUrl,
		}

		body, err := createExpectedSignUpMessage(m)
		if err != nil {
			t.Fatal("Error Failed to parse HTML Template", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body.String(),
			EmailSubject: "アカウント作成 メールアドレスの確認をお願いします。",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_ResendCode' の場合はCustomMessageが返却される
	t.Run("Return ResendCode CustomMessage", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_ResendCode",
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		kimonoAppFrontendUrl := os.Getenv("KIMONO_APP_FRONTEND_URL")
		confirmUrl := fmt.Sprintf(
			"%v/accounts/create/confirm?code=%v&userName=%v",
			kimonoAppFrontendUrl,
			createEventParams.CodeParameter,
			createEventParams.UserName,
		)

		m := SignUpMessage{
			ConfirmUrl: confirmUrl,
		}

		body, err := createExpectedSignUpMessage(m)
		if err != nil {
			t.Fatal("Error Failed to parse HTML Template", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body.String(),
			EmailSubject: "アカウント作成 メールアドレスの確認をお願いします。",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_ForgotPassword' の場合はCustomMessageが返却される
	t.Run("Return ForgotPassword CustomMessage", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_ForgotPassword",
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		kimonoAppFrontendUrl := os.Getenv("KIMONO_APP_FRONTEND_URL")
		confirmUrl := fmt.Sprintf(
			"%v/password/reset/confirm?code=%v&userName=%v",
			kimonoAppFrontendUrl,
			createEventParams.CodeParameter,
			createEventParams.UserName,
		)

		m := ForgotPasswordMessage{ConfirmUrl: confirmUrl}

		body, err := createExpectedForgotPasswordMessageMessage(m)
		if err != nil {
			t.Fatal("Error Failed to parse HTML Template", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "認証コードは {####} です。",
			EmailMessage: body.String(),
			EmailSubject: "パスワードをリセットします。メールの確認をお願いします。",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_SignUp' だがUserPoolIDが一致しないのでDefaultのメッセージが返却される
	t.Run("Return Signup DefaultMessage Because the UserPoolID doesn't match", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_SignUp",
			UserPoolID:    "OtherUserPoolID",
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "SMSMessage",
			EmailMessage: "EmailMessage",
			EmailSubject: "EmailSubject",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_ResendCode' だがUserPoolIDが一致しないのでDefaultのメッセージが返却される
	t.Run("Return ResendCode DefaultMessage Because the UserPoolID doesn't match", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_ResendCode",
			UserPoolID:    "OtherUserPoolID",
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "SMSMessage",
			EmailMessage: "EmailMessage",
			EmailSubject: "EmailSubject",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが 'CustomMessage_ForgotPassword' だがUserPoolIDが一致しないのでDefaultのメッセージが返却される
	t.Run("Return ForgotPassword DefaultMessage Because the UserPoolID doesn't match", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "CustomMessage_ForgotPassword",
			UserPoolID:    "OtherUserPoolID",
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "SMSMessage",
			EmailMessage: "EmailMessage",
			EmailSubject: "EmailSubject",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})

	// TriggerSourceが指定した値以外の場合はDefaultのメッセージが返却される
	t.Run("Return DefaultMessage Because the TriggerSource is not a specified value", func(t *testing.T) {
		createEventParams := &createUserPoolsCustomMessageEventParams{
			TriggerSource: "Unknown",
			UserPoolID:    os.Getenv("TARGET_USER_POOL_ID"),
			UserName:      "keitakn",
			Sub:           "dba1d5db-1d94-45b6-8f1b-fad23bb94cd5",
			CodeParameter: "123456789",
		}

		event := createUserPoolsCustomMessageEvent(createEventParams)
		handlerResult, err := handler(*event)
		if err != nil {
			t.Fatal("Error failed to trigger with an invalid request", err)
		}

		expected := &events.CognitoEventUserPoolsCustomMessageResponse{
			SMSMessage:   "SMSMessage",
			EmailMessage: "EmailMessage",
			EmailSubject: "EmailSubject",
		}

		if reflect.DeepEqual(&handlerResult.Response, expected) == false {
			t.Error("\nActually: ", &handlerResult.Response, "\nExpected: ", expected)
		}
	})
}
