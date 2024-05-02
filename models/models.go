package models

import (
	"time"
)

type DashboardData struct {
	IsAdmin      bool
	Username     string
	TotalUsers   int
	TotalAttacks int
	OnlineBots   int
	NewsItems    []NewsItem
	GraphData    GraphData
}

type NewsItem struct {
	ID      int
	Title   string
	Content string
	Date    time.Time
}

type GraphData struct {
	Labels []string
	Data   []int
}
type Domain struct {
	ID         int    `json:"id"`
	Ip         string `json:"ip"`
	Name       string `json:"name"`
	Cloudflare bool   `json:"cloudflare"`
	RateLimit  int    `json:"ratelimit"`
	UserID     int    `json:"user_id"`
}

type AttackInfo struct {
	AttackDate        time.Time
	MaxPower          int
	RequestsPerMinute int
	AttackerIP        []string
}
type StatisticData struct {
	DomainName              string
	TotalAttacks            int
	ReflectedAttacks        int
	MaxAttackPower          int
	AvgRequestsPerIP        int
	VisitCountriesData      []CountryData // Данные для карты посещений
	AttackCountriesData     []CountryData // Данные для карты атак
	VisitorIPs              []IPData      // IP адреса посетителей
	AttackerIPs             []IPData      // IP адреса атакующих
	VisitCountriesDataJSON  string
	AttackCountriesDataJSON string
}

type CountryData struct {
	CountryCode string // Код страны, например "RU", "US"
	Value       int    // Значение, например количество посещений или атак
}

type IPData struct {
	IP      string
	ASN     string
	Country string
}
type CachedIPInfo struct {
	DomainID    int       `json:"domain_id"`
	IP          string    `json:"ip"`
	ASN         string    `json:"asn"`
	CountryCode string    `json:"country_code"`
	IsAttacker  bool      `json:"is_attacker"` // true для attacker_ips, false для visitor_ips
	CreatedAt   time.Time `json:"created_at"`
}
