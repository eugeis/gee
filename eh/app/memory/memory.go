package memory

import (
	commandbus "github.com/looplab/eventhorizon/commandbus/local"
	eventbus "github.com/looplab/eventhorizon/eventbus/local"
	eventstore "github.com/looplab/eventhorizon/eventstore/memory"
	eventpublisher "github.com/looplab/eventhorizon/publisher/local"
	repo "github.com/looplab/eventhorizon/repo/memory"
	"github.com/eugeis/gee/eh/app"
	"github.com/eugeis/gee/eh"
	"github.com/looplab/eventhorizon"
)

func NewAppMemory(productName string, appName string, secure bool) *app.AppBase {
	// Create the event store.
	eventStore := eventstore.NewEventStore()

	// Create the event bus that distributes events.
	eventBus := eventbus.NewEventBus()
	eventPublisher := eventpublisher.NewEventPublisher()
	eventBus.SetPublisher(eventPublisher)

	// Create the command bus.
	commandBus := commandbus.NewCommandBus()

	repos := make(map[string]eventhorizon.ReadWriteRepo)
	readRepos := func(name string, factory func() interface{}) (ret eventhorizon.ReadWriteRepo) {
		if item, ok := repos[name]; !ok {
			ret = &eh.ReadWriteRepoDelegate{Factory: func() (ret eventhorizon.ReadWriteRepo, err error) {
				return repo.NewRepo(), nil
			}}
			repos[name] = ret
		} else {
			ret = item
		}
		return
	}
	return app.NewAppBase(productName, appName, secure, eventStore, eventBus, eventPublisher, commandBus, readRepos)
}
