package schedule_controllers

import (
	"context"
	"io"
	event_broker "weather_notification/src/modules/notification_schedule/domain/event"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type StreamController struct {
	event_broker event_broker.EventBroker
	ctx          context.Context
}

func (controller *StreamController) StartStream(c *gin.Context) {
	controller.setSSEHeaders(c)

	c.Stream(func(w io.Writer) bool {
		err := controller.event_broker.Subscribe(c.GetString("AccountId"), func(message []byte) error {
			c.SSEvent("message", string(message))
			c.Writer.Flush()
			return nil
		})

		if err != nil {
			c.JSON(400, gin.H{"message": "Could not connect"})
			return false
		}

		return false
	})
}

func (controller *StreamController) setSSEHeaders(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
}

func NewStreamController(redisClient redis.Client, eventBroker event_broker.EventBroker) StreamController {
	return StreamController{
		event_broker: eventBroker,
		ctx:          context.Background(),
	}
}
