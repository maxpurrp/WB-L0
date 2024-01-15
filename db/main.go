package db

import (
	"client/order"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	connStr = "postgresql://max:1234@localhost:5432/WB?sslmode=disable"
	db      = openConn()
	ordr    = len(GetAllOrders())
)

func openConn() *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func saveOrder(db *sql.DB, data *order.Order) {
	_, err := db.Exec(`
	INSERT INTO orders (order_uid, track_number, entry, delivery_name, delivery_phone, delivery_zip,
	delivery_city, delivery_address, delivery_region, delivery_email, locale, internal_signature, customer_id,
	delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		data.OrderUID, data.TrackNumber, data.Entry, data.Delivery.Name, data.Delivery.Phone, data.Delivery.Zip,
		data.Delivery.City, data.Delivery.Address, data.Delivery.Region, data.Delivery.Email, data.Locale, data.InternalSig, data.CustomerID,
		data.DeliveryService, data.ShardKey, data.SMID, data.DateCreated, data.OofShard,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func savePayments(db *sql.DB, id int, data *order.Order) {
	_, err := db.Exec(`
	INSERT INTO payments (order_id, transaction, request_id, currency, provider, amount,
	payment_dt, bank, delivery_cost, goods_total, custom_fee)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		ordr+1, data.Payment.Transaction, data.Payment.RequestID, data.Payment.Currency,
		data.Payment.Provider, data.Payment.Amount, data.Payment.PaymentDt, data.Payment.Bank,
		data.Payment.DeliveryCost, data.Payment.GoodsTotal, data.Payment.CustomFee,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func saveItems(db *sql.DB, id int, data *order.ItemData) {
	_, err := db.Exec(`
	INSERT INTO order_items (order_id, chrt_id, track_number, price, rid, name, sale, size,
	total_price, nm_id, brand, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		ordr+1, data.ChrtID, data.TrackNumber, data.Price,
		data.RID, data.Name, data.Sale, data.Size,
		data.TotalPrice, data.NMID, data.Brand, data.Status,
	)

	if err != nil {
		log.Fatal(err)
	}
}
func SaveInDB(data *order.Order) {
	err := db.Ping()
	if err != nil {
		db = openConn()
	}
	saveOrder(db, data)
	savePayments(db, ordr+1, data)
	for i := 0; i < len(data.Items); i++ {
		item := data.Items[i]
		saveItems(db, ordr+1, &item)
	}
	ordr++
	defer db.Close()
}

func getOrderData(uid string, db *sql.DB) order.Order {
	var orderData order.Order
	err := db.QueryRow(`
        SELECT *
        FROM orders
        WHERE order_uid = $1
    `, uid).Scan(
		&orderData.OrderUID, &orderData.OrderUID, &orderData.TrackNumber, &orderData.Entry, &orderData.Delivery.Name,
		&orderData.Delivery.Phone, &orderData.Delivery.Zip, &orderData.Delivery.City, &orderData.Delivery.Address,
		&orderData.Delivery.Region, &orderData.Delivery.Email, &orderData.Locale, &orderData.InternalSig,
		&orderData.CustomerID, &orderData.DeliveryService, &orderData.ShardKey, &orderData.SMID,
		&orderData.DateCreated, &orderData.OofShard,
	)
	if err != nil {
		log.Fatal(err)
	}
	return orderData
}

func getPaymentData(id string, db *sql.DB) (order.PaymentData, int) {
	var paymentData order.PaymentData
	var linkID int
	var garbage_1 string
	rows, err := db.Query(`
			SELECT *
			FROM payments
			WHERE transaction = $1;
			`, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&garbage_1, &linkID, &paymentData.Transaction, &paymentData.RequestID,
			&paymentData.Currency, &paymentData.Provider, &paymentData.Amount, &paymentData.PaymentDt,
			&paymentData.Bank, &paymentData.DeliveryCost, &paymentData.GoodsTotal, &paymentData.CustomFee,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
	return paymentData, linkID
}

func getitemsData(linkID int, db *sql.DB) []order.ItemData {
	var ItemsData []order.ItemData
	var itemData order.ItemData
	var garbage int
	var garbage_1 int
	rows, err := db.Query(`
    SELECT *
    FROM order_items
    WHERE order_id = $1`, linkID)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&garbage, &garbage_1, &itemData.ChrtID, &itemData.TrackNumber, &itemData.Price, &itemData.RID,
			&itemData.Name, &itemData.Sale, &itemData.Size, &itemData.TotalPrice, &itemData.NMID,
			&itemData.Brand, &itemData.Status,
		)
		if err != nil {
			log.Fatal(err)
		}
		ItemsData = append(ItemsData, itemData)
	}
	return ItemsData
}

func GetRecord(orderUID string) order.Order {
	db := openConn()
	orderData := getOrderData(orderUID, db)
	paymentsData, link := getPaymentData(orderData.OrderUID, db)
	itemData := getitemsData(link, db)
	record := order.Order{
		OrderUID:        orderData.OrderUID,
		TrackNumber:     orderData.TrackNumber,
		Entry:           orderData.Entry,
		Delivery:        orderData.Delivery,
		Payment:         paymentsData,
		Items:           itemData,
		Locale:          orderData.Locale,
		InternalSig:     orderData.InternalSig,
		CustomerID:      orderData.CustomerID,
		DeliveryService: orderData.DeliveryService,
		ShardKey:        orderData.ShardKey,
		SMID:            orderData.SMID,
		DateCreated:     orderData.DateCreated,
		OofShard:        orderData.OofShard,
	}
	return record

}

func GetAllOrders() []order.Order {
	db = openConn()
	var orderModel order.Order
	var orderUIDs []string
	var orders []order.Order
	row, err := db.Query(`SELECT order_uid from orders;`)
	if err != nil {
		fmt.Println(err)
	}
	for row.Next() {
		row.Scan(&orderModel.OrderUID)
		orderUIDs = append(orderUIDs, orderModel.OrderUID)
	}
	for i := 0; i < len(orderUIDs); i++ {
		record := GetRecord(orderUIDs[i])
		orders = append(orders, record)
	}
	return orders
}
