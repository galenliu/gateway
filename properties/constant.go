package properties

type Access string

const (
	PermRead   = "pr"
	PermWrite  = "pw"
	PermEvents = "ev"
)

const (
	UnitPercentage = "percentage" //百分比
	UnitCelsius    = "celsius"    //摄氏度
	UnitLux        = "lux"        //亮度
	UnitSeconds    = "seconds"    //时间 秒
)
