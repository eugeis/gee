package mongo

import (
	"github.com/looplab/eventhorizon"
	commandbus "github.com/looplab/eventhorizon/commandbus/local"
	eventbus "github.com/looplab/eventhorizon/eventbus/local"
	eventstore "github.com/looplab/eventhorizon/eventstore/mongodb"
	eventpublisher "github.com/looplab/eventhorizon/publisher/local"
	repo "github.com/looplab/eventhorizon/repo/mongodb"
	"github.com/eugeis/gee/eh"
	"github.com/eugeis/gee/eh/app"
)

func NewAppMongo(productName string, appName string, secure bool, mongoUrl string) *app.AppBase {
	// Create the event store.
	eventStore := &eh.EventStoreDelegate{Factory:
	func() (ret eventhorizon.EventStore, err error) {
		return eventstore.NewEventStore("localhost", productName)
	}}

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
				var retRepo *repo.Repo
				if retRepo, err = repo.NewRepo(mongoUrl, productName, name); err == nil {
					retRepo.SetModel(factory)
					ret = retRepo
				}
				return
			}}
			repos[name] = ret
		} else {
			ret = item
		}
		return
	}
	return app.NewAppBase(productName, appName, secure, eventStore, eventBus, eventPublisher, commandBus, readRepos)
}
