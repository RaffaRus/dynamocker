package webserver

type ApiInterface interface {
	checkVersion() error
	perform()
}

type Api struct {
	resource string
	versions []uint16
	handler map[uint16]func(http.ResponseWriter, *http.Request)
}

// api ctor
func NewApi(resource string, ver ApiVersion) (Api, error) {
	if ver == null {
		return (Api{
			resource,
			
		}, null)

	}
}