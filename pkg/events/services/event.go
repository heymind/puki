package services

import (
	"github.com/lantu-dev/puki/pkg/events/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type EventService struct {
	db *gorm.DB
}

func NewEventService(db *gorm.DB) *EventService {
	return &EventService{db: db}
}

type GetEventsListReq struct {
	EventIDs []int64
}
type GetEventsListRes []struct {
	ID          int64
	Organizer   string
	Title       string
	Description string
	ImageUrl    string
	StartedAt   time.Time
	EndedAt     time.Time
	Location    string
	EventType   uint16
}

func (s EventService) GetEventsList(r *http.Request, req *GetEventsListReq, res *GetEventsListRes) (err error) {
	err = s.db.Model(&models.Event{}).Where(req.EventIDs).Find(res).Error

	return
}

type GetEventMoreInfoReq struct {
	EventID int64
}
type GetEventMoreInfoRes struct {
	Schedules []struct {
		Title             string
		StartedAt         time.Time
		EndedAt           time.Time
		TalkerName        string
		TalkerTitle       string
		TalkerAvatarURL   string
		TalkerDescription string
	}
	Hackathon struct {
		Steps string
	}
}

func (s EventService) GetEventMoreInfo(r *http.Request, req *GetEventMoreInfoReq, res *GetEventMoreInfoRes) (err error) {
	var Event struct {
		EventType uint16
	}

	s.db.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Model(&models.Event{}).First(&Event, req.EventID).Error; err != nil {
			return
		}
		var target *gorm.DB
		switch Event.EventType {
		case models.EventTypeSalon:
			fallthrough

		case models.EventTypeLecture:
			target = tx.Model(&models.Schedule{}).Where(&models.Schedule{EventID: req.EventID})
			if err = target.Find(&res.Schedules).Error; err != nil {
				return
			}

		case models.EventTypeHackathon:
			target = tx.Model(&models.Hackathon{}).Where(&models.Hackathon{EventID: req.EventID})
			if err = target.First(&res.Hackathon).Error; err != nil {
				return
			}

		case models.EventTypeOther:

		case models.EventTypeNull:

		default:
		}

		return
	})

	return
}
