package grpc

import (
	"context"

	pb "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/server/grpc/eventpb"
	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *CalendarService) CreateEvent(ctx context.Context, r *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	response := &pb.CreateEventResponse{}

	id, err := uuid.Parse(r.Event.ID)
	if err != nil {
		return response, err
	}

	userID, err := uuid.Parse(r.Event.UserID)
	if err != nil {
		return response, err
	}

	event := storage.Event{
		ID:          id,
		Title:       r.Event.Title,
		UserID:      userID,
		Description: r.Event.Description,
		StartDate:   r.Event.StartDate.AsTime(),
		EndDate:     r.Event.EndDate.AsTime(),
	}

	if err := s.app.InsertEvent(ctx, event); err != nil {
		return response, err
	}

	return response, nil
}

func (s *CalendarService) GetEventByID(ctx context.Context, r *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	response := &pb.GetEventResponse{}
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return response, err
	}

	e, err := s.app.GetEventByID(ctx, id)
	if err != nil {
		return response, err
	}

	response.Event = &pb.Event{
		ID:          e.ID.String(),
		Title:       e.Title,
		Description: e.Description,
		StartDate:   timestamppb.New(e.StartDate),
		EndDate:     timestamppb.New(e.EndDate),
	}

	return response, nil
}

func (s *CalendarService) UpdateEvent(ctx context.Context, r *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	response := &pb.UpdateEventResponse{}

	id, err := uuid.Parse(r.Event.ID)
	if err != nil {
		return response, err
	}

	userID, err := uuid.Parse(r.Event.UserID)
	if err != nil {
		return response, err
	}

	event := storage.Event{
		ID:          id,
		Title:       r.Event.Title,
		UserID:      userID,
		Description: r.Event.Description,
		StartDate:   r.Event.StartDate.AsTime(),
		EndDate:     r.Event.EndDate.AsTime(),
	}

	if err := s.app.UpdateEvent(ctx, event); err != nil {
		return response, err
	}

	return response, nil
}

func (s *CalendarService) DeleteEvent(ctx context.Context, r *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	response := &pb.DeleteEventResponse{}
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return response, err
	}

	if err := s.app.DeleteEvent(ctx, id); err != nil {
		return response, err
	}

	return response, nil
}

func (s *CalendarService) GetEvents(ctx context.Context, r *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	response := &pb.GetEventsResponse{}
	events, err := s.app.GetEventsForPeriod(ctx, r.StartDate.AsTime(), r.EndDate.AsTime())
	if err != nil {
		return response, err
	}

	response.Events = make([]*pb.Event, len(events))
	for i, e := range events {
		response.Events[i] = &pb.Event{
			ID:          e.ID.String(),
			Title:       e.Title,
			Description: e.Description,
			StartDate:   timestamppb.New(e.StartDate),
			EndDate:     timestamppb.New(e.EndDate),
		}
	}

	return response, nil
}
