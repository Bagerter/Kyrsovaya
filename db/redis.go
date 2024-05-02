package db

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func InitRd() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // адрес сервера Redis
		Password: "",               // пароль (если есть)
		DB:       0,                // номер базы данных
	})
	pong, err := Client.Ping(Ctx).Result()
	fmt.Println(pong, err)
	// Output: PONG <nil> если Redis работает и доступен
	// В противном случае, pong будет пустым и err будет содержать информацию об ошибке
	if err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
}

func Fraud(fraudscore int, ip string) int {
	val, err := Client.Get(Ctx, "fraudscore:"+ip).Result()
	if err == redis.Nil {
		err = Client.Set(Ctx, "fraudscore:"+ip, fraudscore, time.Hour*24).Err()
		if err != nil {
			// Обрабатываем ошибку
			fmt.Println("Ошибка при установке значения в Redis:", err)
		}
		return fraudscore
	} else {
		addedScore, convErr := strconv.Atoi(val)
		if convErr != nil {
			// Обрабатываем ошибку
			fmt.Println("Ошибка при преобразовании строки в число:", convErr)
			// Возвращаем исходное значение fraudscore или какое-то другое значение в случае ошибки
			return fraudscore
		}
		fraudscore += addedScore
		err = Client.Set(Ctx, "fraudscore:"+ip, fraudscore, time.Hour*24).Err()
		if err != nil {
			fmt.Println("Ошибка при перезаписи значения в Redis:", err)
		}
	}
	return fraudscore
}

func Fraudja3(ja3hash string, ip string) bool {
	val, err := Client.Get(Ctx, "ja3hash:"+ip).Result()
	if err == redis.Nil {
		err = Client.Set(Ctx, "ja3hash:"+ip, ja3hash, time.Hour*1).Err()
		if err != nil {
			fmt.Println("Ошибка при установке значения в Redis:", err)
		}
		return true
	} else {
		if val != ja3hash {
			return false
		}
	}
	return true
}

func GetFraudIP(ip string) int {
	var fraudscore int

	val, err := Client.Get(Ctx, "fraudscore:"+ip).Result()
	if err == redis.Nil {
		fraudscore = 0
	} else {
		addedScore, convErr := strconv.Atoi(val)
		if convErr != nil {
			// Обрабатываем ошибку
			fmt.Println("Ошибка при преобразовании строки в число:", convErr)
			// Возвращаем исходное значение fraudscore или какое-то другое значение в случае ошибки
			return fraudscore
		}
		fraudscore = addedScore
	}

	return fraudscore
}

func CaptchaCache(keyForCaching string, captchaData string, privatePart string) error {

	err := Client.HSet(Ctx, keyForCaching, "value1", captchaData).Err()
	if err != nil {
		return err
	}

	err = Client.HSet(Ctx, keyForCaching, "value2", privatePart).Err()
	if err != nil {
		return err
	}

	ttl := time.Hour

	err = Client.Expire(Ctx, keyForCaching, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetCaptchaCache(keyForCaching string) (captchaData string, privatePart string) {

	captchaData, err := Client.HGet(Ctx, keyForCaching, "value1").Result()
	if err != nil {
		return "error", "error"
	}

	privatePart, err = Client.HGet(Ctx, keyForCaching, "value2").Result()
	if err != nil {
		return "error", "error"
	}
	return captchaData, privatePart
}

func GetCanvas(ip, canvas string) string {
	if canvas == "empty" {
		val, err := Client.Get(Ctx, "canvas:"+ip).Result()
		if err == nil {
			return "empty"
		}
		return val
	} else {
		err := Client.Set(Ctx, "canvas:"+ip, canvas, time.Hour*1).Err()
		if err != nil {
			// Обрабатываем ошибку
			fmt.Println("Ошибка при установке значения в Redis:", err)
		}
		return "done"
	}
}
