package item

type Usable interface {
	Use(*EffectData)

	String() string
}
