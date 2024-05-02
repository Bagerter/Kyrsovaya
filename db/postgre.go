package db

import (
	"context"
	"crypto/sha256"
	"fmt"
	"newproject/models"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
)

var DBConn *pgxpool.Pool
var C *cache.Cache

func init() {
	C = cache.New(10*time.Minute, 10*time.Minute)
	go periodicFlushToDB()
}

func InitPg(dsn string) error {
	var err error
	DBConn, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println("lox2")
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	var version string
	err = DBConn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		return fmt.Errorf("unable to execute query: %v", err)
	}

	fmt.Printf("Connected to PostgreSQL %s\n", version)
	return nil
}
func GetUserIDByUsername(username string) (int, error) {
	var userID int
	query := `SELECT userid FROM users WHERE username = $1`
	err := DBConn.QueryRow(context.Background(), query, username).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить ID пользователя: %v", err)
	}
	return userID, nil
}
func Signup(username, password, email string) error {
	userExists, err := UserExists(username)
	if err != nil {
		return err
	}
	if userExists {
		return fmt.Errorf("username already exists")
	}

	emailExists, err := EmailExists(email)
	if err != nil {
		return err
	}
	if emailExists {
		return fmt.Errorf("email already exists")
	}

	_, err = DBConn.Exec(context.Background(), `INSERT INTO users (username, mail, password) VALUES ($1, $2, $3)`, username, email, hash(password, "keygen"))
	if err != nil {
		return err
	}
	return nil
}

func GetUserCount() (int, error) {
	var count int
	err := DBConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func Login(username, password string) (bool, error) {
	var dbPassword string
	err := DBConn.QueryRow(context.Background(), `SELECT password FROM users WHERE username=$1`, username).Scan(&dbPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, fmt.Errorf("user does not exist")
		}
		return false, err
	}

	// Сравнение хешей
	if hash(password, "keygen") == dbPassword {
		return true, nil
	}
	return false, fmt.Errorf("incorrect password")
}

func hash(value, salt string) string {
	var s = append([]byte(value), []byte(salt)...)
	hash := sha256.Sum256(s)
	return fmt.Sprintf("%x", hash) // возвращает шестнадцатеричное представление хеша
}
func UserExists(username string) (bool, error) {
	var exists bool
	err := DBConn.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
func EmailExists(email string) (bool, error) {
	var exists bool
	err := DBConn.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM users WHERE mail=$1)`, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// В файле db/news.go

// AddNews добавляет новость в базу данных.
func AddNews(title, content string) error {
	commandTag, err := DBConn.Exec(context.Background(), `INSERT INTO NewsItems (title, content) VALUES ($1, $2)`, title, content)
	if err != nil || commandTag.RowsAffected() == 0 {
		return fmt.Errorf("failed to add news item: %v", err)
	}
	return nil
}

// В файле db/news.go

// DeleteNews удаляет новость из базы данных по идентификатору.
func DeleteNews(id string) error {
	commandTag, err := DBConn.Exec(context.Background(), `DELETE FROM NewsItems WHERE id = $1`, id)
	if err != nil || commandTag.RowsAffected() == 0 {
		return fmt.Errorf("failed to delete news item: %v", err)
	}
	return nil
}

// В файле db/news.go

// GetAllNews извлекает все новости из базы данных.
func GetAllNews() ([]models.NewsItem, error) {
	var news []models.NewsItem
	rows, err := DBConn.Query(context.Background(), "SELECT id, title, content, date FROM NewsItems ORDER BY date DESC")
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ni models.NewsItem
		if err := rows.Scan(&ni.ID, &ni.Title, &ni.Content, &ni.Date); err != nil {
			return nil, fmt.Errorf("rows scan failed: %v", err)
		}
		news = append(news, ni)
	}
	return news, nil
}
func AddDomain(domain *models.Domain) error {
	sql := `INSERT INTO domains (name, cloudflare, ratelimit, user_id, Ip) VALUES ($1, $2, $3, $4, $5)`
	_, err := DBConn.Exec(context.Background(), sql, domain.Name, domain.Cloudflare, domain.RateLimit, domain.UserID, domain.Ip)
	if err != nil {
		return err
	}
	return nil
}
func DeleteDomain(id int) error {
	sql := `DELETE FROM domains WHERE id = $1`
	_, err := DBConn.Exec(context.Background(), sql, id)
	if err != nil {
		return err
	}
	return nil
}
func UpdateDomain(domain *models.Domain) error {
	sql := `UPDATE domains SET name = $1, cloudflare = $2, ratelimit = $3, user_id = $4, Ip = $5 WHERE id = $6`
	_, err := DBConn.Exec(context.Background(), sql, domain.Name, domain.Cloudflare, domain.RateLimit, domain.UserID, domain.Ip, domain.ID)
	if err != nil {
		return err
	}
	return nil
}
func GetAllDomains(userID int) ([]models.Domain, error) {
	var domains []models.Domain
	sql := `SELECT id, name, ip, cloudflare, ratelimit FROM domains WHERE user_id = $1`
	rows, err := DBConn.Query(context.Background(), sql, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d models.Domain
		if err := rows.Scan(&d.ID, &d.Name, &d.Ip, &d.Cloudflare, &d.RateLimit); err != nil {
			fmt.Println(err)
			return nil, err
		}
		domains = append(domains, d)
	}

	return domains, nil
}
func GetAttackCountsForLast7Days() ([]int, error) {
	// Массив для хранения количества атак за каждый из последних 7 дней
	counts := make([]int, 7)

	// SQL-запрос для получения количества атак за каждый из последних 7 дней
	query := `
    SELECT date_trunc('day', attack_date) as day, count(*)
    FROM attack_info
    WHERE attack_date >= current_date - interval '7 days'
    GROUP BY day
    ORDER BY day ASC
    `

	rows, err := DBConn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос к БД: %v", err)
	}
	defer rows.Close()

	// Инициализируем счётчик для обхода последних 7 дней
	dayCounts := make(map[string]int)
	for rows.Next() {
		var day time.Time
		var count int
		if err := rows.Scan(&day, &count); err != nil {
			return nil, fmt.Errorf("не удалось прочитать данные из БД: %v", err)
		}
		dayCounts[day.Format("02-01-2006")] = count
	}

	// Заполняем массив counts данными, используя dayCounts
	for i, day := range GetLast7DaysLabels() {
		if count, ok := dayCounts[day]; ok {
			counts[i] = count
		} else {
			counts[i] = 0 // Если за день атак не было, устанавливаем 0
		}
	}

	return counts, nil
}
func GetTotalAttackCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM attack_info"
	err := DBConn.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить общее количество атак: %v", err)
	}
	return count, nil
}
func GetLast7DaysLabels() []string {
	labels := []string{}
	now := time.Now()
	for i := 0; i > -7; i-- {
		day := now.AddDate(0, 0, i).Format("02/01")
		labels = append([]string{day}, labels...)
	}
	return labels
}
func DomainBelongsToUser(domainID, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM domains WHERE id = $1 AND user_id = $2)`
	err := DBConn.QueryRow(context.Background(), query, domainID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking domain ownership: %v", err)
	}
	return exists, nil
}

func GetDomainAttacks(domainID int) ([]models.AttackInfo, error) {
	var attacks []models.AttackInfo
	rows, err := DBConn.Query(context.Background(),
		`SELECT attack_date, max_power, requests_per_minute FROM attack_info WHERE domain_id = $1 ORDER BY attack_date DESC`,
		domainID)
	if err != nil {
		return nil, fmt.Errorf("error querying attack info: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var attack models.AttackInfo
		if err := rows.Scan(&attack.AttackDate, &attack.MaxPower, &attack.RequestsPerMinute); err != nil {
			return nil, fmt.Errorf("error scanning attack info: %v", err)
		}
		attacks = append(attacks, attack)
	}

	return attacks, nil
}
func GetAttackerIPs(domainID int) ([]models.IPData, error) {
	var ips []models.IPData
	rows, err := DBConn.Query(context.Background(), `SELECT ip, asn, country_code FROM attacker_ips WHERE domain_id = $1`, domainID)
	if err != nil {
		return nil, fmt.Errorf("error querying attacker ips: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ip, asn, countryCode string
		if err := rows.Scan(&ip, &asn, &countryCode); err != nil {
			return nil, fmt.Errorf("error scanning attacker ips: %v", err)
		}
		ips = append(ips, models.IPData{
			IP:      ip,
			ASN:     asn,
			Country: countryCode,
		})
	}
	return ips, nil
}
func GetVisitorIPs(domainID int) ([]models.IPData, error) {
	var ips []models.IPData
	rows, err := DBConn.Query(context.Background(), `SELECT ip, asn, country_code FROM visitor_ips WHERE domain_id = $1`, domainID)
	if err != nil {
		return nil, fmt.Errorf("error querying visitor ips: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ip, asn, countryCode string
		if err := rows.Scan(&ip, &asn, &countryCode); err != nil {
			return nil, fmt.Errorf("error scanning visitor ips: %v", err)
		}
		ips = append(ips, models.IPData{
			IP:      ip,
			ASN:     asn,
			Country: countryCode,
		})
	}
	return ips, nil
}
func GetIPCountsByCountry(domainID int, table string) ([]models.CountryData, error) {
	var countryData []models.CountryData

	// Формируем SQL-запрос в зависимости от таблицы (посетители или атакующие)
	query := fmt.Sprintf(`SELECT country_code, COUNT(*) FROM %s WHERE domain_id = $1 GROUP BY country_code`, table)

	rows, err := DBConn.Query(context.Background(), query, domainID)
	if err != nil {
		return nil, fmt.Errorf("error querying %s: %v", table, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cd models.CountryData
		if err := rows.Scan(&cd.CountryCode, &cd.Value); err != nil {
			return nil, fmt.Errorf("error scanning country data: %v", err)
		}
		countryData = append(countryData, cd)
	}

	return countryData, nil
}
func CacheAllDomains() error {
	var domains []models.Domain
	sql := `SELECT id, name, Ip, cloudflare, ratelimit FROM domains`
	rows, err := DBConn.Query(context.Background(), sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var d models.Domain
		if err := rows.Scan(&d.ID, &d.Name, &d.Ip, &d.Cloudflare, &d.RateLimit); err != nil {
			return err
		}
		domains = append(domains, d)
		// Добавляем каждый домен в кеш, используя его имя в качестве ключа
		C.Set(fmt.Sprintf("domain_%s", d.Name), d, cache.DefaultExpiration)
	}

	// Дополнительно кешируем весь список доменов для быстрого доступа
	C.Set("all_domains", domains, cache.DefaultExpiration)

	return nil
}

func CacheDomainByName(domainName string) (*models.Domain, error) {
	var domain models.Domain
	sql := `SELECT id, name, Ip, cloudflare, ratelimit FROM domains WHERE name = $1`
	err := DBConn.QueryRow(context.Background(), sql, domainName).Scan(&domain.ID, &domain.Name, &domain.Ip, &domain.Cloudflare, &domain.RateLimit)
	if err != nil {
		return nil, err // Возвращаем ошибку, если домен не найден или произошла другая ошибка
	}

	// Добавляем домен в кеш
	C.Set(fmt.Sprintf("domain_%s", domain.Name), domain, cache.DefaultExpiration)
	return &domain, nil
}

// GetDomainFromCacheByName - функция для поиска домена в кеше по имени
func GetDomainFromCacheByName(domainName string) (*models.Domain, bool) {
	domainInterface, found := C.Get(fmt.Sprintf("domain_%s", domainName))
	if found {
		if domain, ok := domainInterface.(models.Domain); ok {
			return &domain, true
		}
	}

	// Если домен не найден в кэше, пытаемся кэшировать и возвращаем его
	domain, err := CacheDomainByName(domainName)
	if err != nil {
		fmt.Println("Ошибка при кэшировании домена:", err)
		return nil, false
	}

	return domain, true
}
func AddToCache(domainID int, ip, asn, countryCode string, attacker bool) {
	key := fmt.Sprintf("stats_%d_%s", domainID, ip)
	value := models.CachedIPInfo{DomainID: domainID, IP: ip, ASN: asn, CountryCode: countryCode, IsAttacker: attacker}
	C.Set(key, value, cache.DefaultExpiration)
}

// Функция для периодической выгрузки данных из кэша в БД
func periodicFlushToDB() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	fmt.Println("Запуск периодической выгрузки данных...")
	for range ticker.C {
		fmt.Println("Попытка выгрузки данных...")
		if err := flushToDB(); err != nil {
			fmt.Printf("Ошибка при выгрузке данных в БД: %v\n", err)
		} else {
			fmt.Println("Данные успешно выгружены.")
		}
	}
}

func flushToDB() error {
	tx, err := DBConn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	commitRequired := false

	for k, v := range C.Items() {
		// Проверяем, начинается ли ключ с "stats_"
		if strings.HasPrefix(k, "stats_") {
			fmt.Printf("Ключ: %s, Значение: %+vn", k, v.Object)
			ipInfo, ok := v.Object.(models.CachedIPInfo)
			if !ok {
				fmt.Printf("Неверный тип для ключа %s, ожидался models.CachedIPInfo\n", k)
				continue
			}

			// Определяем, к какой таблице относится запись
			tableName := "visitor_ips"
			if ipInfo.IsAttacker {
				tableName = "attacker_ips"
			}

			// Вставляем данные в соответствующую таблицу
			query := fmt.Sprintf(`INSERT INTO %s (domain_id, ip, asn, country_code, created_at) VALUES ($1, $2, $3, $4, $5)`, tableName)
			fmt.Printf("Выполнение запроса к таблице %s...\n", tableName)
			result, err := tx.Exec(context.Background(), query, ipInfo.DomainID, ipInfo.IP, ipInfo.ASN, ipInfo.CountryCode, time.Now())
			if err != nil {
				tx.Rollback(context.Background())
				return fmt.Errorf("failed to insert ip info into %s: %v", tableName, err)
			}
			fmt.Printf("Запрос выполнен успешно, затронуто строк: %d\n", result.RowsAffected())

			// Удаляем обработанный элемент из кэша
			C.Delete(k)
			commitRequired = true
		}
	}

	// Фиксируем транзакцию, если были изменения
	if commitRequired {
		if err := tx.Commit(context.Background()); err != nil {
			return fmt.Errorf("failed to commit transaction: %v", err)
		}
	} else {
		tx.Rollback(context.Background())
	}

	return nil
}
