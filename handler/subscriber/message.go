package subscriber

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
)

type MessageService interface {
	SubscribeMessage() http.Handler
}

type messageHandler struct {
	messageApplication adapter.MessageApplication
	logger             adapter.CompositeLogger
}

func NewMessageHandler(
	messageApplication adapter.MessageApplication,
	logger adapter.CompositeLogger) MessageService {
	return &messageHandler{
		messageApplication: messageApplication,
		logger:             logger,
	}
}

type PubsubMessageParams struct {
	Message struct {
		Data string `json:"data"`
	} `json:"message"`
}

type MessageParams struct {
	ID        string `json:"id"`
	RoomID    string `json:"roomId"`
	FromID    string `json:"fromId"`
	ToID      string `json:"toId"`
	Text      string `json:"text"`
	ImageURL  string `json:"imageUrl"`
	FileURL   string `json:"fileUrl"`
	CreatedAt int64  `json:"createdAt"`
}

func (h *messageHandler) SubscribeMessage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		b, err := ioutil.ReadAll(r.Body)
		defer func() {
			_ = r.Body.Close()
		}()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		pubsubParams := &PubsubMessageParams{}

		if err := json.Unmarshal(b, pubsubParams); err != nil {
			h.logger.Error().With(ctx).Printf("Error json parse: %v", err)
			http.Error(w, err.Error(), 400)
			return
		}

		data, err := base64.StdEncoding.DecodeString(pubsubParams.Message.Data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		messageParams := &MessageParams{}

		if err := json.Unmarshal(data, messageParams); err != nil {
			h.logger.Error().With(ctx).Printf("Error json parse: %v", err)
			http.Error(w, err.Error(), 400)
			return
		}

		imageURL, _ := url.Parse(messageParams.ImageURL)
		fileURL, _ := url.Parse(messageParams.FileURL)

		if err := h.messageApplication.Create(
			ctx,
			domain.MessageID(messageParams.ID),
			domain.MessageRoomID(messageParams.RoomID),
			messageParams.FromID,
			messageParams.ToID,
			adapter.MessageParams{
				Text:     messageParams.Text,
				ImageURL: imageURL,
				FileURL:  fileURL,
			},
			time.Now()); err != nil {
			h.logger.Error().With(ctx).Printf("Error create message: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}
	})
}
