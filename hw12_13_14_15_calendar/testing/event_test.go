package testing_test

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/server/grpc/eventpb"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("gRpc adding a new event", func() {
	BeforeEach(func() {
		ctx := context.Background()
		connectionString := "postgres://pguser:pgpwd@localhost:5432/calendar?sslmode=disable"
		db, err := sqlx.Open("postgres", connectionString)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		_, _ = db.ExecContext(ctx, "DELETE FROM events;")

		fmt.Println("All event were removed")
	})

	Context("when event is correct", func() {
		It("adds a new event successfully", func() {
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "0.0.0.0:50051",
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			client := eventpb.NewCalendarClient(conn)

			id := uuid.New().String()
			e, err := client.CreateEvent(ctx, &eventpb.CreateEventRequest{
				Event: &eventpb.Event{
					ID:          id,
					Title:       "title",
					Description: "description",
					UserID:      uuid.New().String(),
					StartDate:   timestamppb.New(time.Now().AddDate(0, 0, -3)),
					EndDate:     timestamppb.New(time.Now().AddDate(0, 0, 1)),
				},
			})

			Expect(e).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())

			e1, err := client.GetEventByID(ctx, &eventpb.GetEventRequest{
				ID: id,
			})
			Expect(err).To(BeNil())
			Expect(e1).NotTo(BeNil())
		})
	})

	Context("when id of event is incorrect GUID", func() {
		It("returns error", func() {
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "0.0.0.0:50051",
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			client := eventpb.NewCalendarClient(conn)

			_, err = client.CreateEvent(ctx, &eventpb.CreateEventRequest{
				Event: &eventpb.Event{
					ID:          "wrong GUID",
					Title:       "title",
					Description: "description",
					UserID:      uuid.New().String(),
					StartDate:   timestamppb.New(time.Now().AddDate(0, 0, -3)),
					EndDate:     timestamppb.New(time.Now().AddDate(0, 0, 1)),
				},
			})

			Expect(err).NotTo(BeNil())
		})
	})

	Context("when there are events for today", func() {
		It("returns only events for one day", func() {
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "0.0.0.0:50051",
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			client := eventpb.NewCalendarClient(conn)

			_, err = client.CreateEvent(ctx, &eventpb.CreateEventRequest{
				Event: &eventpb.Event{
					ID:          uuid.New().String(),
					Title:       "title",
					Description: "description",
					UserID:      uuid.New().String(),
					StartDate:   timestamppb.New(time.Now().AddDate(0, 0, -1)),
					EndDate:     timestamppb.New(time.Now().AddDate(0, 0, 1)),
				},
			})
			Expect(err).To(BeNil())
			_, err = client.CreateEvent(ctx, &eventpb.CreateEventRequest{
				Event: &eventpb.Event{
					ID:          uuid.New().String(),
					Title:       "title1",
					Description: "description1",
					UserID:      uuid.New().String(),
					StartDate:   timestamppb.New(time.Now()),
					EndDate:     timestamppb.New(time.Now().AddDate(0, 0, 1)),
				},
			})
			Expect(err).To(BeNil())
			_, err = client.CreateEvent(ctx, &eventpb.CreateEventRequest{
				Event: &eventpb.Event{
					ID:          uuid.New().String(),
					Title:       "title2",
					Description: "description2",
					UserID:      uuid.New().String(),
					StartDate:   timestamppb.New(time.Now()),
					EndDate:     timestamppb.New(time.Now().AddDate(0, 0, 1)),
				},
			})
			Expect(err).To(BeNil())

			now := time.Now()
			year, month, day := now.Date()
			begin := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
			end := begin.AddDate(0, 0, 1)

			response, err := client.GetEvents(ctx, &eventpb.GetEventsRequest{
				StartDate: timestamppb.New(begin),
				EndDate:   timestamppb.New(end),
			})

			Expect(err).To(BeNil())
			Expect(response.Events).Should(HaveLen(2))
		})
	})
})
