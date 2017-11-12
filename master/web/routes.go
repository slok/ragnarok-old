package web

// DefaultAPIV1Routes are the default routes for the rest v1 api paths.
var DefaultAPIV1Routes = APIV1Routes{
	Debug:            "/api/v1/debug",
	CreateExperiment: "/api/v1/experiment",
}

// APIV1Routes are the rest api v1 routes.
type APIV1Routes struct {
	Debug            string
	CreateExperiment string
}

// DefaultHTTPRoutes are the default routes for the HTTP paths.
var DefaultHTTPRoutes = HTTPRoutes{
	APIV1: DefaultAPIV1Routes,
}

// HTTPRoutes are the routes that will serve the handlers.
type HTTPRoutes struct {
	APIV1 APIV1Routes
}
