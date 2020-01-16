package scene

type Scene interface {
	Preload(func())
	Update(float32)
	Draw()
	Unload()
	String() string
}
