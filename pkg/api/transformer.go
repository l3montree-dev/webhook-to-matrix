package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-jsonnet"

	_ "embed"
)

type AppModel string

const (
	GlitchTip              AppModel = "glitchtip"
	Botkube                AppModel = "botkube"
	DevGuard               AppModel = "devguard"
	Github                 AppModel = "github"
	Gitlab                 AppModel = "gitlab"
	DocumentationAssigment AppModel = "documentationassignment"
)

//go:embed models/glitchtip.libsonnet
var mappingCodeGlitchTip string

//go:embed models/botkube.libsonnet
var mappingCodeBotKube string

//go:embed models/devguard.libsonnet
var mappingCodeDevGuard string

//go:embed models/github.libsonnet
var mappingCodeGithub string

//go:embed models/gitlab-opencode.libsonnet
var mappingCodeGitlab string

//go:embed models/documentationassignment.libsonnet
var mappingCodeDocumentationAssigment string

func TransformGlitchTip(res http.ResponseWriter, req *http.Request) {
	transform(res, req, GlitchTip, mappingCodeGlitchTip)
}

func TransformBotKube(res http.ResponseWriter, req *http.Request) {
	transform(res, req, Botkube, mappingCodeBotKube)
}

func TransformDevGuard(res http.ResponseWriter, req *http.Request) {
	transform(res, req, DevGuard, mappingCodeDevGuard)
}

func TransformGithub(res http.ResponseWriter, req *http.Request) {
	transform(res, req, Github, mappingCodeGithub)
}

func TransformGitlab(res http.ResponseWriter, req *http.Request) {
	transform(res, req, Gitlab, mappingCodeGitlab)
}

func TransformDocumentationAssigment(res http.ResponseWriter, req *http.Request) {
	transform(res, req, DocumentationAssigment, mappingCodeDocumentationAssigment)
}

func bodyToString(req *http.Request) (*string, error) {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(body)
	return &bodyStr, nil
}

func convertRawJsonToMatrixMessage(jsonStr string, transformationType AppModel, mappingCode string) (*MatrixMessage, error) {
	debugName := fmt.Sprintf("%s.libsonnet", transformationType)

	vm := jsonnet.MakeVM()
	vm.ExtVar("input", jsonStr)
	output, err := vm.EvaluateAnonymousSnippet(debugName, string(mappingCode))
	if err != nil {
		return nil, err
	}

	// Check if output is null (for ignored events/users)
	if strings.TrimSpace(output) == "null" {
		return nil, nil
	}

	var msg MatrixMessage
	err = json.Unmarshal([]byte(output), &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func transform(res http.ResponseWriter, req *http.Request, transformationType AppModel, mappingCode string) {
	roomID := req.URL.Query().Get("roomid")
	if roomID == "" {
		http.Error(res, "missing roomid", http.StatusBadRequest)
		return
	}

	body, err := bodyToString(req)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		http.Error(res, "cannot read body", http.StatusBadRequest)
		return
	}

	msg, err := convertRawJsonToMatrixMessage(*body, transformationType, mappingCode)
	if err != nil {
		log.Printf("failed to convert: %v", err)
		http.Error(res, "failed to convert message", http.StatusBadRequest)
		return
	}

	// If msg is nil, the event should be ignored (filtered user or unsupported event)
	if msg == nil {
		log.Printf("Event ignored (filtered user or unsupported event type)")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("ignored"))
		return
	}

	err = SendMessage(*msg, roomID)
	if err != nil {
		log.Printf("Matrix send error: %v", err)
		http.Error(res, "Matrix send error", 500)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("ok"))
}
