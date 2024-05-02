package config

import (
	"fmt"
	"newproject/db"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	mutex sync.Mutex
)

const (
	ModeNone = iota
	ModeCookieBased
	ModeJSBased
	ModeCaptchaBased
)

func IncrementCounter(key string) {
	mutex.Lock()
	defer mutex.Unlock()

	val, found := db.C.Get(key)
	if found {
		db.C.Set(key, val.(int)+1, 10*time.Second)
	} else {
		db.C.Set(key, 1, 10*time.Second)
	}
}

func CheckAndUpdateMode(ThresholdSuccess int, domainName string) int {
	ThresholdFail := ThresholdSuccess / 2
	mutex.Lock()
	defer mutex.Unlock()

	successKey := fmt.Sprintf("%s_success", domainName)
	failKey := fmt.Sprintf("%s_fail", domainName)

	// Получаем значения по сформированным ключам
	success, _ := db.C.Get(successKey)
	fail, _ := db.C.Get(failKey)

	successCount := 0
	failCount := 0
	if success != nil {
		successCount = success.(int)
	}
	if fail != nil {
		failCount = fail.(int)
	}

	currentMode, found := db.C.Get("mode")
	if !found {
		currentMode = ModeNone
	}

	// Повышение режима
	if successCount >= ThresholdSuccess && currentMode.(int) < ModeCaptchaBased {
		newMode := currentMode.(int) + 1
		// Гарантируем, что новый режим не выходит за пределы допустимого диапазона
		if newMode > ModeCaptchaBased {
			newMode = ModeCaptchaBased
		}
		db.C.Set("mode", newMode, cache.NoExpiration)
		db.C.Delete("success")
		db.C.Delete("fail")
		fmt.Println("Режим повышен до", newMode)
		return newMode
	}

	if failCount >= ThresholdFail && currentMode.(int) > ModeNone {
		newMode := currentMode.(int) - 1
		if newMode < ModeNone {
			newMode = ModeNone
		}
		db.C.Set("mode", newMode, cache.NoExpiration)
		db.C.Delete("success")
		db.C.Delete("fail")
		fmt.Println("Режим понижен до", newMode)
		return newMode
	}
	return currentMode.(int)
}
