package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

type MockChannel struct {
	mock.Mock
}

func (m *MockChannel) QueueDeclare(
	name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table,
) (amqp.Queue, error) {
	argsCalled := m.Called(name, durable, autoDelete, exclusive, noWait, args)
	return amqp.Queue{Name: name}, argsCalled.Error(1)
}

func (m *MockChannel) Publish(
	exchange, key string, mandatory, immediate bool, msg amqp.Publishing,
) error {
	argsCalled := m.Called(exchange, key, mandatory, immediate, msg)
	return argsCalled.Error(0)
}

func (m *MockChannel) QueueBind(
	name, key, exchange string, noWait bool, args amqp.Table,
) error {
	argsCalled := m.Called(name, key, exchange, noWait, args)
	return argsCalled.Error(0)
}

func (m *MockChannel) Consume(
	queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table,
) (<-chan amqp.Delivery, error) {
	argsCalled := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	return make(<-chan amqp.Delivery), argsCalled.Error(1)
}

func (m *MockChannel) ExchangeDeclare(
	name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table,
) error {
	argsCalled := m.Called(name, kind, durable, autoDelete, internal, noWait, args)
	return argsCalled.Error(0)
}

func (m *MockChannel) Close() error {
	argsCalled := m.Called()
	return argsCalled.Error(0)
}

type MockConnection struct {
	mock.Mock
}

func (m *MockConnection) Channel() (*amqp.Channel, error) {
	argsCalled := m.Called()
	return argsCalled.Get(0).(*amqp.Channel), argsCalled.Error(1)
}

func (m *MockConnection) Close() error {
	argsCalled := m.Called()
	return argsCalled.Error(0)
}
