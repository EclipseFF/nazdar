package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Deal struct {
	Fields DealFields `json:"fields"`
}

type DealFields struct {
	Title        string `json:"TITLE"`
	StageID      string `json:"STAGE_ID"`
	Opportunity  int    `json:"OPPORTUNITY"`
	CurrencyID   string `json:"CURRENCY_ID"`
	ContactID    string `json:"CONTACT_ID,omitempty"`
	CompanyID    string `json:"COMPANY_ID,omitempty"`
	Opened       string `json:"OPENED"`
	AssignedByID string `json:"ASSIGNED_BY_ID"`
}

func SendApiReq(user *User, items []*CartItem) error {
	webhookURL := "https://nazdar.bitrix24.kz/rest/60087/0l0bq8l6noka8xx4/crm.deal.add.json"
	body := " "
	totalPrice := 0
	for _, item := range items {
		body += *item.Name + " " + fmt.Sprintf("%d", *item.Count) + " шт. "
		totalPrice += *item.Price * *item.Count
	}

	title := "Заказ с сайта, номер " + *user.Phone + " имя " + *user.Name + body + " " + "Сумма: " + fmt.Sprintf("%d", totalPrice)
	deal := Deal{
		Fields: DealFields{
			Title:       title,
			StageID:     "NEW",
			Opportunity: totalPrice,
			CurrencyID:  "KZT",
			Opened:      "Y",
		},
	}

	jsonData, err := json.Marshal(deal)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New("error sending request")
	}
}
