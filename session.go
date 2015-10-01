package martiniframework

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	ActiveSession *SessionConnector
)

type SessionDriver interface {
	Connect(params map[string]interface{}) error
	Get(sid string) (string, error)
	Del(sid string) error
	Set(sid string, data string, expiry time.Duration) error
}

type SessionConnector struct {
	Address       string
	Password      string
	DatabaseId    int64
	SessionLength int64
	Driver        SessionDriver
}

type SessionModel struct {
	SessionId   string `json:"session_id" db:"session_id"`
	UserId      int64  `json:"user_id" db:"user_id"`
	Expires     int64  `json:"expiry_time" db:"expiry_time"`
	SessionData []byte `json:"session_data" db:"session_data"`
}

func (s *SessionConnector) Connect() error {
	return s.Driver.Connect(map[string]interface{}{
		"Address":    s.Address,
		"Password":   s.Password,
		"DatabaseId": s.DatabaseId,
	})
}

func (s *SessionConnector) CreateSession(uid int64) (SessionModel, error) {
	hn, _ := os.Hostname()
	sid := fmt.Sprintf("%d-%s", time.Now().Unix(), Md5hash(fmt.Sprintf("%d.%s.%d", time.Now().Unix(), hn, time.Now().UnixNano())))
	sm := SessionModel{
		SessionId: sid,
		UserId:    uid,
	}
	log.Printf("CreateSession(%d): key %s = %v", uid, sid, sm)
	err := s.StoreFreshenSession(sm)
	if err != nil {
		log.Print(err.Error())
		return SessionModel{}, err
	}
	return sm, nil
}

func (s *SessionConnector) GetSession(sid string) (SessionModel, error) {
	log.Printf("GetSession(%s)", sid)
	val, err := s.Driver.Get(sid)
	if err != nil {
		return SessionModel{}, err
	}
	var m SessionModel
	err = json.Unmarshal([]byte(val), &m)
	if err != nil {
		return SessionModel{}, err
	}
	return m, nil
}

func (s *SessionConnector) ExpireSession(sid string) error {
	log.Printf("ExpireSession(%s)", sid)
	return s.Driver.Del(sid)
}

func (s *SessionConnector) StoreFreshenSession(sm SessionModel) error {
	log.Printf("StoreFreshenSession(): %v", sm)
	sm.Expires = s.SessionLength
	b, err := json.Marshal(sm)
	if err != nil {
		return err
	}
	return s.Driver.Set(sm.SessionId, string(b), time.Duration(sm.Expires)*time.Second)
}

func TokenAuthFunc(sid string) (bool, SessionModel) {
	if sid == "" {
		return false, SessionModel{}
	}
	s, err := ActiveSession.GetSession(sid)
	log.Printf("TokenAuthFunc(): %s returned %v", sid, s)
	if err != nil {
		return false, SessionModel{}
	}
	if s.SessionId != "" {
		return true, s
	}
	return false, SessionModel{}
}
