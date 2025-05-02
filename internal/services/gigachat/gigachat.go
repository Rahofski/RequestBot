package gigachat

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetAccessToken() string{

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка получения рабочей директории:", err)
		return ""
	}

	// Загружаем публичный сертификат Минцифры
	rootCA, err := os.ReadFile(cwd + "\\russian_trusted_root_ca.pem")
	if err != nil {
		fmt.Println("Ошибка чтения сертификата:", err)
		return ""
	}

	// Создаём кастомный пул доверенных сертификатов
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(rootCA) {
		fmt.Println("Ошибка добавления сертификата в пул")
		return ""
	}

	// Создаём HTTP-клиент с кастомными сертификатами
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	// Отправляем запрос на получение токена
	data := url.Values{}
	data.Set("scope", "GIGACHAT_API_PERS")

	// Формируем запрос
	url := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return ""
	}

	// Устанавливаем заголовки
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("RqUID", "4ce218e4-bb77-454a-ba7b-ebc2505c1ca3")
	//gotta move into .env file
	req.Header.Add("Authorization", "Basic NTljZWQwNGItNjUyYS00ZTk5LTk4NjgtZDkxMDFjZTEyNDI5OjA0MzA1ODQ2LWUwYjUtNDM1My04MjdmLWE2OTNjM2EzMWM0NA==")

	// Добавляем данные в тело запроса
	req.Body = io.NopCloser(strings.NewReader(data.Encode()))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка HTTP-запроса:", err)
		return ""
	}
	defer res.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return ""
	}

	return ExtractAccessToken(string(body))
}


func ExtractAccessToken(response string) string {
	// Структура для хранения данных из JSON
	type TokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	// Убираем возможные пробелы в начале/конце строки
	response = strings.TrimSpace(response)

	// Создаём переменную для хранения результата
	var tokenData TokenResponse

	// Декодируем JSON
	err := json.Unmarshal([]byte(response), &tokenData)
	if err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		return ""
	}


	return tokenData.AccessToken
}

func DefineField(message string ,accessToken string) (int, error) { 
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка получения рабочей директории:", err)
		return -1, err 
	}

	// Загружаем публичный сертификат Минцифры
	rootCA, err := os.ReadFile(cwd + "\\russian_trusted_root_ca.pem")
	if err != nil {
		fmt.Println("Ошибка чтения сертификата:", err)
		return -1, err 
	}

	// Создаём кастомный пул доверенных сертификатов
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(rootCA) {
		fmt.Println("Ошибка добавления сертификата в пул")
		return -1, err  
	}

	// Создаём HTTP-клиент с кастомными сертификатами
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	method := http.MethodPost

	url := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
  
	payload := fmt.Sprintf(`{
		"model": "GigaChat",
		"messages": [
		  {
			"role": "system",
			"content": "Ты - професстональный определитель сферы заявки. Твоя задача - определить, к какой категории относится заявка. Доступные категории: Водоснабжение , Электричество , Инфраструктура , Компютеры, Плотничество. В качестве ответа верни одно слово только из этих вариантов!"
		  },
		  {
			"role": "user",
			"content": "%s"
		  }
		],
		"stream": false,
		"update_interval": 0
	}`, message)
  
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
  
	if err != nil {
	  fmt.Println(err)
	  return -1, err 
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer " + accessToken)
  
	res, err := client.Do(req)
	if err != nil {
	  fmt.Println(err)
	  return -1, err 
	}

	defer res.Body.Close()
  
	body, err := io.ReadAll(res.Body)
	if err != nil {
	  fmt.Println(err)
	  return -1, err 
	}

	fmt.Println("Ответ сервера:", string(body))

	answer, err := ParseDefineField(body)
	if err != nil {
		fmt.Println(err)
		return -1, err 
	}
	return answer, nil
}


func CheckIfValid(message string ,accessToken string) (bool, error) { 

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка получения рабочей директории:", err)
		return false, err
	}

	// Загружаем публичный сертификат Минцифры
	rootCA, err := os.ReadFile(cwd + "\\russian_trusted_root_ca.pem")
	if err != nil {
		fmt.Println("Ошибка чтения сертификата:", err)
		return false, err 
	}

	// Создаём кастомный пул доверенных сертификатов
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(rootCA) {
		fmt.Println("Ошибка добавления сертификата в пул")
		return false, err 
	}

	// Создаём HTTP-клиент с кастомными сертификатами
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	method := http.MethodPost

	url := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
  
	payload := fmt.Sprintf(`{
		"model": "GigaChat",
		"messages": [
		  {
			"role": "system",
			"content": "Ты - предварительный модератор заявок на ремонт. Твоя задача - ответить да или нет на вопрос, является ли заявка валидной, или в ней написано что-то, неотносящееся к заявке. Если заявка пуста - верни Нет"
		  },
		  {
			"role": "user",
			"content": "%s"
		  }
		],
		"stream": false,
		"update_interval": 0
	}`, message)
  
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
  
	if err != nil {
	  fmt.Println(err)
	  return false, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer " + accessToken)
  
	res, err := client.Do(req)
	if err != nil {
	  fmt.Println(err)
	  return false, err
	}

	defer res.Body.Close()
  
	body, err := io.ReadAll(res.Body)
	if err != nil {
	  fmt.Println(err)
	  return false, err
	}

	answer, err := ParseIfValidResponse(body)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return answer, nil
  }

  type GigaChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// ParseIfValidResponse парсит JSON-ответ и возвращает true, если content != "Нет", иначе false
func ParseIfValidResponse(jsonData []byte) (bool, error) {
	var response GigaChatResponse
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return false, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	if len(response.Choices) > 0 && response.Choices[0].Message.Content == "Нет" {
		return false, nil
	}
	return true, nil
}

func ParseDefineField(jsonData []byte) (int, error) {
	var response GigaChatResponse
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return -1, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	if len(response.Choices) > 0 && response.Choices[0].Message.Content == "Водоснабжение" {
		return 1, nil
	}

	if len(response.Choices) > 0 && response.Choices[0].Message.Content == "Электричество" {
		return 2, nil
	}

	if len(response.Choices) > 0 && response.Choices[0].Message.Content == "Инфраструктура" {
		return 3, nil
	}

	if len(response.Choices) > 0 && response.Choices[0].Message.Content == "Компьютеры" {
		return 5, nil
	}

	if len(response.Choices) > 0 && response.Choices[0].Message.Content == "Плотничество" {
		return 7, nil
	}

	return -1, nil
}