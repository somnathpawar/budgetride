package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer

	queryParams := request.QueryStringParameters
	if queryParams["start_latitude"] == "" ||
		queryParams["start_longitude"] == "" ||
		queryParams["end_latitude"] == "" ||
		queryParams["end_longitude"] == "" {
		//return Response{StatusCode: 400}, errors.New("Please provide Pickup and Drop Points!")
		return ResponseError(http.StatusBadRequest, "Please provide Pickup and Drop Points!")
	}

	allCabs, err := getCabs(
		queryParams["start_latitude"],
		queryParams["start_longitude"],
		queryParams["end_latitude"],
		queryParams["end_longitude"],
		"",
		"")
	if err != nil {
		return ServerError(err)
	}
	if len(allCabs) == 0 {
		return ClientError(http.StatusNotFound)
	}

	body, err := json.Marshal(map[string]interface{}{
		"cabs": allCabs,
	})
	if err != nil {
		return ServerError(err)
	}

	json.HTMLEscape(&buf, body)

	return Success(buf.String())
}

func getCabs(startLn, startLg, endLn, endLg, sortBy, sortOrder string) (CabList, error) {
	avaialableCabs := CabList{}
	// Get uber cabs
	uberCabs, err := UberCabs(startLn, startLg, endLn, endLg)
	if err != nil {
		return nil, err
	}
	avaialableCabs = append(avaialableCabs, uberCabs...)

	// Get Lyft cabs
	lyftCabs, err := LyftCabs(startLn, startLg, endLn, endLg)
	if err != nil {
		return nil, err
	}
	avaialableCabs = append(avaialableCabs, lyftCabs...)

	switch sortBy {
	case "arrival":
		SortBy(func(p1, p2 *CabResponse) bool {
			return p1.Arriving < p2.Arriving
		}).Sort(avaialableCabs, sortOrder)
	default:
		SortBy(func(p1, p2 *CabResponse) bool {
			return p1.Estimate < p2.Estimate
		}).Sort(avaialableCabs, sortOrder)
	}
	return avaialableCabs, nil
}

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func ServerError(err error) (Response, error) {
	errorLogger.Println(err.Error())
	return ResponseError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func ResponseError(status int, data string) (Response, error) {
	var buf bytes.Buffer
	body, _ := json.Marshal(map[string]interface{}{
		"message": data,
	})
	json.HTMLEscape(&buf, body)
	return Response{
		StatusCode: status,
		Body:       buf.String(),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}

func ClientError(status int) (Response, error) {
	return ResponseError(status, http.StatusText(status))
}

func Success(resp string) (Response, error) {
	return Response{
		StatusCode: http.StatusOK,
		Body:       resp,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
	//allCabs, err := getCabs("40.7485413", "-73.98575770000002", "40.6892494", "-74.0445004", "", "")
	//fmt.Println(allCabs, err)
}
