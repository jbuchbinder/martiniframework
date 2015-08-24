package martiniframework

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"net/http"
)

// ContentMiddleware is a Martini handler which specifies the proper
// serialization (XML/JSON) depending on the "Content-Type" header
// presented.
func ContentMiddleware(c martini.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case "application/xml":
		c.MapTo(encoder.XmlEncoder{}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	default:
		c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}
