package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/yoshio/gambituser/awsgo"
	"github.com/yoshio/gambituser/bd"
	"github.com/yoshio/gambituser/models"
)

func main() {
	lambda.Start(EjecutoLambda)
}

func EjecutoLambda(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	awsgo.InicializoAWS()

	//Las variables de entorno viene de 'os' configurar en aws

	//Validar el parámetro
	if !ValidoParametros() {
		fmt.Println("Error en los parámetros debe de enviar 'SecretManager'")
		err := errors.New("Error en los parámetros debe de enviar SecretName")
		return event, err
	}
	var datos models.SignUp

	for row, att := range event.Request.UserAttributes {
		switch row {
		case "email":
			datos.UserEmail = att
			fmt.Println("Email = " + datos.UserEmail)
		case "sub":
			datos.UserUUID = att
			fmt.Println("Sub = " + datos.UserUUID)
		}
	}

	err := bd.ReadSecret()

	if err != nil {
		fmt.Println("Error al leer el Secret " + err.Error())
		return event, err
	}

	err = bd.SignUp(datos)
	return event, nil

}

func ValidoParametros() bool {
	var traeParametro bool

	_, traeParametro = os.LookupEnv("SecretName") //Variable de entorno de nombre SecretName
	return traeParametro
}
