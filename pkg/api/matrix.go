package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// MatrixMessage represents the JSON payload for m.room.message
type InternalMatrixMessage struct {
	MsgType       string `json:"msgtype"`
	Body          string `json:"body"`
	Format        string `json:"format"`
	FormattedBody string `json:"formatted_body"`
}

type MatrixMessage struct {
	Plain string `json:"plain"`
	Html  string `json:"html"`
}

func SendMessage(matrixMessage MatrixMessage, roomID string) error {
	var AccessToken = os.Getenv("ACCESS_TOKEN")
	var Homeserver = os.Getenv("HOME_SERVER")

	matrixURL := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message?access_token=%s", Homeserver, roomID, AccessToken)

	matMsg := InternalMatrixMessage{
		MsgType:       "m.text",
		Body:          matrixMessage.Plain,
		Format:        "org.matrix.custom.html",
		FormattedBody: matrixMessage.Html,
	}
	jsonMsg, _ := json.Marshal(matMsg)

	resp, err := http.Post(matrixURL, "application/json", bytes.NewBuffer(jsonMsg))
	if err != nil {
		log.Printf("Matrix send error: %v", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		bb, _ := io.ReadAll(resp.Body)
		log.Printf("Matrix API error: %s", string(bb))
		return err
	}

	return nil
}
