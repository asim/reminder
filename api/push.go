package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"log"
)

type PushSubscription struct {
	Endpoint string                 `json:"endpoint"`
	Keys     map[string]interface{} `json:"keys"`
}

var pushFile = ReminderPath("push_subscriptions.json")
var pushMtx sync.RWMutex
var pushSubscriptions = make(map[string]PushSubscription)

// VAPID keys (generate your own for production)
var VAPIDPublicKey string
var VAPIDPrivateKey string
var VAPIDEmail = "mailto:admin@reminder.local"

func LoadOrGenerateVAPIDKeys() error {
	dir := ReminderPath("keys")
	privPath := filepath.Join(dir, "vapid_private.pem")
	pubPath := filepath.Join(dir, "vapid_public.b64")

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Check if keys exist
	privExists := false
	pubExists := false
	if _, err := os.Stat(privPath); err == nil {
		privExists = true
	}
	if _, err := os.Stat(pubPath); err == nil {
		pubExists = true
	}

	if !privExists || !pubExists {
		priv, pub, err := webpush.GenerateVAPIDKeys()
		if err != nil {
			return err
		}
		if err := os.WriteFile(pubPath, []byte(pub), 0600); err != nil {
			return err
		}
		if err := os.WriteFile(privPath, []byte(priv), 0600); err != nil {
			return err
		}
	}
	// Load public key
	pub, err := os.ReadFile(pubPath)
	if err != nil {
		return err
	}
	VAPIDPublicKey = string(pub)
	// Debug: print decoded length and first byte
	decoded, err := decodeBase64URL(VAPIDPublicKey)
	if err == nil {
		log.Printf("[VAPID] Decoded public key length: %d, first byte: %d", len(decoded), decoded[0])
	} else {
		log.Printf("[VAPID] Failed to decode public key: %v", err)
	}
	// Load private key
	priv, err := os.ReadFile(privPath)
	if err != nil {
		return err
	}
	VAPIDPrivateKey = string(priv)
	return nil
}

// Helper to decode base64url
func decodeBase64URL(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	return base64.RawURLEncoding.DecodeString(s)
}

func LoadPushSubscriptions() error {
	f, err := os.Open(pushFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&pushSubscriptions)
}

func SavePushSubscriptions() error {
	pushMtx.RLock()
	defer pushMtx.RUnlock()
	f, err := os.Create(pushFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(pushSubscriptions)
}

func AddPushSubscription(sub PushSubscription) error {
	pushMtx.Lock()
	pushSubscriptions[sub.Endpoint] = sub
	pushMtx.Unlock()
	return SavePushSubscriptions()
}

func RemovePushSubscription(endpoint string) error {
	pushMtx.Lock()
	delete(pushSubscriptions, endpoint)
	pushMtx.Unlock()
	return SavePushSubscriptions()
}

func ListPushSubscriptions() []PushSubscription {
	pushMtx.RLock()
	defer pushMtx.RUnlock()
	result := make([]PushSubscription, 0, len(pushSubscriptions))
	for _, sub := range pushSubscriptions {
		result = append(result, sub)
	}
	return result
}

func RegisterPushRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/push/subscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var sub PushSubscription
		b, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(b, &sub); err != nil || sub.Endpoint == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := AddPushSubscription(sub); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/api/push/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req struct{ Endpoint string }
		b, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(b, &req); err != nil || req.Endpoint == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := RemovePushSubscription(req.Endpoint); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Expose VAPID public key to frontend
	mux.HandleFunc("/api/push/vapidPublicKey", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(VAPIDPublicKey))
	})
}

func SendPushNotification(sub PushSubscription, payload string) error {
	subscription := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.Keys["p256dh"].(string),
			Auth:   sub.Keys["auth"].(string),
		},
	}
	resp, err := webpush.SendNotification([]byte(payload), subscription, &webpush.Options{
		Subscriber:      VAPIDEmail,
		VAPIDPublicKey:  VAPIDPublicKey,
		VAPIDPrivateKey: VAPIDPrivateKey,
		TTL:             60,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func SendPushToAll(payload string) []string {
	subs := ListPushSubscriptions()
	errors := []string{}
	for _, sub := range subs {
		err := SendPushNotification(sub, payload)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to send push to %s: %v", sub.Endpoint, err)
			log.Println(errMsg)
			errors = append(errors, errMsg)
		}
		// avoid rate limits
		time.Sleep(100 * time.Millisecond)
	}
	return errors
}
