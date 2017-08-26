package eh

import (
	"github.com/looplab/eventhorizon"
	"github.com/eugeis/gee/enum"
	"context"
	"github.com/pkg/errors"
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"io"
	"html"
	"github.com/gorilla/schema"
	"github.com/looplab/eventhorizon/eventhandler/projector"
)

type AggregateInitializer struct {
	aggregateType    eventhorizon.AggregateType
	aggregateFactory func(id eventhorizon.UUID) eventhorizon.Aggregate
	modelFactory     func() interface{}
	commands         []enum.Literal
	events           []enum.Literal

	eventStore              eventhorizon.EventStore
	eventBus                eventhorizon.EventBus
	eventPublisher          eventhorizon.EventPublisher
	commandBus              eventhorizon.CommandBus
	repository              *eventhorizon.EventSourcingRepository
	commandHandler          *eventhorizon.AggregateCommandHandler
	projectorListener       DelegateEventHandler
	setupCallbacks          []func() error
	readRepos               func(name string) eventhorizon.ReadWriteRepo
	DefaultProjectorEnabled bool
	ProjectorRepo           eventhorizon.ReadRepo
}

func NewAggregateInitializer(aggregateType eventhorizon.AggregateType,
	aggregateFactory func(id eventhorizon.UUID) eventhorizon.Aggregate,
	modelFactory func() interface{},
	commands []enum.Literal, events []enum.Literal,
	projectorListener DelegateEventHandler,
	setupCallbacks []func() error, eventStore eventhorizon.EventStore, eventBus eventhorizon.EventBus,
	eventPublisher eventhorizon.EventPublisher, commandBus eventhorizon.CommandBus,
	readRepos func(name string) eventhorizon.ReadWriteRepo) (ret *AggregateInitializer) {
	ret = &AggregateInitializer{
		aggregateType:     aggregateType,
		aggregateFactory:  aggregateFactory,
		modelFactory:      modelFactory,
		commands:          commands,
		events:            events,
		projectorListener: projectorListener,
		setupCallbacks:    setupCallbacks,

		eventStore:              eventStore,
		eventBus:                eventBus,
		eventPublisher:          eventPublisher,
		commandBus:              commandBus,
		readRepos:               readRepos,
		DefaultProjectorEnabled: true,
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

	if o.DefaultProjectorEnabled {
		if err = o.registerProjector(); err != nil {
			return
		}
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

func (o *AggregateInitializer) registerProjector() (err error) {
	o.ProjectorRepo, err = o.RegisterProjector(o.projectorListener)
	return
}

func (o *AggregateInitializer) RegisterProjector(listener DelegateEventHandler) (ret eventhorizon.ReadRepo, err error) {
	projectorType := string(o.aggregateType)
	repo := o.readRepos(projectorType)
	projector := projector.NewEventHandler(NewProjector(projectorType, listener), repo)
	projector.SetModel(o.modelFactory)
	o.RegisterForAllEvents(projector)
	ret = repo
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
	Apply(event eventhorizon.Event, model interface{}) error
}

type AggregateBase struct {
	*eventhorizon.AggregateBase
	DelegateCommandHandler
	DelegateEventHandler
	Model interface{}
}

func (o *AggregateBase) HandleCommand(ctx context.Context, cmd eventhorizon.Command) error {
	return o.Execute(cmd, o.Model, o.AggregateBase)
}

func (o *AggregateBase) ApplyEvent(ctx context.Context, event eventhorizon.Event) error {
	return o.Apply(event, o.Model)
}

func NewAggregateBase(aggregateType eventhorizon.AggregateType, id eventhorizon.UUID,
	commandHandler DelegateCommandHandler, eventHandler DelegateEventHandler,
	model interface{}) *AggregateBase {
	return &AggregateBase{
		AggregateBase:          eventhorizon.NewAggregateBase(aggregateType, id),
		DelegateCommandHandler: commandHandler,
		DelegateEventHandler:   eventHandler,
		Model:                  model,
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

func IdNotDefined(currentId eventhorizon.UUID, aggregateType eventhorizon.AggregateType) error {
	return errors.New(fmt.Sprintf("Id not defined for aggregateType=%v", aggregateType))
}

func IdsDismatch(entityId eventhorizon.UUID, currentId eventhorizon.UUID, aggregateType eventhorizon.AggregateType) error {
	return errors.New(fmt.Sprintf("Dismatch entity id and current id, %v != %v, for aggregateType=%v",
		entityId, currentId, aggregateType))
}

func ValidateNewId(entityId eventhorizon.UUID, currentId eventhorizon.UUID, aggregateType eventhorizon.AggregateType) (ret error) {
	if len(entityId) > 0 {
		ret = EntityAlreadyExists(entityId, aggregateType)
	} else if len(currentId) == 0 {
		ret = IdNotDefined(currentId, aggregateType)
	}
	return
}

func ValidateIdsMatch(entityId eventhorizon.UUID, currentId eventhorizon.UUID, aggregateType eventhorizon.AggregateType) (ret error) {
	if len(entityId) == 0 {
		ret = EntityNotExists(currentId, aggregateType)
	} else if entityId != currentId {
		ret = IdsDismatch(entityId, currentId, aggregateType)
	}
	return
}

type HttpQueryHandler struct {
}

func (o *HttpQueryHandler) HandleResult(ret interface{}, err error, method string, w http.ResponseWriter, r *http.Request) {
	if err == nil {
		var js []byte
		if js, err = json.Marshal(ret); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

type HttpCommandHandler struct {
	Context    context.Context
	CommandBus eventhorizon.CommandBus
}

func NewHttpCommandHandler(context context.Context, commandBus eventhorizon.CommandBus) (ret *HttpCommandHandler) {
	ret = &HttpCommandHandler{
		Context:    context,
		CommandBus: commandBus,
	}
	return
}

func (o *HttpCommandHandler) HandleCommand(command eventhorizon.Command, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(command)
	defer r.Body.Close()
	if err == io.EOF {
		if err = r.ParseForm(); err == nil {
			err = schema.NewDecoder().Decode(command, r.Form)
		}
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't decode body to command %T %v because of %v", command, command, err),
			http.StatusBadRequest)
		return
	}

	if err := o.CommandBus.HandleCommand(o.Context, command); err != nil {
		http.Error(w, fmt.Sprintf("Can't execute command %T %v because of %v", command, command, err),
			http.StatusExpectationFailed)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	fmt.Fprintf(w, "Succefully executed command %T %v from %v", command, command, html.EscapeString(r.URL.Path))
}

type ProjectorEventHandler struct {
	DelegateEventHandler
	projectorType projector.Type
}

func NewProjector(projectorType string, eventHandler DelegateEventHandler) (ret *ProjectorEventHandler) {
	ret = &ProjectorEventHandler{
		projectorType:        projector.Type(projectorType),
		DelegateEventHandler: eventHandler,
	}
	return
}

func (o *ProjectorEventHandler) ProjectorType() projector.Type {
	return o.projectorType
}

func (o *ProjectorEventHandler) Project(ctx context.Context, event eventhorizon.Event, model interface{}) (ret interface{}, err error) {
	ret = model
	err = o.Apply(event, model)
	return
}
