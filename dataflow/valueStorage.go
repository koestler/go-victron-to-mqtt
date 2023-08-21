package dataflow

type State map[string]ValueMap

type ValueStorageInstance struct {
	// this represents the state of the storage instance and must only be access by the main go routine

	// state: 1. dimension: device.Name, 2. dimension: register.Name
	state         State
	subscriptions map[*Subscription]struct{}

	// communication channels to/from the main go routine
	inputChannel            chan Value
	subscriptionChannel     chan *Subscription
	readStateRequestChannel chan *readStateRequest

	shutdown chan struct{}
}

type SkipRegisterNameStruct struct {
	Device   string
	Register string
}

type SkipRegisterCategoryStruct struct {
	Device   string
	Category string
}

type OnlyOnceKey struct {
	deviceName   string
	registerName string
}

type Filter struct {
	IncludeDevices         map[string]bool
	SkipRegisterNames      map[SkipRegisterNameStruct]bool
	SkipRegisterCategories map[SkipRegisterCategoryStruct]bool
	SkipNull               bool
	OnlyOnce               bool
}

type readStateRequest struct {
	filter   Filter
	response chan State
}

func (instance *ValueStorageInstance) mainStorageRoutine() {
	for {
		select {
		case <-instance.shutdown:
			return
		case newValue := <-instance.inputChannel:
			instance.handleNewValue(newValue)
		case newSubscription := <-instance.subscriptionChannel:
			instance.subscriptions[newSubscription] = struct{}{}
		case newReadStateRequest := <-instance.readStateRequestChannel:
			instance.handleNewReadStateRequest(newReadStateRequest)
		}
	}
}

func (instance *ValueStorageInstance) handleNewValue(newValue Value) {
	// make sure device exists
	if _, ok := instance.state[newValue.DeviceName()]; !ok {
		instance.state[newValue.DeviceName()] = make(ValueMap)
	}

	// check if the newValue is not present or has been changed
	if currentValue, ok := instance.state[newValue.DeviceName()][newValue.Register().Name()]; !ok || !currentValue.Equals(newValue) {
		// copy the input value to all subscribed output channels
		for subscription := range instance.subscriptions {
			// check if Subscription was shut down
			select {
			case <-subscription.shutdownChannel:
				delete(instance.subscriptions, subscription)
			default:
				// Subscription was not shut down -> forward new value
				if !subscription.filter.OnlyOnce {
					subscription.forward(newValue)
				} else {
					k := OnlyOnceKey{newValue.DeviceName(), newValue.Register().Name()}
					if _, ok := subscription.sentOnce[k]; !ok {
						subscription.forward(newValue)
						subscription.sentOnce[k] = true
					}
				}
			}
		}

		if _, ok := newValue.(NullRegisterValue); ok {
			delete(instance.state[newValue.DeviceName()], newValue.Register().Name())
		} else {
			// and save the new state
			instance.state[newValue.DeviceName()][newValue.Register().Name()] = newValue
		}
	}
}

func (instance *ValueStorageInstance) handleNewReadStateRequest(newReadStateRequest *readStateRequest) {
	filter := &newReadStateRequest.filter

	response := make(State)

	for deviceName, deviceState := range instance.state {
		if !filterByDevice(filter, deviceName) {
			continue
		}

		response[deviceName] = make(ValueMap)

		for registerName, value := range deviceState {
			if !filterByRegister(filter, deviceName, value.Register()) {
				continue
			}

			response[deviceName][registerName] = value
		}
	}

	newReadStateRequest.response <- response
}

func ValueStorageCreate() (valueStorageInstance *ValueStorageInstance) {
	valueStorageInstance = &ValueStorageInstance{
		state:                   make(State),
		subscriptions:           make(map[*Subscription]struct{}),
		inputChannel:            make(chan Value, 1024),
		subscriptionChannel:     make(chan *Subscription),
		readStateRequestChannel: make(chan *readStateRequest, 16),
		shutdown:                make(chan struct{}),
	}

	// start main go routine
	go valueStorageInstance.mainStorageRoutine()

	return
}

func (instance *ValueStorageInstance) Shutdown() {
	close(instance.shutdown)
}

func (instance *ValueStorageInstance) GetState(filter Filter) State {
	response := make(chan State)

	request := readStateRequest{
		filter:   filter,
		response: response,
	}

	instance.readStateRequestChannel <- &request

	return <-request.response
}

func (instance *ValueStorageInstance) GetSlice(filter Filter) (result []Value) {
	state := instance.GetState(filter)

	// create result slice of correct capacity
	capacity := 0
	for _, deviceState := range state {
		capacity += len(deviceState)
	}
	result = make([]Value, 0, capacity)

	for _, deviceState := range state {
		for _, value := range deviceState {
			result = append(result, value)
		}
	}
	return
}

func (instance *ValueStorageInstance) Fill(value Value) {
	instance.inputChannel <- value
}

func (instance *ValueStorageInstance) Subscribe(filter Filter) Subscription {
	s := Subscription{
		shutdownChannel: make(chan struct{}),
		outputChannel:   make(chan Value, 128),
		filter:          filter,
		sentOnce:        make(map[OnlyOnceKey]bool),
	}

	instance.subscriptionChannel <- &s

	return s
}
