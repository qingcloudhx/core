package engine

import (
	"fmt"
	"github.com/qingcloudhx/core/data/property"
	"strings"

	"github.com/qingcloudhx/core/action"
	"github.com/qingcloudhx/core/app"
	"github.com/qingcloudhx/core/engine/channels"
	"github.com/qingcloudhx/core/engine/runner"
	"github.com/qingcloudhx/core/engine/secret"
	"github.com/qingcloudhx/core/support"
	"github.com/qingcloudhx/core/support/log"
	"github.com/qingcloudhx/core/support/managed"
)

// engineImpl is the type for the Default Engine Implementation
type engineImpl struct {
	config         *Config
	flogoApp       *app.App
	actionRunner   action.Runner
	serviceManager *support.ServiceManager
	logger         log.Logger
}

type Option func(*engineImpl)

// New creates a new Engine
func New(appConfig *app.Config, options ...Option) (Engine, error) {
	if appConfig == nil {
		return nil, fmt.Errorf("no App configuration provided")
	}
	if len(appConfig.Name) == 0 {
		return nil, fmt.Errorf("no App name provided")
	}
	//log.Init(appConfig.Name)
	if len(appConfig.Version) == 0 {
		return nil, fmt.Errorf("no App version provided")
	}

	engine := &engineImpl{}
	logger := log.ChildLogger(log.RootLogger(), "engine")
	engine.logger = logger
	//log.SetLogLevel(log.DebugLevel, logger)

	if engine.config == nil {
		config := &Config{}
		config.StopEngineOnError = true
		config.RunnerType = GetRunnerType()
		//config.LogLevel = DefaultLogLevel

		engine.config = config
	}

	for _, option := range options {
		option(engine)
	}

	//add engine channels - todo should these be moved to app
	channelDescriptors := appConfig.Channels
	if len(channelDescriptors) > 0 {
		for _, descriptor := range channelDescriptors {
			name, buffSize := channels.Decode(descriptor)

			logger.Debugf("Creating Engine Channel '%s'", name)

			_, err := channels.New(name, buffSize)
			if err != nil {
				return nil, err
			}
		}
	}

	if engine.actionRunner == nil {
		var actionRunner action.Runner

		runnerType := engine.config.RunnerType
		if strings.EqualFold(ValueRunnerTypePooled, runnerType) {
			actionRunner = runner.NewPooled(NewPooledRunnerConfig())
		} else if strings.EqualFold(ValueRunnerTypeDirect, runnerType) {
			actionRunner = runner.NewDirect()
		} else {
			return nil, fmt.Errorf("unknown runner type: %s", runnerType)
		}

		logger.Debugf("Using '%s' Action Runner", runnerType)

		engine.actionRunner = actionRunner
	}

	var appOptions []app.Option
	if !engine.config.StopEngineOnError {
		appOptions = append(appOptions, app.ContinueOnError)
	}

	propResolvers := GetAppPropertyValueResolvers(logger)
	enablePropertiesResolution := false
	if len(propResolvers) > 0 {
		err := property.EnablePropertyResolvers(propResolvers)
		if err != nil {
			return nil, err
		}

		enablePropertiesResolution = true
	}

	// properties post processors (properties resolver if enabled, secret properties replacer)
	var postProcessors []property.PostProcessor
	if enablePropertiesResolution {
		postProcessors = append(postProcessors, property.PropertyResolverProcessor)
	}
	postProcessors = append(postProcessors, secret.PropertyProcessor)

	option := app.FinalizeProperties(postProcessors...)
	appOptions = append(appOptions, option)

	flogoApp, err := app.New(appConfig, engine.actionRunner, appOptions...)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Creating app [ %s ] with version [ %s ]", appConfig.Name, appConfig.Version)
	engine.flogoApp = flogoApp
	engine.serviceManager = support.GetDefaultServiceManager()

	return engine, nil
}

func (e *engineImpl) App() *app.App {
	return e.flogoApp
}

//Start initializes and starts the Triggers and initializes the Actions
func (e *engineImpl) Start() error {

	logger := e.logger

	logger.Infof("Starting app [ %s ] with version [ %s ]", e.flogoApp.Name(), e.flogoApp.Version())

	logger.Info("Engine Starting...")

	logger.Info("Starting Services...")

	actionRunner := e.actionRunner.(interface{})

	if managedRunner, ok := actionRunner.(managed.Managed); ok {
		_ = managed.Start("ActionRunner Service", managedRunner)
	}

	err := e.serviceManager.Start()

	if err != nil {
		logger.Error("Error Starting Services - " + err.Error())
	} else {
		logger.Info("Started Services")
	}

	if len(managedServices) > 0 {
		for _, mService := range managedServices {
			err = mService.Start()
			if err != nil {
				logger.Error("Error Starting Services - " + err.Error())
				//TODO Should we exit here?
			}
		}
	}

	logger.Info("Starting Application...")
	err = e.flogoApp.Start()
	if err != nil {
		return err
	}
	logger.Info("Application Started")

	if channels.Count() > 0 {
		logger.Info("Starting Engine Channels...")
		_ = channels.Start()
		logger.Info("Engine Channels Started")
	}

	logger.Info("Engine Started")

	return nil
}

func (e *engineImpl) Stop() error {

	logger := e.logger

	logger.Info("Engine Stopping...")

	if channels.Count() > 0 {
		logger.Info("Stopping Engine Channels...")
		_ = channels.Stop()
		logger.Info("Engine Channels Stopped...")
	}

	logger.Info("Stopping Application...")
	_ = e.flogoApp.Stop()
	logger.Info("Application Stopped")

	//TODO temporarily add services
	logger.Info("Stopping Services...")

	actionRunner := e.actionRunner.(interface{})

	if managedRunner, ok := actionRunner.(managed.Managed); ok {
		_ = managed.Stop("ActionRunner", managedRunner)
	}

	err := e.serviceManager.Stop()

	if err != nil {
		logger.Error("Error Stopping Services - " + err.Error())
	} else {
		logger.Info("Stopped Services")
	}

	if len(managedServices) > 0 {
		for _, mService := range managedServices {
			err = mService.Stop()
			if err != nil {
				logger.Error("Error Stopping Services - " + err.Error())
			}
		}
	}

	logger.Info("Engine Stopped")
	log.Sync()

	return nil
}
