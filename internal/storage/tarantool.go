package storage

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Storage struct {
	conn *tarantool.Connection
	log  *logrus.Entry
}

// New создает новое подключение к Tarantool
func New(addr string) (*Storage, error) {
	conn, err := tarantool.Connect(addr, tarantool.Opts{
		User:          "admin",
		Pass:          "",
		Timeout:       5 * time.Second,
		Reconnect:     3 * time.Second,
		MaxReconnects: 3,
	})
	if err != nil {
		return nil, err
	}

	if _, err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Storage{
		conn: conn,
		log:  logrus.WithField("module", "tarantool"),
	}, nil
}

// Close закрывает соединение
func (s *Storage) Close() {
	if s.conn != nil {
		s.conn.Close()
		s.log.Info("Connection closed")
	}
}

// InitSpace создает пространство для хранения голосований
func (s *Storage) InitSpace() error {
	script := `
	box.schema.create_space('polls', {
		if_not_exists = true,
		format = {
			{name = 'id', type = 'string'},
			{name = 'creator', type = 'string'},
			{name = 'question', type = 'string'},
			{name = 'options', type = 'array'},
			{name = 'votes', type = 'map'}
		}
	})
	box.space.polls:create_index('primary', {
		type = 'hash',
		parts = {'id'}
	})
	`

	_, err := s.conn.Eval(script, []interface{}{})
	return err
}

// CreatePoll создает новое голосование
func (s *Storage) CreatePoll(pollID, creatorID, question string, options []string) error {
	s.log.WithFields(logrus.Fields{
		"pollID":   pollID,
		"creator":  creatorID,
		"question": question,
	}).Debug("Creating new poll")

	_, err := s.conn.Insert("polls", []interface{}{
		pollID,
		creatorID,
		question,
		options,
		map[string]int{},
	})

	if err != nil {
		s.log.WithError(err).Error("Failed to create poll")
	}
	return err
}

// Vote добавляет голос к указанному варианту
func (s *Storage) Vote(pollID, userID, option string) error {
	s.log.WithFields(logrus.Fields{
		"pollID": pollID,
		"userID": userID,
		"option": option,
	}).Debug("Processing vote")

	_, err := s.conn.Get("polls", pollID)
	if err != nil {
		s.log.WithError(err).Error("Poll not found")
		return err
	}

	_, err = s.conn.Update("polls", "primary", []interface{}{pollID}, []interface{}{
		[]interface{}{"+", 4, []interface{}{option, 1}},
	})

	if err != nil {
		s.log.WithError(err).Error("Failed to register vote")
	}
	return err
}

// GetResults возвращает текущие результаты голосования
func (s *Storage) GetResults(pollID string) (map[string]int, error) {
	s.log.WithField("pollID", pollID).Debug("Getting poll results")

	resp, err := s.conn.Select("polls", "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil {
		s.log.WithError(err).Error("Failed to get poll results")
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, nil
	}

	tuple := resp.Data[0].([]interface{})
	return tuple[4].(map[string]int), nil
}

// ClosePoll завершает голосование
func (s *Storage) ClosePoll(pollID string) error {
	s.log.WithField("pollID", pollID).Info("Closing poll")

	_, err := s.conn.Delete("polls", "primary", []interface{}{pollID})
	if err != nil {
		s.log.WithError(err).Error("Failed to close poll")
	}
	return err
}
