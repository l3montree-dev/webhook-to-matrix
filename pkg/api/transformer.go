package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/go-jsonnet"
)

type AppModel string

const (
	GlitchTip AppModel = "glitchtip"
	Botkube            = "botkube"
)

func TransformGlitchTip(res http.ResponseWriter, req *http.Request) {
	transform(res, req, GlitchTip)
}

func TransformBotKube(res http.ResponseWriter, req *http.Request) {
	transform(res, req, Botkube)
}

func bodyToString(req *http.Request) (*string, error) {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(body)
	print(bodyStr)
	return &bodyStr, nil
}

func convertRawJsonToMatrixMessage(jsonStr string, transformationType AppModel) (*MatrixMessage, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(b))

	transformationFilePath := fmt.Sprintf("%s/models/%s.libsonnet", basepath, transformationType)

	// Read the Jsonnet mapping code
	mappingCode, err := os.ReadFile(transformationFilePath)
	if err != nil {
		return nil, err
	}

	vm := jsonnet.MakeVM()
	vm.ExtVar("input", jsonStr)
	output, err := vm.EvaluateAnonymousSnippet(transformationFilePath, string(mappingCode))
	if err != nil {
		return nil, err
	}

	var msg MatrixMessage
	err = json.Unmarshal([]byte(output), &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func transform(res http.ResponseWriter, req *http.Request, transformationType AppModel) {
	body, err := bodyToString(req)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		http.Error(res, "cannot read body", http.StatusBadRequest)
		return
	}

	msg, err := convertRawJsonToMatrixMessage(*body, transformationType)
	if err != nil {
		log.Printf("failed to convert: %v", err)
		http.Error(res, "failed to convert message", http.StatusBadRequest)
		return
	}

	err = SendMessage(*msg)
	if err != nil {
		log.Printf("Matrix send error: %v", err)
		http.Error(res, "Matrix send error", 500)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("ok"))
}
