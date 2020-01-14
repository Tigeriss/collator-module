package session

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/sessions"
)

const (
	cookieName = "collator-sessions"
	defaultKey = "default"
)

var (
	key         = []byte("f9abac83a0db40d2238d09fc22d0fce4")
	sessionPath = path.Join(".", "collator-sessions")
	store = sessions.NewCookieStore(key)
)

func init() {
	err := os.MkdirAll(sessionPath, 0755)
	if err != nil {
		log.Panicf("unable to create session directory: %s\n", err)
	}
}

func Load(val interface{}, writer http.ResponseWriter, request *http.Request) error {
	// get a session
	sess, err := store.Get(request, cookieName)
	if err != nil {
		return fmt.Errorf("unable to init a default session storage: %s", err)
	}

	def, ok := sess.Values[defaultKey].([]byte)
	if !ok {
		// that's ok, it's new session
		return nil
	}

	// otherwise deserialize it
	err = json.Unmarshal(def, val)
	if err != nil {
		// it's broken? so clear it in that case
		log.Printf("unable to unmarshal the session with %s ssid: %s", sess.ID, err)
		delete(sess.Values, defaultKey)
		err = sess.Save(request, writer)
		if err != nil {
			return fmt.Errorf("unable to delete the broken session with %s ssid: %s", sess.ID, err)
		}
	}

	return nil
}

func Save(val interface{}, writer http.ResponseWriter, request *http.Request) error {
	sess, err := store.Get(request, cookieName)
	if err != nil {
		return fmt.Errorf("unable to init a default session storage: %s", err)
	}

	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("unable to unmarshalize the session with %s ssid: %s\n", sess.ID, err)
	}

	sess.Values[defaultKey] = data
	err = sess.Save(request, writer)
	if err != nil {
		return fmt.Errorf("unable to store the session with %s ssid: %s", sess.ID, err)
	}

	return nil
}
