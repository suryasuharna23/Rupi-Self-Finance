package telegram

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/suryasuharna23/rupi/internal/conversation"
)

// Handler adalah sisi Conversation Service yang dipanggil adapter.
type Handler interface {
	HandleMessage(ctx context.Context, msg conversation.Incoming) ([]conversation.Reply, error)
}

type Adapter struct {
	client  *Client
	handler Handler
}

func NewAdapter(client *Client, handler Handler) *Adapter {
	return &Adapter{client: client, handler: handler}
}

// Poll menjalankan long polling sampai ctx dibatalkan.
func (a *Adapter) Poll(ctx context.Context) error {
	if err := a.client.DeleteWebhook(ctx); err != nil {
		return err
	}
	var offset int64
	for ctx.Err() == nil {
		updates, err := a.client.GetUpdates(ctx, offset)
		if err != nil {
			if ctx.Err() == nil {
				slog.Error("getUpdates gagal", "err", err)
				time.Sleep(3 * time.Second)
			}
			continue
		}
		for _, u := range updates {
			offset = u.UpdateID + 1
			a.process(ctx, u)
		}
	}
	return nil
}

// WebhookHandler memverifikasi secret_token sebelum memproses update (SEC-01).
func (a *Adapter) WebhookHandler(secretToken string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
		if subtle.ConstantTimeCompare([]byte(got), []byte(secretToken)) != 1 {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		var u Update
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&u); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		a.process(r.Context(), u)
		w.WriteHeader(http.StatusOK)
	})
}

func (a *Adapter) process(ctx context.Context, u Update) {
	msg, ok := toIncoming(u)
	if !ok {
		return
	}
	if u.CallbackQuery != nil {
		if err := a.client.AnswerCallback(ctx, u.CallbackQuery.ID); err != nil {
			slog.Warn("answerCallbackQuery gagal", "err", err)
		}
	}
	replies, err := a.handler.HandleMessage(ctx, msg)
	if err != nil {
		slog.Error("handle message gagal", "err", err, "message_id", msg.MessageID)
		replies = []conversation.Reply{{Text: "Maaf, ada gangguan. Coba lagi sebentar lagi ya."}}
	}
	for _, r := range replies {
		if err := a.client.SendMessage(ctx, msg.ChatID, r.Text, toButtons(r.Buttons)); err != nil {
			slog.Error("sendMessage gagal", "err", err)
		}
	}
}

func toIncoming(u Update) (conversation.Incoming, bool) {
	switch {
	case u.Message != nil && u.Message.From != nil && u.Message.Text != "":
		m := u.Message
		return conversation.Incoming{
			Channel:    "telegram",
			MessageID:  strconv.FormatInt(m.Chat.ID, 10) + ":" + strconv.FormatInt(m.MessageID, 10),
			ChatID:     m.Chat.ID,
			ExternalID: m.From.ID,
			Text:       m.Text,
			ReceivedAt: time.Unix(m.Date, 0),
		}, true
	case u.CallbackQuery != nil && u.CallbackQuery.Message != nil:
		c := u.CallbackQuery
		return conversation.Incoming{
			Channel:      "telegram",
			MessageID:    "cb:" + c.ID,
			ChatID:       c.Message.Chat.ID,
			ExternalID:   c.From.ID,
			CallbackData: c.Data,
			ReceivedAt:   time.Now(),
		}, true
	}
	return conversation.Incoming{}, false
}

func toButtons(rows [][]conversation.Button) [][]InlineButton {
	out := make([][]InlineButton, 0, len(rows))
	for _, row := range rows {
		r := make([]InlineButton, 0, len(row))
		for _, b := range row {
			r = append(r, InlineButton{Text: b.Text, CallbackData: b.Data})
		}
		out = append(out, r)
	}
	return out
}
