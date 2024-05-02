package firewall

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var (
	//Known fingerprints along with what browser/tool/bot/etc they belong to

	//READONLY
	KnownFingerprints = map[string]string{
		//Windows
		"0x1301,0x1302,0x1303,0xc02b,0xc02f,0xc02c,0xc030,0xcca9,0xcca8,0xc013,0xc014,0x9c,0x9d,0x2f,0x35,0x583235353139,0x437572766550323536,0x437572766550333834,0x0,":                                                                        "Chromium",
		"0x1303,0x1302,0xc02b,0xc02f,0xcca9,0xcca8,0xc02c,0xc030,0xc00a,0xc009,0xc013,0xc014,0x9c,0x9d,0x2f,0x35,0x437572766550323536,0x437572766550333834,0x437572766550353231,0x437572766549442832353629,0x437572766549442832353729,0x0,":     "Firefox",
		"0x1303,0x1302,0xc02b,0xc02f,0xcca9,0xcca8,0xc02c,0xc030,0xc013,0xc014,0x9c,0x9d,0x2f,0x35,0x437572766550323536,0x437572766550333834,0x437572766550353231,0x437572766549442832353629,0x437572766549442832353729,0x0,":                   "Firefox-Dev",
		"0x1301,0x1302,0x1302,0x1303,0xc02b,0xc02f,0xc02c,0xc030,0xcca9,0xcca8,0xc013,0xc014,0x9c,0x9d,0x2f,0x35,0x583235353139,0x437572766550323536,0x437572766550333834,0x0,":                                                                 "Edge",
		"0x1303,0x1302,0xc02b,0xc02f,0xcca9,0xcca8,0xc02c,0xc030,0xc00a,0xc009,0xc013,0xc014,0x9c,0x9d,0x2f,0x35,0xa,0x437572766550323536,0x437572766550333834,0x437572766550353231,0x437572766549442832353629,0x437572766549442832353729,0x0,": "Tor",

		//IPhone
		"0x1301,0x1302,0x1303,0xc02c,0xc02b,0xcca9,0xc030,0xc02f,0xcca8,0xc00a,0xc009,0xc014,0xc013,0x9d,0x9c,0x35,0x2f,0xc008,0xc012,0xa,0x583235353139,0x437572766550323536,0x437572766550333834,0x437572766550353231,0x0,": "Safari",

		//Android
		"0xc02c,0xc02f,0xc02b,0x9f,0x9e,0xc032,0xc02e,0xc031,0xc02d,0xa5,0xa1,0xa4,0xa0,0xc028,0xc024,0xc014,0xc00a,0xc02a,0xc026,0xc00f,0xc005,0xc027,0xc023,0xc013,0xc009,0xc029,0xc025,0xc00e,0xc004,0x6b,0x69,0x68,0x39,0x37,0x36,0x67,0x3f,0x3e,0x33,0x31,0x30,0x9d,0x9c,0x3d,0x35,0x3c,0x2f,0xff,0x437572766550353231,0x437572766550333834,0x4375727665494428323229,0x0,": "Dalvik",
	}

	//READONLY
	BotFingerprints = map[string]string{
		//Bots
		"0xc030,0x9f,0xcca9,0xcca8,0xccaa,0xc02b,0xc02f,0x9e,0xc024,0xc028,0x6b,0xc023,0xc027,0x67,0xc00a,0xc014,0x39,0xc009,0xc013,0x33,0x9d,0x9c,0x3d,0x3c,0x35,0x2f,0xff,0x437572766550323536,0x437572766550353231,0x437572766550333834,0x0,":                                                                                                                                                                                                                                       "Checkhost",
		"0x1303,0x1301,0x1305,0x1304,0xc030,0xc02c,0xc028,0xc024,0xc014,0xc00a,0xa3,0x9f,0x6b,0x6a,0x39,0x38,0x88,0x87,0x9d,0x3d,0x35,0x84,0xc02f,0xc02b,0xc027,0xc023,0xc013,0xc009,0xa2,0x9e,0x67,0x40,0x33,0x32,0x9a,0x99,0x45,0x44,0x9c,0x3c,0x2f,0x96,0x41,0xff,0x437572766550323536,0x4375727665494428333029,0x437572766550353231,0x437572766550333834,0x0,":                                                                                                                     "Host-Tracker (http)",
		"(0xcca9,0xcca8,0xc02b,0xc02f,0xc02c,0xc030,0xc013,0xc014,0x9c,0x9d,0x2f,0x35,0xa,0x583235353139,0x437572766550323536,0x437572766550333834,0x0,":                                                                                                                                                                                                                                                                                                                               "Host-Tracker (page-speed)",
		"0x1303,0x1301,0xc02f,0xc02b,0xc030,0xc02c,0x9e,0xc027,0x67,0xc028,0x6b,0xa3,0x9f,0xcca9,0xcca8,0xccaa,0xc0af,0xc0ad,0xc0a3,0xc09f,0xc05d,0xc061,0xc057,0xc053,0xa2,0xc0ae,0xc0ac,0xc0a2,0xc09e,0xc05c,0xc060,0xc056,0xc052,0xc024,0x6a,0xc023,0x40,0xc00a,0xc014,0x39,0x38,0xc009,0xc013,0x33,0x32,0x9d,0xc0a1,0xc09d,0xc051,0x9c,0xc0a0,0xc09c,0xc050,0x3d,0x3c,0x35,0x2f,0xff,0x437572766550323536,0x4375727665494428333029,0x437572766550353231,0x437572766550333834,0x0,": "Postman",

		//Tools
		"0x1303,0x1301,0xc02c,0xc030,0x9f,0xcca9,0xcca8,0xccaa,0xc02b,0xc02f,0x9e,0xc024,0xc028,0x6b,0xc023,0xc027,0x67,0xc00a,0xc014,0x39,0xc009,0xc013,0x33,0x9d,0x9c,0x3d,0x3c,0x35,0x2f,0xff,0x437572766550323536,0x4375727665494428333029,0x437572766550353231,0x437572766550333834,0x0,":                                                                                                                          "Curl",
		"0xc02c,0xc028,0xc024,0xc014,0xc00a,0xa5,0xa1,0x9f,0x6b,0x69,0x68,0x39,0x37,0x36,0x88,0x86,0x85,0xc032,0xc02e,0xc02a,0xc026,0xc00f,0xc005,0x9d,0x3d,0x35,0x84,0xc02f,0xc02b,0xc027,0xc023,0xc013,0xc009,0xa4,0xa0,0x9e,0x67,0x3f,0x3e,0x33,0x31,0x30,0x45,0x43,0x42,0xc031,0xc02d,0xc029,0xc025,0xc00e,0xc004,0x9c,0x3c,0x2f,0x41,0xff,0x437572766550353231,0x437572766550333834,0x4375727665494428323229,0x0,": "Aio-http",

		//Crawler
		"0x1303,0x1301,0xc02c,0xc030,0x9f,0xcca9,0xcca8,0xccaa,0xc02b,0xc02f,0x9e,0xc024,0xc028,0x6b,0xc023,0xc027,0x67,0xc00a,0xc014,0x39,0xc009,0xc013,0x33,0x9d,0x9c,0x3d,0x3c,0x35,0x2f,0xff,0x437572766550323536,0x4375727665494428333029,0x437572766550353231,0x437572766550333834,0x437572766549442832353629,0x437572766549442832353729,0x437572766549442832353829,0x437572766549442832353929,0x437572766549442832363029,0x0,":                                                                                                                                                                                                                                                                                                                                                                              "DataForSeo",
		"0x1303,0x1301,0xc02c,0xc030,0xc02b,0xc02f,0xcca9,0xcca8,0x9f,0x9e,0xccaa,0xc0af,0xc0ad,0xc0ae,0xc0ac,0xc024,0xc028,0xc023,0xc027,0xc00a,0xc014,0xc009,0xc013,0xc0a3,0xc09f,0xc0a2,0xc09e,0x6b,0x67,0x39,0x33,0x9d,0x9c,0xc0a1,0xc09d,0xc0a0,0xc09c,0x3d,0x3c,0x35,0x2f,0xff,0x437572766550323536,0x4375727665494428333029,0x437572766550353231,0x437572766550333834,0x0,":                                                                                                                                                                                                                                                                                                                                                                                                                                 "Python-Requests",
		"0xc087,0xcca9,0xc0ad,0xc00a,0xc02b,0xc086,0xc0ac,0xc009,0xc008,0xc030,0xc08b,0xcca8,0xc014,0xc02f,0xc08a,0xc013,0xc012,0x9d,0xc07b,0xc09d,0x35,0x84,0x9c,0xc07a,0xc09c,0x2f,0x41,0xa,0x9f,0xc07d,0xccaa,0xc09f,0x39,0x88,0x9e,0xc07c,0xc09e,0x33,0x45,0x16,0x437572766550333834,0x437572766550353231,0x4375727665494428323129,0x4375727665494428313929,0x0,":                                                                                                                                                                                                                                                                                                                                                                                                                                              "Unsolicited Cralwer",
		"0xc02c,0xc028,0xc024,0xc014,0xc00a,0xa5,0xa3,0xa1,0x9f,0x6b,0x6a,0x69,0x68,0x39,0x38,0x37,0x36,0x88,0x87,0x86,0x85,0xc032,0xc02e,0xc02a,0xc026,0xc00f,0xc005,0x9d,0x3d,0x35,0x84,0xc02f,0xc02b,0xc027,0xc023,0xc013,0xc009,0xa4,0xa2,0xa0,0x9e,0x67,0x40,0x3f,0x3e,0x33,0x32,0x31,0x30,0x9a,0x99,0x98,0x97,0x45,0x44,0x43,0x42,0xc031,0xc02d,0xc029,0xc025,0xc00e,0xc004,0x9c,0x3c,0x2f,0x96,0x41,0x7,0xc011,0xc007,0xc00c,0xc002,0x5,0x4,0xc012,0xc008,0x16,0x13,0x10,0xd,0xc00d,0xc003,0xa,0xff,0x437572766550353231,0x4375727665494428323829,0x4375727665494428323729,0x437572766550333834,0x4375727665494428323629,0x4375727665494428323229,0x4375727665494428313429,0x4375727665494428313329,0x4375727665494428313129,0x4375727665494428313229,0x43757276654944283929,0x4375727665494428313029,0x0,": "Unsolicited Crawler",
	}

	//READONLY
	ForbiddenFingerprints = map[string]string{
		"0x1303,0x1302,0xc02f,0xc02b,0xc030,0xc02c,0x9e,0xc027,0x67,0xc028,0x6b,0x9f,0xcca9,0xcca8,0xccaa,0xc0af,0xc0ad,0xc0a3,0xc09f,0xc05d,0xc061,0xc053,0xc0ae,0xc0ac,0xc0a2,0xc09e,0xc05c,0xc060,0xc052,0xc024,0xc023,0xc00a,0xc014,0x39,0xc009,0xc013,0x33,0x9d,0xc0a1,0xc09d,0xc051,0x9c,0xc0a0,0xc09c,0xc050,0x3d,0x3c,0x35,0x2f,0xff,0x437572766550323536,0x4375727665494428333029,0x437572766550353231,0x437572766550333834,0x437572766549442832353629,0x437572766549442832353729,0x437572766549442832353829,0x437572766549442832353929,0x437572766549442832363029,0x0,":                                             "Http-Flood (1)",
		"0x1303,0x1301,0xc02f,0xc02b,0xc030,0xc02c,0x9e,0xc027,0x67,0xc028,0x6b,0xa3,0x9f,0xcca9,0xcca8,0xccaa,0xc0af,0xc0ad,0xc0a3,0xc09f,0xc05d,0xc061,0xc057,0xc053,0xa2,0xc0ae,0xc0ac,0xc0a2,0xc09e,0xc05c,0xc060,0xc056,0xc052,0xc024,0x6a,0xc023,0x40,0xc00a,0xc014,0x39,0x38,0xc009,0xc013,0x33,0x32,0x9d,0xc0a1,0xc09d,0xc051,0x9c,0xc0a0,0xc09c,0xc050,0x3d,0x3c,0x35,0x2f,0xff,0x437572766550323536,0x4375727665494428333029,0x437572766550353231,0x437572766550333834,0x437572766549442832353629,0x437572766549442832353729,0x437572766549442832353829,0x437572766549442832353929,0x437572766549442832363029,0x0,": "Headless Browser",
		"0x1301,0x1302,0x1303,0xc02b,0xc02f,0xc02c,0xc030,0xcca9,0xcca8,0xc013,0xc014,0x9c,0x9d,0x2f,0x35,0xa,0x4375727665494428313636393629,0x583235353139,0x437572766550323536,0x437572766550333834,0x0,": "Headless Browser",
	}
)

func GetFP(c *fiber.Ctx) string {
	ua := c.Get("User-Agent")
	return ua
}

func GenerateJA3Hash(c *fiber.Ctx) (string, error) {
	tlsInfo := c.ClientHelloInfo()
	if tlsInfo == nil {
		return "", fmt.Errorf("eror")
	}

	var fingerprint string

	for _, suite := range tlsInfo.CipherSuites[1:] {
		fingerprint += fmt.Sprintf("0x%x,", suite)
	}

	for _, curve := range tlsInfo.SupportedCurves[1:] {
		fingerprint += fmt.Sprintf("0x%x,", curve)
	}

	for _, point := range tlsInfo.SupportedPoints[:1] {
		fingerprint += fmt.Sprintf("0x%x,", point)
	}

	return fingerprint, nil
}
