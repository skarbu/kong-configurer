package logger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	logFileName = "log_" + rawTimestamp() + ".txt"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s \n", msg, err)
	}
}

func LogOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s : %s \n", msg, err)
	}
}

func FailOnErrors(errs []error, msg string) {
	if len(errs) != 0 {
		for _, err := range errs {
			log.Print(err)
		}
		log.Fatal(msg)
	}
}

func LogToFile(request *http.Request, response *http.Response, requestBody string) {
	responseBody, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	response.Body = ioutil.NopCloser(bytes.NewBuffer(responseBody))
	log := fmt.Sprintf("%s: \n request: %s %s \n request body: %s \n response status: %s \n response responseBody: %s \n",
		prettyTimeStamp(),
		request.Method,
		request.URL.Path,
		requestBody,
		response.Status,
		string(responseBody))
	appendToFile(log)
}

func LogMsgToFile(msg string) {
	log := fmt.Sprintf("%s: \n %s \n",
		prettyTimeStamp(),
		msg)
	appendToFile(log)
}

func FatalInvalidArgs() {
	msg := "\nInvalid connection arguments pass to application. Ensure you run application with all required arguments\n" +
		"-u {user} - kong username with admin privileges \n" +
		"-p {port} - password (optional) \n" +
		"-h {host} - kong host \n" +
		"-f {configfile} - relative path to json config file \n"
	log.Fatal(msg)
}

func appendToFile(log string) {
	f, err := os.OpenFile(logFileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	LogOnError(err, "error during creating log file")
	defer f.Close()
	_, err = f.WriteString(log)
	LogOnError(err, "error during saving log file")
}

func rawTimestamp() string {
	now := time.Now()
	return fmt.Sprintf("%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func prettyTimeStamp() string {
	return time.Now().Format(time.RFC3339)
}

func TimeTrack(start time.Time, operationName string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", operationName, elapsed)
}
