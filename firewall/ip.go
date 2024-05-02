package firewall

import (
	"encoding/json"
	"io"
	"net/http"
	"newproject/db"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type IpInfo struct {
	Country struct {
		Code string `json:"alpha2_code"`
	} `json:"country"`
	AS struct {
		Num string `json:"number"`
		Org string `json:"org"`
	} `json:"as"`
}

var ipLocks sync.Map

func getIpLock(ip string) *sync.Mutex {
	mutex, _ := ipLocks.LoadOrStore(ip, &sync.Mutex{})
	return mutex.(*sync.Mutex)
}

func GetIpInfoforme(c *fiber.Ctx) (ip string, country string, asn string, org string) {
	ip = c.Get("CF-Connecting-IP")
	if ip == "" {
		ip = c.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.Get("X-Real-IP")
	}
	if ip == "" {
		ip = c.IP()
	}
	var data IpInfo
	mutex := getIpLock(ip)
	mutex.Lock()
	defer mutex.Unlock()
	val, err := db.Client.Get(db.Ctx, ip).Result()
	if err == redis.Nil {
		resp, err := http.Get("http://apimon.de/ip/" + ip)
		if err != nil {
			return ip, "UNK", "UNK", "UNK"
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return ip, "UNK", "UNK", "UNK"
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			return ip, "UNK", "UNK", "UNK"
		}
		ipInfo, err := json.Marshal(data)
		if err != nil {
			return "", "UNK", "UNK", "UNK"
		}
		err = db.Client.Set(db.Ctx, ip, ipInfo, 0).Err()
		if err != nil {
			return "", "UNK", "UNK", "UNK"
		}
	} else if err != nil {
		// Если возникла другая ошибка при получении значения из Redis, вернем "UNK" значения
		return ip, "UNK", "UNK", "UNK"
	} else {
		// Информация об IP присутствует в Redis, используем её
		err = json.Unmarshal([]byte(val), &data)
		if err != nil {
			return ip, "UNK", "UNK", "UNK"
		}
	}

	return ip, data.Country.Code, data.AS.Num, data.AS.Org
}
