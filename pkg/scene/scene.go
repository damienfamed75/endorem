package scene

type Scene interface {
	Preload()
	Update(float32)
	Draw()
	Unload()
	String() string
}
