package eh

import (
	"github.com/looplab/eventhorizon"
	"github.com/eugeis/gee/enum"
	"context"
	"github.com/pkg/errors"
	"fmt"
	"time"
)

type AggregateInitializer struct {
	aggregateType    eventhorizon.AggregateType
	aggregateFactory func(id eventhorizon.UUID) eventhorizon.Aggregate
	commands         []enum.Literal
	events           []enum.Literal

	eventStore     eventhorizon.EventStore
	eventBus       eventhorizon.EventBus
	eventPublisher eventhorizon.EventPublisher
	commandBus     eventhorizon.CommandBus
	repository     *eventhorizon.EventSourcingRepository
	commandHandler *eventhorizon.AggregateCommandHandler
	setupCallbacks []func() error
}

func NewAggregateInitializer(aggregateType eventhorizon.AggregateType,
	aggregateFactory func(id eventhorizon.UUID) eventhorizon.Aggregate, commands []enum.Literal, events []enum.Literal,
	setupCallbacks []func() error, eventStore eventhorizon.EventStore, eventBus eventhorizon.EventBus,
	eventPublisher eventhorizon.EventPublisher, commandBus eventhorizon.CommandBus) (ret *AggregateInitializer) {
	ret = &AggregateInitializer{
		aggregateType:    aggregateType,
		aggregateFactory: aggregateFactory,
		commands:         commands,
		events:           events,
		setupCallbacks:   setupCallbacks,

		eventStore:     eventStore,
		eventBus:       eventBus,
		eventPublisher: eventPublisher,
		commandBus:     commandBus,

	}
	return
}

func (o *AggregateInitializer) Setup() (err error) {
	//register aggregate factory
	eventhorizon.RegisterAggregate(o.aggregateFactory)

	if o.repository, err = eventhorizon.NewEventSourcingRepository(o.eventStore, o.eventBus); err != nil {
		return
	}

	if o.commandHandler, err = eventhorizon.NewAggregateCommandHandler(o.repository); err != nil {
		return
	}

	if err = o.registerCommands(); err != nil {
		return
	}

	if o.setupCallbacks != nil {
		for _, callback := range o.setupCallbacks {
			callback()
		}
	}
	return
}

func (o *AggregateInitializer) registerCommands() (err error) {
	for _, item := range o.commands {
		if err = o.commandHandler.SetAggregate(o.aggregateType, eventhorizon.CommandType(item.Name())); err != nil {
			return
		}
		if err = o.commandBus.SetHandler(o.commandHandler, eventhorizon.CommandType(item.Name())); err != nil {
			return
		}
	}
	return
}

func (o *AggregateInitializer) RegisterForAllEvents(handler eventhorizon.EventHandler) {
	for _, item := range o.events {
		o.eventBus.AddHandler(handler, eventhorizon.EventType(item.Name()))
	}
}

func (o *AggregateInitializer) RegisterForEvent(handler eventhorizon.EventHandler, event enum.Literal) {
	o.eventBus.AddHandler(handler, eventhorizon.EventType(event.Name()))
}

type AggregateStoreEvent interface {
	StoreEvent(eventhorizon.EventType, eventhorizon.EventData, time.Time) eventhorizon.Event
}

type DelegateCommandHandler interface {
	Execute(cmd eventhorizon.Command, entity interface{}, store AggregateStoreEvent) error
}

type DelegateEventHandler interface {
	Apply(event eventhorizon.Event, entity interface{}) error
}

type AggregateBase struct {
	*eventhorizon.AggregateBase
	DelegateCommandHandler
	DelegateEventHandler
	Entity interface{}
}

func (o *AggregateBase) HandleCommand(ctx context.Context, cmd eventhorizon.Command) error {
	return o.Execute(cmd, o.Entity, o.AggregateBase)
}

func (o *AggregateBase) ApplyEvent(ctx context.Context, event eventhorizon.Event) error {
	return o.Apply(event, o.Entity)
}

func NewAggregateBase(aggregateType eventhorizon.AggregateType, id eventhorizon.UUID,
	commandHandler DelegateCommandHandler,
	eventHandler DelegateEventHandler, entity interface{}) *AggregateBase {
	return &AggregateBase{
		AggregateBase:          eventhorizon.NewAggregateBase(aggregateType, id),
		DelegateCommandHandler: commandHandler,
		DelegateEventHandler:   eventHandler,
		Entity:                 entity,
	}
}

func CommandHandlerNotImplemented(commandType eventhorizon.CommandType) error {
	return errors.New(fmt.Sprintf("Handler not implemented for %v", commandType))
}

func EventHandlerNotImplemented(eventType eventhorizon.EventType) error {
	return errors.New(fmt.Sprintf("Handler not implemented for %v", eventType))
}

func EntityAlreadyExists(entityId eventhorizon.UUID, aggregateType eventhorizon.AggregateType) error {
	return errors.New(fmt.Sprintf("Entity already exists with id=%v and aggregateType=%v", entityId, aggregateType))
}

func EntityNotExists(entityId eventhorizon.UUID, aggregateType eventhorizon.AggregateType) error {
	return errors.New(fmt.Sprintf("Entity not exists with id=%v and aggregateType=%v", entityId, aggregateType))
}

func IdsDismatch(entityId eventhorizon.UUID, currentId eventhorizon.UUID, aggregateType eventhorizon.AggregateType) error {
	return errors.New(fmt.Sprintf("Dismatch entity id and current id, %v != %v, for aggregateType=%v",
		entityId, currentId, aggregateType))
}
