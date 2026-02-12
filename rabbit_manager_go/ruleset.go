package rabbit_manager_go

/**
THIS FILE CAN BE MOVED TO ANOTHER PACKAGE TO SPLIT THE LOGIC
*/
import (
	"asiatix/internal/crossdomain_events"
)

type PublishRule struct {
	RoutingKey string
	Mandatory  bool
	Immediate  bool
	Exchange   *Exchange
	Queue      *Queue
}

var ruleset = map[crossdomain_events.EventName]PublishRule{
	crossdomain_events.ReceiptResultEventName: {
		Exchange: DefaultExchange,
		Queue:    DefaultQueue,
	},
}
