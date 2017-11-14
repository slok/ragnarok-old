package web

// APIV1Routes are the rest api v1 routes.
type APIV1Routes struct {
	Debug           string
	WriteExperiment string
}

// DefaultAPIV1Routes are the default routes for the rest v1 api paths.
var DefaultAPIV1Routes = APIV1Routes{
	Debug:           "/api/v1/debug",
	WriteExperiment: "/api/v1/experiment",
}

// DefaultHTTPRoutes are the default routes for the HTTP paths.
var DefaultHTTPRoutes = HTTPRoutes{
	APIV1: DefaultAPIV1Routes,
}

// HTTPRoutes are the routes that will serve the handlers.
type HTTPRoutes struct {
	APIV1 APIV1Routes
}
