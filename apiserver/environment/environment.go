package environment

import (
	"fmt"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

const FRONT_END_PORT int = 5173
const PORT int = 3333

var REDIS_HOST string = getEnv("REDIS_HOST", "localhost")
var DOMAIN string = getEnv("DOMAIN", "blackjack.gg")
var FRONT_END_URL string = fmt.Sprintf("https://%s:%d", DOMAIN, FRONT_END_PORT)
var BLACKJACK_SERVER_PATH string = fmt.Sprintf("%s/blackjack/", DOMAIN)
