package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Message struct {
	Alias       string       `json:"alias"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
	Sections    []Section    `json:"sections"`
}

type Attachment struct {
	Title     string   `json:"title"`
	TitleLink string   `json:"title_link"`
	Text      *string  `json:"text"`
	ImageURL  *string  `json:"image_url"`
	Color     string   `json:"color"`
	Fields    []Field  `json:"fields"`
	MrkdownIn []string `json:"mrkdown_in"` // note spelling from your sample
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Section struct {
	ActivityTitle    string `json:"activityTitle"`
	ActivitySubtitle string `json:"activitySubtitle"`
}

func GlitchTipHandler(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var msg Message
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "cannot read body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &msg); err != nil || msg.Text == "" {
		http.Error(res, "invalid or missing field", http.StatusBadRequest)
		return
	}

	matMsg := MatrixMessage{
		Body: msg.Text,
		Html: fmt.Sprintf(`<b>%s:</b> %s (<a href="%s">%s<a/>)`, msg.Text, msg.Attachments[0].Title, msg.Attachments[0].TitleLink, "View Issue"),
	}
	err = SendMessage(matMsg)
	if err != nil {
		log.Printf("Matrix send error: %v", err)
		http.Error(res, "Matrix send error", 500)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("ok"))

	/*
		str, _ := json.Marshal(matMsg)
		res.Write([]byte(str))
	*/
}
