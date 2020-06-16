package libs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Rhymen/go-whatsapp"
)

type wpResponse struct {
	Id          string
	RemoteJid   string
	SenderJid   string
	FromMe      bool
	Timestamp   uint64
	PushName    string
	MessageType string
}
type wpTextResponse struct {
	Info    wpResponse
	Text    string
	Context whatsapp.ContextInfo
}
type wpDataResponse struct {
	Info    wpResponse
	Data    string
	Caption string
	Context whatsapp.ContextInfo
}

type wpLocationResponse struct {
	Info      wpResponse
	Latitude  float64
	Longitude float64
	Name      string
	Address   string
	Url       string
	Context   whatsapp.ContextInfo
}
type myHandler struct {
	WebhookUrl string
}

func (myHandler) HandleError(err error) {
	if strings.Contains(err.Error(), "error processing data: received invalid data") || strings.Contains(err.Error(), "invalid string with tag 174") {
		return
	}
	fmt.Fprintf(os.Stderr, "%v", err)
}
func (h *myHandler) HandleLiveLocationMessage(message whatsapp.LiveLocationMessage) {
	if message.Info.FromMe {
		return
	} else {
		data, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(data))
		_, err = http.Post(h.WebhookUrl, "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func (h *myHandler) HandleLocationMessage(message whatsapp.LocationMessage) {
	if message.Info.FromMe {
		return
	} else {
		response := wpLocationResponse{
			Info: wpResponse{
				Id:          message.Info.Id,
				RemoteJid:   message.Info.RemoteJid,
				SenderJid:   message.Info.SenderJid,
				FromMe:      message.Info.FromMe,
				Timestamp:   message.Info.Timestamp,
				PushName:    message.Info.PushName,
				MessageType: "location",
			},
			Latitude:  message.DegreesLatitude,
			Longitude: message.DegreesLongitude,
			Name:      message.Name,
			Address:   message.Address,
			Url:       message.Url,
			Context:   message.ContextInfo,
		}

		data, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(data))
		_, err = http.Post(h.WebhookUrl, "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func (h *myHandler) HandleStickerMessage(message whatsapp.StickerMessage) {
	if message.Info.FromMe == true {
		return
	} else {
		data, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(data))
		// _, err = http.Post(h.WebhookUrl, "application/json", bytes.NewBuffer(data))
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
	}

}

func (h *myHandler) HandleTextMessage(message whatsapp.TextMessage) {
	if message.Info.FromMe == true {
		return
	} else {
		response := wpTextResponse{
			Info: wpResponse{
				Id:          message.Info.Id,
				RemoteJid:   message.Info.RemoteJid,
				SenderJid:   message.Info.SenderJid,
				FromMe:      message.Info.FromMe,
				Timestamp:   message.Info.Timestamp,
				PushName:    message.Info.PushName,
				MessageType: "text",
			},
			Text:    message.Text,
			Context: message.ContextInfo,
		}
		data, err := json.Marshal(response)
		fmt.Println(h.WebhookUrl)
		_, err = http.Post(h.WebhookUrl, "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Println(err)
			return

		}
	}

}

func (h *myHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	if message.Info.FromMe == true {
		return
	} else {
		response := wpDataResponse{
			Info: wpResponse{
				Id:          message.Info.Id,
				RemoteJid:   message.Info.RemoteJid,
				SenderJid:   message.Info.SenderJid,
				FromMe:      message.Info.FromMe,
				Timestamp:   message.Info.Timestamp,
				PushName:    message.Info.PushName,
				MessageType: "data/image",
			},
			Caption: message.Caption,
			Data:    "",
			Context: message.ContextInfo,
		}
		data, err := message.Download()
		if err != nil {
			fmt.Println(err)
			return
		}
		f, err1 := os.Create("./share/store/" + string(message.Info.Id))
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		defer f.Close()
		_, err = f.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		response.Data = string(message.Info.Id)
		payload, err2 := json.Marshal(response)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		fmt.Println(string(data))

		_, err = http.Post(h.WebhookUrl, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func (h *myHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	//fmt.Println(message.Info)
}

func (h *myHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
	if message.Info.FromMe == true {
		return
	} else {
		response := wpDataResponse{
			Info: wpResponse{
				Id:          message.Info.Id,
				RemoteJid:   message.Info.RemoteJid,
				SenderJid:   message.Info.SenderJid,
				FromMe:      message.Info.FromMe,
				Timestamp:   message.Info.Timestamp,
				PushName:    message.Info.PushName,
				MessageType: "data/image",
			},
			Caption: message.Caption,
			Data:    "",
			Context: message.ContextInfo,
		}
		data, err := message.Download()
		if err != nil {
			fmt.Println(err)
			return
		}

		f, err1 := os.Create("./share/store/" + string(message.Info.Id))
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		defer f.Close()
		_, err = f.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		response.Data = string(message.Info.Id)
		payload, err2 := json.Marshal(response)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		fmt.Println(string(data))

		_, err = http.Post(h.WebhookUrl, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func (h *myHandler) HandleAudioMessage(message whatsapp.AudioMessage) {
	if message.Info.FromMe == true {
		return
	} else {
		data, err := message.Download()
		if err != nil {
			fmt.Println(err)
			return
		}
		f, err1 := os.Create("./share/store/" + string(message.Info.Id))
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		defer f.Close()
		n, err := f.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("wrote %d bytes\n", n)
	}
}

func (h *myHandler) HandleJsonMessage(message string) {
	fmt.Println(message)
}

func (h *myHandler) HandleContactMessage(message whatsapp.ContactMessage) {
	fmt.Println(message)
}
