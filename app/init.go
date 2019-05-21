package app

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/revel/revel"
	"github.com/revel/revel/logger"
	"net/http"
	"os"
	"strconv"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	logger.LogFunctionMap["stdoutjson"]=
		func(c *logger.CompositeMultiHandler, options *logger.LogOptions) {
			// Set the json formatter to os.Stdout, replace any existing handlers for the level specified
			c.SetJson(os.Stdout, options)
		}

	logger.LogFunctionMap["stderrjson"]=
		func(c *logger.CompositeMultiHandler, options *logger.LogOptions) {
			// Set the json formatter to os.Stdout, replace any existing handlers for the level specified
			c.SetJson(os.Stderr, options)
		}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
	revel.OnAppStart(InitLogger)
	revel.OnAppStart(InitConfig)
	revel.OnAppStart(InitTemplateEngine)
	revel.OnAppStart(initPrometheus)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

func GetConfigString(key, defaultValue string) (ret string) {
	ret = defaultValue

	// try to get config
	if val, exists := revel.Config.String(key); exists && val != "" {
		return val
	}

	// try to get config default
	if val, exists := revel.Config.String(key + ".default"); exists && val != "" {
		return val
	}

	return
}

func GetConfigInt(key string, defaultValue int) (ret int) {
	ret = defaultValue

	// try to get config
	if val, exists := revel.Config.String(key); exists && val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}

	// try to get config default
	if val, exists := revel.Config.String(key + ".default"); exists && val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}

	return
}

func GetConfigBoolean(key string, defaultValue bool) (ret bool) {
	ret = defaultValue

	// try to get config
	if val, exists := revel.Config.String(key); exists && val != "" {
		if i, err := strconv.ParseBool(val); err == nil {
			return i
		}
	}

	// try to get config default
	if val, exists := revel.Config.String(key + ".default"); exists && val != "" {
		if i, err := strconv.ParseBool(val); err == nil {
			return i
		}
	}

	return
}


func InitLogger() {
}

func InitConfig() {
}

func InitTemplateEngine() {
	revel.TemplateFuncs["config"] = func(option string) string {
		return GetConfigString(option, "")
	}
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(GetConfigString("metrics.listen", ":9001"), nil)
	}()
}
