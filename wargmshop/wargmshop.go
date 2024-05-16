package wargmshop

import (
    "encoding/json"
    "fmt"
    "net/http"
    "draw/logger"
	"draw/givepak"
	"draw/restjs"	
)

// Purchase представляет информацию о покупке.
type Purchase struct {
    ID           int    `json:"id"`
    ServerID     int    `json:"server_id"`
    UserID       int    `json:"user_id"`
    OfferID      int    `json:"offer_id"`
    UserSteamID  int64  `json:"user_steam_id"`
    Created      int64  `json:"created"`
    Name         string `json:"name"`
    Item         string `json:"item"`
    SetCount     int    `json:"set_count"`
    BuyCount     int    `json:"buy_count"`
    Status       string `json:"status"`
    Description  string `json:"desc"`
    Type         string `json:"type"`
    Price        float64 `json:"price"`
    Amount       int    `json:"amount"`
    Currency     string `json:"cy"`
    Delivery     string `json:"delivery"`
    Player       string `json:"player"`
    Claim        int    `json:"claim"`
}
type MatchedData struct {
    OfferID      int    `json:"offer_id"`
    UserID   string `json:"userId"`
    Item     string `json:"item"`
	Count     int    `json:"count"`
    ID       int    `json:"id"`
}
const (
    baseURL     = "https://api.wargm.ru/v1/shop/"
    serverID    = "3597"
    serverAPIKey = "fEj7ngwolYgUWEh_u5h(yhQYgMKl2qcvrvn)BL7vFH_Pn_f0"
)

var log = logger.NewInfoLogger()
// GetOpenPurchases возвращает список открытых покупок пользователя.
func getOpenPurchases() ([]Purchase, error) {
    url := fmt.Sprintf("%soperations?client=%s:%s&status=pending", baseURL, serverID, serverAPIKey)
    resp, err := http.Get(url)
    if err != nil {
        log.Error("Ошибка при запросе списка открытых покупок: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    var response struct {
        Responce struct {
            Status    string          `json:"status"`
            Message   string          `json:"msg"`
            Data      json.RawMessage `json:"data"`
            DataCount int             `json:"data_count"`
        } `json:"responce"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        log.Error("Ошибка при декодировании ответа: %v", err)
        return nil, err
    }

    if response.Responce.Status != "ok" {
        return nil, fmt.Errorf("ошибка в ответе API: %s", response.Responce.Message)
    }

    // Проверка на пустой объект "{}"
    if string(response.Responce.Data) == "{}" {
        return []Purchase{}, nil
    }

    // Декодируем данные покупок
    purchases := make(map[string]Purchase)
    if err := json.Unmarshal(response.Responce.Data, &purchases); err != nil {
        log.Error("Ошибка при декодировании данных о покупках: %v", err)
        return nil, err
    }

    // Преобразуем map в срез
    var purchasesSlice []Purchase
    for _, purchase := range purchases {
        purchasesSlice = append(purchasesSlice, purchase)
    }

    return purchasesSlice, nil
}
// ConfirmPurchase подтверждает покупку с указанным ID операции.
func confirmPurchase(operationID int) error {
    url := fmt.Sprintf("%soperation_success?client=%s:%s&operation_id=%d", baseURL, serverID, serverAPIKey, operationID)
	
    resp, err := http.Get(url)
    if err != nil {
        log.Error("Ошибка при запросе подтверждения покупки: %v", err)
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Error("Не удалось подтвердить покупку: статус %d", resp.StatusCode)
        return fmt.Errorf("не удалось подтвердить покупку: статус %d", resp.StatusCode)
    }

    log.Info("Покупка успешно подтверждена: operationID=%s", operationID)
    return nil
}


type MatchedPurchase struct {
    UserID    string // ID пользователя
    Item      string // Наименование товара
    SetCount  int    // Количество товаров
}

func findMatchedPurchases(players []restjs.Player, purchases []Purchase) []MatchedData {
    var matchedData []MatchedData

    // Перебор игроков из restjs
    for _, player := range players {
        playerName := player.Name

        // Поиск соответствующей покупки в purchases
        for _, purchase := range purchases {
            if playerName == purchase.Player {
                // Если найдено совпадение, добавляем данные в массив matchedData
                matchedData = append(matchedData, MatchedData{
                    UserID:   player.UserID,
                    Item:     purchase.Item,
                    OfferID: purchase.OfferID,
                    Count: purchase.BuyCount * purchase.SetCount,
                    ID:       purchase.ID,
                })
                break // Прерываем внутренний цикл, чтобы не искать дальше
            }
        }
    }

    return matchedData
}

func giveAndConfirmItems(matchedPurchases []MatchedData) {
    for _, data := range matchedPurchases {
        if data.Item != "" {
		processMatchedData(data)
        confirmPurchase(data.ID)
        }
    }
}

func Handler(players []restjs.Player) {
    // Получение списка покупок из wargmshop
	        log.Info("test")
    purchases, err := getOpenPurchases()
    if err != nil {
        log.Error("Ошибка при получении списка покупок из wargmshop:", err)
		        fmt.Printf("", err)
        return
    }
    matchedPurchases := findMatchedPurchases(players, purchases)
    giveAndConfirmItems(matchedPurchases)
}
func processMatchedData(matchedData MatchedData) {
    switch matchedData.Item {
    case "exp_10k":
	count := matchedData.Count * 10000
            givepak.GiveExp(matchedData.UserID, count)
    case "exp_10k2":
fmt.Printf("dsadad")
    default:
            givepak.GiveItem(matchedData.Item, matchedData.UserID, matchedData.Count)
    }
}

