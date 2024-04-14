package webserver

type ApiInterface interface {
	checkVersion() error
	perform()
}

type Api struct {
	resource string
	versions []uint16
}

// api ctor
func NewApi(resource string, ver ApiVersion) (Api, error) {
	if ver == null {
		return (Api{
			resource,
			
		}, null)

	}
}