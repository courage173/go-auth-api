package healthcheck

import routing "github.com/go-ozzo/ozzo-routing/v2"

func RegisterHandlers(r *routing.Router, version string){
	r.To("GET,HEAD", "/healthcheck", healthcheck(version))
}

func healthcheck (version string) routing.Handler {
	return  func(c *routing.Context) error {
		return c.Write("OK " + version)
	}
}