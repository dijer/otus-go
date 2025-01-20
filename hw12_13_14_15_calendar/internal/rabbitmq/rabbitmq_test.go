package rabbitmq

import (
	"testing"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRabbitClient(t *testing.T) {
	mockConn := new(MockConnection)
	mockChan := new(MockChannel)
	log := logger.New("INFO")

	mockConn.On("Channel").Return(mockChan, nil)
	mockChan.On("QueueDeclare", "testQueue", true, false, false, false, nil).Return(amqp.Queue{Name: "testQueue"}, nil)
	mockChan.On("Publish", "", "testQueue", false, false, mock.Anything).Return(nil)
	mockChan.On("QueueBind", "testQueue", "", "testExchange", false, nil).Return(nil)
	mockChan.On("ExchangeDeclare", "testExchange", "fanout", true, false, false, false, nil).Return(nil)
	mockChan.On("Close").Return(nil)
	mockConn.On("Close").Return(nil)

	client, err := New("amqp://guest:guest@localhost:5672/", log)
	require.Nil(t, err)
	require.NotNil(t, client)

	err = client.ExchangeDeclare("testExchange")
	require.Nil(t, err)

	err = client.QueueDeclare("testQueue")
	require.Nil(t, err)

	err = client.Publish("testExchange", []byte("testMessage"))
	require.Nil(t, err)

	err = client.QueueBind("testExchange")
	require.Nil(t, err)

	client.Close()
}
