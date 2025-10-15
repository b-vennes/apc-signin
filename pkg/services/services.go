package services

import (
	"database/sql"
	"fmt"
	"log"

	"signinapc.app/api/pkg/events"
)

type EventService interface {
	Persist(event events.Event) (bool, error)
	NextSequence(stream string) (int64, error)
	ReadAll(stream string) ([]events.Event, error)
}

type sqliteEventService struct {
	db  *sql.DB
	ids events.IDs
}

func MakeSqliteEventService(db *sql.DB, ids events.IDs) EventService {
	return sqliteEventService{
		db,
		ids,
	}
}

func (s sqliteEventService) Persist(event events.Event) (bool, error) {
	meta := event.Meta(s.ids)
	body := event.Body()
	fmt.Printf("Persisting event %s to stream %s at sequence %d\n", body, meta.Stream, meta.SequenceNum)

	insertQuery := fmt.Sprintf(
		"INSERT INTO events(streamID, type, sequence, body) SELECT '%s', '%s', %d, '%s' WHERE NOT EXISTS(SELECT 1 FROM events WHERE streamID = '%s' AND sequence = %d)", meta.Stream,
		meta.Type,
		meta.SequenceNum,
		body,
		meta.Stream,
		meta.SequenceNum,
	)

	result, err := s.db.Exec(insertQuery)

	if err != nil {
		return false, fmt.Errorf("failed to execute event insert query: %s", err.Error())
	}

	changes, err := result.RowsAffected()

	if err != nil {
		return false, fmt.Errorf("failed to check if database was affected: %s", err.Error())
	}

	return changes > 0, nil
}

func (s sqliteEventService) NextSequence(stream string) (int64, error) {

	values, err := s.db.Query(fmt.Sprintf("SELECT sequence FROM events WHERE streamID = '%s' ORDER BY sequence DESC LIMIT 1", stream))

	if err != nil {
		return 0, fmt.Errorf("failed to request sequence rows from DB: %s", err.Error())
	}

	defer values.Close()

	var sequence int64

	hadNext := values.Next()

	if hadNext {
		err = values.Scan(&sequence)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to query for next sequence number: %s", err.Error())
	}

	if hadNext {
		sequence++
	}

	log.Println(sequence)

	return sequence, nil
}

type databaseEvent struct {
	meta events.EventMeta
	body string
}

func (d databaseEvent) Meta(_ events.IDs) events.EventMeta {
	return d.meta
}

func (d databaseEvent) Body() string {
	return d.body
}

func (s sqliteEventService) ReadAll(stream string) ([]events.Event, error) {
	values, err := s.db.Query(fmt.Sprintf("SELECT * FROM events WHERE streamID = '%s'", stream))

	if err != nil {
		return nil, fmt.Errorf("failed to request stream rows from DB: %s", err.Error())
	}

	defer values.Close()

	response := []events.Event{}

	for values.Next() {
		var (
			streamID  string
			eventType string
			sequence  int64
			body      string
		)

		err := values.Scan(&streamID, &eventType, &sequence, &body)

		if err != nil {
			return nil, fmt.Errorf("failed to stream all rows from DB: %s", err.Error())
		}

		record := databaseEvent{
			meta: events.EventMeta{
				Stream:      streamID,
				Type:        eventType,
				SequenceNum: sequence,
			},
			body: body,
		}

		response = append(response, record)
	}

	return response, nil
}
