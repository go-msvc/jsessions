package memsessions

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-msvc/errors"
	"github.com/go-msvc/jsessions"
	"github.com/satori/uuid"
)

func New() jsessions.ISessions {
	ss := &sessions{
		byID:       map[string]jsessions.ISession{},
		closedChan: make(chan *session),
	}
	go func(ss *sessions) {
		for {
			s, ok := <-ss.closedChan
			if !ok {
				break
			}
			delete(ss.byID, s.id)
		}
	}(ss)
	return ss
}

type sessions struct {
	sync.Mutex
	byID       map[string]jsessions.ISession
	closedChan chan *session
}

//call Get("") to start a new session
//call Get("...") to get existing session or nil if not exist
func (ss *sessions) Get(id string) jsessions.ISession {
	ss.Lock()
	defer ss.Unlock()
	if s, ok := ss.byID[id]; ok {
		return s
	}
	if id != "" {
		return nil //does not exist
	}
	//new session
	s := &session{
		ss:        ss,
		id:        uuid.NewV1().String(),
		startTime: time.Now(),
		lastTime:  time.Now(),
		closed:    false,
		data:      map[string]interface{}{},
	}
	ss.byID[s.id] = s
	return s
}

type session struct {
	sync.Mutex
	ss        *sessions
	id        string
	startTime time.Time
	lastTime  time.Time
	closed    bool
	data      map[string]interface{}
}

func (s *session) ID() string { return s.id }

func (s *session) Close() {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return
	}
	s.closed = true
}

func (s *session) Set(name string, newValue interface{}) (savedValue interface{}, err error) {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return nil, errors.Errorf("session[%s].Set(%s) after closing", s.id, name)
	}
	s.data[name] = newValue
	return newValue, nil
}

func (s *session) SetString(name string, newValue string) (savedValue string, err error) {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return "", errors.Errorf("session[%s].Set(%s) after closing", s.id, name)
	}
	s.data[name] = newValue
	return newValue, nil
}
func (s *session) SetInt(name string, newValue int) (savedValue int, err error) {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return 0, errors.Errorf("session[%s].Set(%s) after closing", s.id, name)
	}
	s.data[name] = newValue
	return newValue, nil
}
func (s *session) SetBool(name string, newValue bool) (savedValue bool, err error) {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return false, errors.Errorf("session[%s].Set(%s) after closing", s.id, name)
	}
	s.data[name] = newValue
	return newValue, nil
}

func (s *session) Get(name string) (interface{}, error) {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return nil, errors.Errorf("session[%s].Get(%s) after closing", s.id, name)
	}
	if value, ok := s.data[name]; ok {
		return value, nil
	}
	return nil, nil
}

func (s *session) GetString(name string) (string, bool) {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return "", false
	}
	if value, ok := s.data[name]; ok {
		if strValue, ok := value.(string); ok {
			return strValue, true
		}
		return fmt.Sprintf("%v", value), ok
	}
	return "", false
}

func (s *session) GetInt(name string) (int, bool) {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return 0, false
	}
	if value, ok := s.data[name]; ok {
		if intValue, ok := value.(int); ok {
			return intValue, true
		}
		strValue := fmt.Sprintf("%v", value)
		if intValue, err := strconv.Atoi(strValue); err == nil {
			return intValue, true
		}
	}
	return 0, false
}

func (s *session) GetBool(name string) (bool, bool) {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return false, false
	}
	if value, ok := s.data[name]; ok {
		if boolValue, ok := value.(bool); ok {
			return boolValue, true
		}
		sv := strings.ToLower(fmt.Sprintf("%v", value))
		if sv == "true" || sv == "yes" || sv == "1" {
			return true, true
		}
		if sv == "false" || sv == "no" || sv == "0" {
			return true, true
		}
	}
	return false, false
}

func (s *session) Data() map[string]interface{} {
	return s.data
}

func (s *session) Save() error {
	return nil
}
