package currency

const (
	Bitcoin = "bitcoin"
	Cosmos  = "cosmos"
)

var (
	PrefixMap = map[string]string{
		Bitcoin: "bc1q",
		Cosmos:  "cosmos1",
	}
	PathMap = map[string]string{
		Bitcoin: "m/84'/0'/0'/0/0",
		Cosmos:  "m/44'/118'/0'/0/0",
	}
)
