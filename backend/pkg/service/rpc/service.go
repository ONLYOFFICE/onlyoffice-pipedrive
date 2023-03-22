package rpc

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/config"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/messaging"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/middleware/wrapper"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/resilience"
	"github.com/go-micro/plugins/v4/wrapper/breaker/hystrix"
	rlimiter "github.com/go-micro/plugins/v4/wrapper/ratelimiter/uber"
	"github.com/go-micro/plugins/v4/wrapper/select/roundrobin"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"go-micro.dev/v4"
	"go-micro.dev/v4/cache"
	gcache "go-micro.dev/v4/cache"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	uber "go.uber.org/ratelimit"
)

type RPCMessageHandler struct {
	Topic   string
	Queue   string
	Handler interface{}
}

type RPCEngine interface {
	BuildMessageHandlers() []RPCMessageHandler
	BuildHandlers(client client.Client, cache gcache.Cache) []interface{}
}

// NewService Initializes an http service with options.
func NewService(
	engine RPCEngine,
	registry registry.Registry,
	broker messaging.BrokerWithOptions,
	cache cache.Cache,
	tracer *oteltrace.TracerProvider,
	rpcConfig *config.ServerConfig,
	resilienceConfig *config.ResilienceConfig,
	tracerConfig *config.TracerConfig,
) micro.Service {
	var wrappers []server.HandlerWrapper = make([]server.HandlerWrapper, 0, 2)

	if err := broker.Broker.Init(); err != nil {
		log.Fatalf("could not initialize a new broker instance: %s", err.Error())
	}

	if err := broker.Broker.Connect(); err != nil {
		log.Fatalf("broker connection error: %s", err.Error())
	}

	if resilienceConfig.Resilience.RateLimiter.Limit > 0 {
		wrappers = append(wrappers, rlimiter.NewHandlerWrapper(int(resilienceConfig.Resilience.RateLimiter.Limit), uber.Per(1*time.Second)))
	}

	if tracerConfig.Tracer.Enable {
		wrappers = append(wrappers, wrapper.TracePropagationHandlerWrapper)
	}

	hystrix.ConfigureDefault(resilience.BuildHystrixCommandConfig(resilienceConfig))

	service := micro.NewService(
		micro.Name(strings.Join([]string{rpcConfig.Namespace, rpcConfig.Name}, ":")),
		micro.Version(strconv.Itoa(rpcConfig.Version)),
		micro.Context(context.Background()),
		micro.Server(server.NewServer(
			server.Name(strings.Join([]string{rpcConfig.Namespace, rpcConfig.Name}, ":")),
			server.Address(rpcConfig.Address),
		)),
		micro.Cache(cache),
		micro.Registry(registry),
		micro.Broker(broker.Broker),
		micro.Client(client.NewClient(
			client.Registry(registry),
			client.Broker(broker.Broker),
		)),
		micro.WrapClient(
			roundrobin.NewClientWrapper(),
			hystrix.NewClientWrapper(),
			opentelemetry.NewClientWrapper(opentelemetry.WithTraceProvider(otel.GetTracerProvider())),
		),
		micro.WrapSubscriber(opentelemetry.NewSubscriberWrapper(opentelemetry.WithTraceProvider(otel.GetTracerProvider()))),
		micro.WrapHandler(wrappers...),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
		micro.AfterStop(func() error {
			if tracer != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				if err := tracer.Shutdown(ctx); err != nil {
					return err
				}

				return nil
			}

			return nil
		}),
	)

	if len(engine.BuildMessageHandlers()) > 0 {
		for _, entry := range engine.BuildMessageHandlers() {
			if entry.Handler != nil && entry.Topic != "" {
				if entry.Queue != "" {
					if err := micro.RegisterSubscriber(
						entry.Topic, service.Server(), entry.Handler, server.SubscriberContext(broker.SubOptions.Context), server.SubscriberQueue(entry.Queue),
					); err != nil {
						log.Fatalf("could not register a new subscriber with topic %s. Reason: %s", entry.Topic, err.Error())
					}
				} else {
					if err := micro.RegisterSubscriber(entry.Topic, service.Server(), entry.Handler, server.SubscriberContext(broker.SubOptions.Context)); err != nil {
						log.Fatalf("could not register a new subscriber with topic %s. Reason: %s", entry.Topic, err.Error())
					}
				}
			}
		}
	}

	for _, handler := range engine.BuildHandlers(service.Client(), service.Options().Cache) {
		if err := micro.RegisterHandler(service.Server(), handler); err != nil {
			log.Fatalf("could not initialize rpc handlers: %s", err.Error())
		}
	}

	return service
}