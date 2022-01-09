package yeelight

func checkBrightnessValue(b int) bool {
	return b < 0 || b > 100
}
