package portal

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	givsoft "giv/givsoft"
	"io"
	"log"
	"net/http"
	net_url "net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/peterbourgon/diskv/v3"
	Jalaali "github.com/yaa110/go-persian-calendar"
)

var DB *diskv.Diskv

var base_url string = "https://batkap.com/site/"

type Order_result struct {
	Success bool `json:"success"`
	Order   struct {
		ID             int      `json:"id"`
		Description    any      `json:"description"`
		Status         []string `json:"status"`
		Quantity       int      `json:"quantity"`
		Weight         int      `json:"weight"`
		Shipping       int      `json:"shipping"`
		Subtotal       int      `json:"subtotal"`
		Discount       int      `json:"discount"`
		Tax            int      `json:"tax"`
		Price          int      `json:"price"`
		RemainingPrice int      `json:"remaining_price"`
		IP             string   `json:"ip"`
		Contact        struct {
			Name    string `json:"name"`
			Mobile  string `json:"mobile"`
			Phone   any    `json:"phone"`
			Email   any    `json:"email"`
			Country struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"country"`
			State struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"state"`
			City struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"city"`
			Zipcode   string  `json:"zipcode"`
			Address   string  `json:"address"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"contact"`
		Items []struct {
			Variant *struct {
				ID  int `json:"id"`
				Sku any `json:"sku"`
			} `json:"variant"`
			Product *struct {
				ID int `json:"id"`
			} `json:"product"`
			Title    string  `json:"title"`
			Price    int     `json:"price"`
			Quantity int     `json:"quantity"`
			Weight   any     `json:"weight"`
			Shipping any     `json:"shipping"`
			Discount any     `json:"discount"`
			Sku      *string `json:"sku"`
			Tax      int     `json:"tax"`
		} `json:"items"`
		Coupons   any `json:"coupons"`
		Shipments any `json:"shipments"`
		Payments  []*struct {
			ID          int      `json:"id"`
			Description any      `json:"description"`
			ReferenceID string   `json:"reference_id"`
			Type        string   `json:"type"`
			Status      []string `json:"status"`
			SubStatus   any      `json:"sub_status"`
			Amount      int      `json:"amount"`
			Created     struct {
				Year      string `json:"year"`
				Month     string `json:"month"`
				MonthName string `json:"month_name"`
				Day       string `json:"day"`
				Date      string `json:"date"`
				Time      string `json:"time"`
				Universal string `json:"universal"`
				Timestamp int    `json:"timestamp"`
				Subtract  string `json:"subtract"`
				Past      bool   `json:"past"`
			} `json:"created"`
			Gateway struct {
				ID    int    `json:"id"`
				Title string `json:"title"`
				Type  string `json:"type"`
				Owner string `json:"owner"`
			} `json:"gateway"`
		} `json:"payments"`
		User *struct {
			ID           int    `json:"id"`
			Username     string `json:"username"`
			Name         any    `json:"name"`
			Nickname     any    `json:"nickname"`
			NationalCode any    `json:"national_code"`
			Avatar       any    `json:"avatar"`
		} `json:"user"`
		Label         any `json:"label"`
		ShippingClass struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
			Type  string `json:"type"`
		} `json:"shipping_class"`
		Created struct {
			Year      string `json:"year"`
			Month     string `json:"month"`
			MonthName string `json:"month_name"`
			Day       string `json:"day"`
			Date      string `json:"date"`
			Time      string `json:"time"`
			Universal string `json:"universal"`
			Timestamp int    `json:"timestamp"`
			Subtract  string `json:"subtract"`
			Past      bool   `json:"past"`
		} `json:"created"`
		DueDate struct {
			Year      string `json:"year"`
			Month     string `json:"month"`
			MonthName string `json:"month_name"`
			Day       string `json:"day"`
			Date      string `json:"date"`
			Time      string `json:"time"`
			Universal string `json:"universal"`
			Timestamp int    `json:"timestamp"`
			Subtract  string `json:"subtract"`
			Past      bool   `json:"past"`
		} `json:"due_date"`
		Delivery any `json:"delivery"`
		Updated  any `json:"updated"`
	} `json:"order"`
}
type Order struct {
	ID          int      `json:"id"`
	Description any      `json:"description"`
	Status      []string `json:"status"`
	Quantity    int      `json:"quantity"`
	Weight      int      `json:"weight"`
	Shipping    int      `json:"shipping"`
	Subtotal    int      `json:"subtotal"`
	Discount    int      `json:"discount"`
	Tax         int      `json:"tax"`
	Price       int      `json:"price"`
	Payments    int      `json:"payments"`
	IP          string   `json:"ip"`
	User        struct {
		ID           int    `json:"id"`
		Username     string `json:"username"`
		Name         any    `json:"name"`
		Nickname     any    `json:"nickname"`
		NationalCode any    `json:"national_code"`
		Avatar       any    `json:"avatar"`
	} `json:"user"`
	Contact struct {
		Name    string `json:"name"`
		Mobile  string `json:"mobile"`
		Phone   any    `json:"phone"`
		Email   any    `json:"email"`
		Country struct {
			ID        int     `json:"id"`
			Name      string  `json:"name"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"country"`
		State struct {
			ID        int     `json:"id"`
			Name      string  `json:"name"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"state"`
		City struct {
			ID        int     `json:"id"`
			Name      string  `json:"name"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"city"`
		Zipcode   string  `json:"zipcode"`
		Address   string  `json:"address"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"contact"`
	Label         any `json:"label"`
	ShippingClass struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"shipping_class"`
	Created struct {
		Year      string `json:"year"`
		Month     string `json:"month"`
		MonthName string `json:"month_name"`
		Day       string `json:"day"`
		Date      string `json:"date"`
		Time      string `json:"time"`
		Universal string `json:"universal"`
		Timestamp int    `json:"timestamp"`
		Subtract  string `json:"subtract"`
		Past      bool   `json:"past"`
	} `json:"created"`
	DueDate struct {
		Year      string `json:"year"`
		Month     string `json:"month"`
		MonthName string `json:"month_name"`
		Day       string `json:"day"`
		Date      string `json:"date"`
		Time      string `json:"time"`
		Universal string `json:"universal"`
		Timestamp int    `json:"timestamp"`
		Subtract  string `json:"subtract"`
		Past      bool   `json:"past"`
	} `json:"due_date"`
	Delivery any `json:"delivery"`
	Updated  any `json:"updated"`
}
type Orders struct {
	Success bool     `json:"success"`
	Total   int      `json:"total"`
	Count   int      `json:"count"`
	Orders  []*Order `json:"orders"`
}
type session_resault struct {
	Success     bool   `json:"success"`
	Description string `json:"description"`
	Token       string `json:"token"`
}

func Make_session() string {
	url := base_url + "api/v1/user/create-session"
	method := "POST"
	user := os.Getenv("PORTAL_USER")
	password := os.Getenv("PORTAL_PASS")
	payload_string := fmt.Sprintf(`{
		"username": "%s",
		"password": "%s"
	}`, user, password)
	payload := strings.NewReader(payload_string)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	req.Close = true
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "Application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	var Token_res session_resault
	err = dec.Decode(&Token_res)
	if err != nil {
		body, _ := io.ReadAll(res.Body)
		fmt.Println(res.StatusCode, "|-|", string(body))
		log.Printf("there was an error while decoding %+v\n", err)
		os.Exit(1)
	}
	if res.StatusCode != 200 {
		log.Println(url, payload_string)
		log.Printf("Token Not Found %s", res.Status)
		os.Exit(-1)
	} else {
		return Token_res.Token
	}
	return ""
}
func Get_orders(token string, wg *sync.WaitGroup) {
	defer wg.Done()
	todayJ := Jalaali.Now().AddDate(0, 0, -1).Format("yyy/MM/dd")
	status := []string{"paid", "cash_on_delivery"}
	for _, s := range status {
		url := base_url + fmt.Sprintf("api/v1/manage/store/orders?page=1&size=20&status=%s&payment=&start=%s&end=&label_id=&user_id=&shipping_id=&ip=&keywords=", s, todayJ)
		method := "GET"
		fmt.Print(url)
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-Type", "Application/json")
		req.Close = true
		if err != nil {
			log.Println(err)
			return
		}
		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		defer res.Body.Close()

		if err != nil {
			log.Printf("somthing is wrong with the Body %s\n", err)
			return
		}

		decoder := json.NewDecoder(res.Body)
		var orders Orders
		err = decoder.Decode(&orders)
		if err != nil {
			log.Printf("There was an error while decoding orders %s\n", err)
			return
		}
		lastPortalPurchase, _ := DB.Read("LAST_PORTAL_PURCHASE")
		lastPortalPurchaseValue := binary.LittleEndian.Uint32(lastPortalPurchase)

		if orders.Count > 0 {
			for _, order := range orders.Orders {
				if uint32(order.ID) == lastPortalPurchaseValue {
					break
				}
				wg.Add(1)
				go get_Order(token, order.ID, wg)
			}
			buf := make([]byte, 4) // adjust size according to your int type
			binary.LittleEndian.PutUint32(buf, uint32(orders.Orders[0].ID))
			DB.Write("LAST_PORTAL_PURCHASE", buf)
		}
	}
}
func get_Order(token string, order_id int, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("https://batkap.com/site/api/v1/manage/store/orders/%d", order_id)
	fmt.Printf("Getting order Detail %d\n", order_id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "Application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	decder := json.NewDecoder(res.Body)
	var order_resault Order_result
	decder.Decode(&order_resault)
	if err != nil {
		body, _ := json.Marshal(order_resault)
		log.Println(body)
		log.Println(err)
		return
	}
	if order_resault.Success {
		var order_detail givsoft.Order_detail
		var ItemDetail []givsoft.Itemdetail
		var date_created_formated string
		if order_resault.Order.Payments != nil {
			date_created, _ := time.Parse("02/01/2006 15:04:05", order_resault.Order.Payments[0].Created.Universal)
			date_created_formated = date_created.Format("2006/01/02 15:04:05")
		} else {
			date_created_formated = time.Now().Format("2006/01/02 15:04:05")
		}
		date := Jalaali.Now().Format("yyyyMMdd")
		total_price := 0
		for _, item := range order_resault.Order.Items {
			if item.Sku != nil {
				total_price += item.Price
				itemId, _ := strconv.ParseInt(*item.Sku, 10, 64)
				ItemDetail = append(ItemDetail, givsoft.Itemdetail{
					ItemDetailID: itemId,
					OrderID:      order_id,
					ItemID:       itemId,
					ItemBarcode:  *item.Sku,
					Quantity:     item.Quantity,
					Fee:          item.Price,
					DateCreated:  date_created_formated,
					DateChanged:  date_created_formated,
				})
			}
		}
		var PersonId string
		if order_resault.Order.User != nil {
			PersonId = strconv.FormatInt(int64(order_resault.Order.User.ID), 10)
			if DB.Has(PersonId) {
				PersonIdb, err := DB.Read(PersonId)
				if err != nil {
					log.Printf("Somthing Wen wrong oops %s\n", err)
				}
				PersonId = strconv.FormatUint(uint64(binary.LittleEndian.Uint32(PersonIdb)), 10)
				log.Println("we have PersonID" + PersonId)
			} else {
				log.Println("we don't have PersonID ")
				PersonId = givsoft.Create_customer(strconv.FormatInt(int64(order_resault.Order.User.ID), 10), order_resault.Order.Contact.Name, order_resault.Order.Contact.City.Name, order_resault.Order.Contact.Address, order_resault.Order.Contact.Mobile, order_resault.Order.Contact.Zipcode)
			}
		} else {
			PersonId = "ناشناخته"
		}
		if PersonId == "" {
			log.Println(PersonId)
			os.Exit(-1)
			return
		}
		order_detail = givsoft.Order_detail{
			OrderID:            -1,
			SourceID:           order_resault.Order.ID,
			Type:               "SALE",
			No:                 order_resault.Order.ID,
			Date:               date,
			EffectiveDate:      date,
			PersonID:           PersonId,
			CouponCode:         "",
			Description:        "سفارش خرید از پرتال",
			TotalQuantity:      order_resault.Order.Quantity,
			TotalPrice:         total_price,
			TotalDiscount:      0,
			PackingCost:        0,
			TransferCost:       0,
			PostRefCode:        strconv.FormatInt(int64(order_resault.Order.ID), 10),
			ReceiverName:       order_resault.Order.Contact.Name,
			ReceiverCity:       order_resault.Order.Contact.City.Name,
			ReceiverAddress:    order_resault.Order.Contact.Address,
			ReceiverMobile:     order_resault.Order.Contact.Mobile,
			ReceiverPostalCode: order_resault.Order.Contact.Zipcode,
			PaymentType:        "",
			PaymentStatus:      "",
			DateCreated:        date_created_formated,
			DateChanged:        date_created_formated,
		}
		if len(order_resault.Order.Payments) > 0 {
			order_detail.PaymentBank = order_resault.Order.Payments[0].Gateway.Title
			order_detail.PaymentType = order_resault.Order.Payments[0].Gateway.Type
			order_detail.PaymentStatus = order_resault.Order.Payments[0].Status[0]
			order_detail.PaymentBankRefCode = order_resault.Order.Payments[0].ReferenceID
		} else {
			order_detail.PaymentType = "ONLINE"
			order_detail.PaymentStatus = "PAYMENT_STATUS_SUCCESSFUL"
		}
		order_detail.ItemDetail = ItemDetail
		wg.Add(1)
		go givsoft.Make_Order(order_detail, wg)
	} else {
		fmt.Println(order_resault.Success)
	}
}
func Update_Product(token string, product_id int, stock int, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("https://batkap.com/site/api/v1/manage/store/products/variants/%d", product_id)
	method := "PUT"
	ch := make(chan Product_resault)
	wg.Add(1)
	go get_product(product_id, token, wg, ch)
	product := <-ch
	close(ch)
	product.Variant.Stock = stock
	product_byte, err := json.Marshal(product.Variant)
	if err != nil {
		log.Println(string(product_byte))
		log.Fatal(err)
	}
	payload := bytes.NewReader(product_byte)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
}
func Quantityonhand(token string, wg *sync.WaitGroup) {
	defer wg.Done()
	lastgiv, _ := DB.Read("LAST_GIV_PURCHASE")
	url := fmt.Sprintf("http://91.92.214.97:8201/api/quantityonhand?lastDate=%s", net_url.QueryEscape(string(lastgiv)))
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Close = true
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("WEB_TOKEN", os.Getenv("WEB_TOKEN"))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	var qoh givsoft.QuantityOnhand
	if err != nil {
		log.Printf("Error while Setting Env variable %s\n", err)
	}
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&qoh)
	err = DB.Write("LAST_GIV_PURCHASE", []byte(qoh.LastDatetime))
	if err != nil {
		log.Print(err)
	}
	qohJ, _ := json.Marshal(qoh)
	log.Print(string(qohJ))
	for _, item := range qoh.Value {
		wg.Add(1)
		ProductID, err := DB.Read(strconv.Itoa(int(item.ItemID))) //getting ItemId from Giv map to Portal Item Id",
		if err != nil {
			wg.Done()
			continue
		}
		time.Sleep(time.Millisecond * 500)
		ProductId := binary.LittleEndian.Uint64(ProductID)
		log.Println(string(item.ItemID), int(ProductId))
		if err != nil {
			wg.Done()
			log.Printf("error while Updating order from giv to portal %s\n", err)
			continue
		}
		go Update_Product(token, int(ProductId), int(item.ItemQuantityOnHand), wg)
	}
}
func QuantityOnhand_byitem(token string, itemId string, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("http://91.92.214.97:8201/api/quantityonhand?itemid=%s", itemId)
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	req.Close = true
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("WEB_TOKEN", os.Getenv("WEB_TOKEN"))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	if err != nil {
		log.Printf("Error while Setting Env variable %s\n", err)
		os.Exit(1)
	}
	var qoh givsoft.QuantityOnhand
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&qoh)
	err = DB.Write("LAST_GIV_PURCHASE", []byte(qoh.LastDatetime))
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	if len(qoh.Value) == 0 {
		log.Println(url)
		return
	}
	item := qoh.Value[0]
	ProductID, err := DB.Read(itemId) //getting ItemId from Giv map to Portal Item Id",
	wg.Add(1)
	if err != nil {
		wg.Done()
		log.Printf("%+v", err)
		log.Printf("%d itemID not Found ", item.ItemID)
		return
	}
	ProductId := binary.LittleEndian.Uint64(ProductID)
	if err != nil {
		wg.Done()
		log.Printf("error while Updating order from giv to portal %s\n", err)
		return
	}
	go Update_Product(token, int(ProductId), int(item.ItemQuantityOnHand), wg)
}
