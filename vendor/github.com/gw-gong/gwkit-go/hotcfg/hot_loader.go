package hotcfg

type HotLoader interface {
	BaseConfigCapable
	LoadConfig()
}
