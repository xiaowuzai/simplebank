package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/xiaowuzai/simplebank/util"
)

// valiadCurrency 参数验证器，验证是否支持货币
var valiadCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}

	return false
}
