package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type OrderItem struct {
	Quantity int    `json:"quantity"`
	Title    string `json:"title"`
}

type Order struct {
	ID           int         `json:"id"`
	CustomerName string      `json:"customerName"`
	Note         string      `json:"note"`
	Items        []OrderItem `json:"items"`
}

type PrintRequest struct {
	PrinterIP string `json:"printerIp"`
	StoreName string `json:"storeName"`
	Order     Order  `json:"order"`
}

const printerPort = "9100"

func buildReceipt(storeName string, order Order) string {
	esc := "\x1b"
	gs := "\x1d"

	if storeName == "" {
		storeName = "سفارش"
	}

	receipt := esc + "@"
	receipt += esc + "a" + "\x01"
	receipt += gs + "!" + "\x11"
	receipt += storeName + "\n"
	receipt += gs + "!" + "\x00"
	receipt += "--------------------------------\n"
	receipt += esc + "a" + "\x00"
	receipt += fmt.Sprintf("سفارش شماره: %d\n", order.ID)
	receipt += fmt.Sprintf("مشتری: %s\n", order.CustomerName)
	receipt += "--------------------------------\n"

	for _, item := range order.Items {
		receipt += fmt.Sprintf("%d x %s\n", item.Quantity, item.Title)
	}

	receipt += "--------------------------------\n"
	if order.Note != "" {
		receipt += fmt.Sprintf("یادداشت: %s\n", order.Note)
	}
	receipt += "\n\n\n"
	receipt += gs + "V" + "\x00"

	return receipt
}

func sendToPrinter(printerIP string, data string) error {
	conn, err := net.DialTimeout("tcp", printerIP+":"+printerPort, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(data))
	return err
}

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	http.HandleFunc("/print", func(w http.ResponseWriter, r *http.Request) {
		var req PrintRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, `{"error":"درخواست نامعتبر"}`, http.StatusBadRequest)
			return
		}

		receipt := buildReceipt(req.StoreName, req.Order)

		err = sendToPrinter(req.PrinterIP, receipt)
		if err != nil {
			fmt.Println("خطای چاپ:", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		fmt.Printf("سفارش #%d با موفقیت چاپ شد.\n", req.Order.ID)
		fmt.Fprintf(w, `{"success":true}`)
	})

	fmt.Println("سرور روشن شد رو http://localhost:9123")
	http.ListenAndServe("127.0.0.1:9123", nil)
}