package portal

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Update_Resault struct {
	Success  bool `json:"success"`
	Total    int  `json:"total"`
	Count    int  `json:"count"`
	Variants []struct {
		ID           int      `json:"id"`
		ProductID    int      `json:"product_id"`
		Title        string   `json:"title"`
		Price        int      `json:"price"`
		ComparePrice any      `json:"compare_price"`
		Tax          any      `json:"tax"`
		Shipping     any      `json:"shipping"`
		Weight       any      `json:"weight"`
		Length       any      `json:"length"`
		Width        any      `json:"width"`
		Height       any      `json:"height"`
		Stock        int      `json:"stock"`
		Minimum      any      `json:"minimum"`
		Maximum      any      `json:"maximum"`
		Sku          string   `json:"sku"`
		Image        any      `json:"image"`
		Type         string   `json:"type"`
		Status       []string `json:"status"`
		Files        any      `json:"files"`
	} `json:"variants"`
}
type Product_resault struct {
	Success bool `json:"success"`
	Variant struct {
		ID           int      `json:"id"`
		ProductID    int      `json:"product_id"`
		Title        string   `json:"title"`
		Price        int      `json:"price"`
		ComparePrice any      `json:"compare_price"`
		Tax          any      `json:"tax"`
		Shipping     any      `json:"shipping"`
		Weight       any      `json:"weight"`
		Length       any      `json:"length"`
		Width        any      `json:"width"`
		Height       any      `json:"height"`
		Stock        int      `json:"stock"`
		Minimum      any      `json:"minimum"`
		Maximum      any      `json:"maximum"`
		Sku          *string  `json:"sku"`
		Image        any      `json:"image"`
		Type         string   `json:"type"`
		Status       []string `json:"status"`
		Files        any      `json:"files"`
	} `json:"variant"`
}

func Update_giv(token string, page int) {
	url := fmt.Sprintf("https://batkap.com/site/api/v1/manage/store/products/variants?size=100&page=%d", page+1)
	method := "GET"
	fmt.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body := res.Body
	decoder := json.NewDecoder(res.Body)
	var Update_Resault Update_Resault
	decoder.Decode(&Update_Resault)
	if !Update_Resault.Success {
		b, _ := io.ReadAll(body)
		fmt.Print(string(b))
		return
	}
	if (page+1)*100 < Update_Resault.Total {
		Update_giv(token, page+1)
	}
	wg := sync.WaitGroup{}
	ch := make(chan Product_resault, len(Update_Resault.Variants))
	for _, p := range Update_Resault.Variants {
		wg.Add(1)
		time.Sleep(time.Millisecond * 400)
		go get_product(p.ID, token, &wg, ch)
	}
	close(ch)
	for prod := range ch {
		if prod.Variant.Sku != nil {
			log.Println(*prod.Variant.Sku, prod.Variant.Title)
			wg.Add(1)
			go QuantityOnhand_byitem(token, *prod.Variant.Sku, &wg)
		}
	}
	wg.Wait()
}
func get_product(id int, token string, wg *sync.WaitGroup, ch chan Product_resault) {
	url := fmt.Sprintf("https://batkap.com/site/api/v1/manage/store/products/variants/%d", id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var Product Product_resault
	decoder.Decode(&Product)
	key := Product.Variant.Sku
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(Product.Variant.ID)) //saving product variants
	if key != nil {
		DB.Write(*key, buf)
		DB.Write(string(Product.Variant.ID), []byte(*key))
	}
	if !Product.Success {
		fmt.Println(res.Status)
	}
	ch <- Product
}
