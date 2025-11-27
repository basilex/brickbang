package rbac

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type IRbacListener interface {
	Start()
	Stop()
}

type RbacListener struct {
	connStr string
	channel string
	cache   IRbacCache
	logger  *log.Logger

	stop    chan struct{}
	stopped chan struct{}
}

func NewRbacListener(connStr, channel string, cache IRbacCache, logger *log.Logger) IRbacListener {
	return &RbacListener{
		connStr: connStr,
		channel: channel,
		cache:   cache,
		logger:  logger,

		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

func (rcv *RbacListener) Start() {
	go rcv.run()
}

func (rcv *RbacListener) Stop() {
	close(rcv.stop)
	<-rcv.stopped
}

func (rcv *RbacListener) run() {
	defer close(rcv.stopped)

	for {
		select {
		case <-rcv.stop:
			rcv.logger.Println("Rbac listener stopped")
			return
		default:
		}

		ctx := context.Background()

		conn, err := pgx.Connect(ctx, rcv.connStr)
		if err != nil {
			rcv.logger.Printf("connection error: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		_, err = conn.Exec(ctx, "LISTEN "+rcv.channel)
		if err != nil {
			rcv.logger.Printf("LISTEN error: %v", err)

			conn.Close(ctx)
			time.Sleep(2 * time.Second)

			continue
		}

		rcv.logger.Printf("LISTEN on '%s' ready", rcv.channel)

		for {
			select {
			case <-rcv.stop:
				conn.Close(ctx)
				return
			default:
			}

			notification, err := conn.WaitForNotification(ctx)
			if err != nil {
				rcv.logger.Printf("notification error: %v", err)
				conn.Close(ctx)
				break
			}

			if notification != nil {
				rcv.handlePayload(notification.Payload)
			}
		}

		time.Sleep(2 * time.Second)
	}
}

func (rcv *RbacListener) handlePayload(payload string) {
	// payload format: "op:entity:id"
	parts := strings.Split(payload, ":")

	if len(parts) != 3 {
		rcv.logger.Printf("invalid payload: %s", payload)
		return
	}

	entity := parts[1]
	id := parts[2]

	switch entity {
	case "role":
		rcv.cache.InvalidateRole(id)
	case "grant":
		rcv.cache.InvalidateGrant(id)
	case "role_grant":
		rcv.cache.InvalidateRoleGrant(id)
	case "user_role":
		rcv.cache.InvalidateUserRole(id)
	default:
		rcv.logger.Printf("unknown entity: %s", entity)
		return
	}

	rcv.logger.Printf("Rbac event: %s %s", entity, id)
}
