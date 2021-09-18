package provider

// IAccesstoken get accesstoken by appid and component_appid
type IAccesstoken interface {
	Get(appid string, componentAppid string) (string, error)
}
