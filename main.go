package main

import (
	"giv/givsoft"
	"giv/portal"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/peterbourgon/diskv/v3"
)

var env_vars = [7]string{
	"WEB_TOKEN",
	"ITEM_DETAIL_ID",
	"PORTAL_PASS",
	"PORTAL_USER",
	"LAST_PORTAL_PURCHASE",
	"LAST_GIV_PURCHASE",
	"SUCCESS",
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error while Loading .env file  %s \n", err)
		os.Exit(1)
	}
	f, err := os.OpenFile("Log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Couldn't Open Log.txt %s \n", err)
		os.Exit(1)
	}
	defer f.Close()
	w := io.MultiWriter(os.Stdout, f)
	log.SetOutput(w)
	flatTransform := func(s string) []string { return []string{} }
	db := diskv.New(diskv.Options{
		BasePath:     "portal_DB",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	portal.DB = db
	givsoft.DB = db
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err != nil {
		log.Printf("Error:error while  loading .env %s\n", err)
	}
	token := portal.Make_session()
	log.Println(token)
	var wg sync.WaitGroup
	wg.Add(1)
	go portal.Get_orders(token, &wg)
	portal.Update_giv(token, 0)
	for range time.Tick(time.Second * 10) {
		log.Printf("Syncing Begineing ...\n")
		wg.Add(1)
		go portal.Get_orders(token, &wg) //Syncing giv Items via portal orders
		portal.Update_giv(token, 0)      //updating portal Product with  giv quantity on hand
		wg.Wait()
	}
}
