package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"mp-generate-test-data/config"
	"mp-generate-test-data/internal/constants"
	"mp-generate-test-data/internal/models"
	"mp-generate-test-data/internal/stores"
	"mp-generate-test-data/internal/utils"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	cfg                *config.Config
	mainStore          *stores.MainStore
	ordersStore        *stores.OrdersStore
	orderAssetsStore   *stores.OrderAssetsStore
	matchedOrdersStore *stores.MatchedOrdersStore
}

func New(
	config *config.Config,
	mainStore *stores.MainStore,
	ordersStore *stores.OrdersStore,
	orderAssetsStore *stores.OrderAssetsStore,
	matchedOrdersStore *stores.MatchedOrdersStore,
) *Service {
	s := &Service{
		cfg:                config,
		mainStore:          mainStore,
		ordersStore:        ordersStore,
		orderAssetsStore:   orderAssetsStore,
		matchedOrdersStore: matchedOrdersStore,
	}

	return s
}

func (s *Service) Run() {
	// s.gen_test_data()
	sendLargeRequest()
}

func (s *Service) genOrdersWithMaxGoRoutines(fromOrderId int64) {
	var wga sync.WaitGroup
	for i := 0; i < constants.MAX_GO_NUM; i++ {
		no := i
		wga.Add(1)
		go func() {
			defer wga.Done()

			err := s.gen_orders(fromOrderId + int64(constants.BATCH_SIZE*no))
			if err != nil {
				fmt.Printf("error here %s\n", err.Error())
				return
			}

			// fmt.Printf("@@@Done order id: %d\n", fromOrderId+int64(constants.BATCH_SIZE*no))
		}()
	}

	wga.Wait()
}

func (s *Service) gen_test_data() {
	// get the largest id in table orders
	order, exist, err := s.ordersStore.GetLargestOrderId()
	if err != nil {
		fmt.Println("error while getting largest order id: " + err.Error())
		return
	}

	fromOrderId := int64(1)
	if exist {
		fromOrderId = order.Id + 1
	}

	for i := 0; i < constants.BATCH_NUM; i++ {
		t1 := time.Now().Unix()
		s.genOrdersWithMaxGoRoutines(fromOrderId + int64(constants.BATCH_SIZE*constants.MAX_GO_NUM*i))
		fmt.Printf("Done order id: %d\n", fromOrderId+int64(constants.BATCH_SIZE*constants.MAX_GO_NUM*i))
		t2 := time.Now().Unix()
		fmt.Printf("about: %d\n", t2-t1)
	}

}

func (s *Service) gen_orders(fromOrderId int64) error {
	// generate 100(batch_size) orders
	orders := make([]*models.OrdersMd, constants.BATCH_SIZE)
	var orderAssets []*models.OrderAssetsMd
	var matchedOrders []*models.MatchedOrdersMd

	// random here
	for i := int64(0); i < constants.BATCH_SIZE; i++ {
		// number of asset in this order
		numOfAssets := utils.RandNumberOfAssets()
		// gen Assets
		assetsJsonB := make([]*models.Assets, numOfAssets)
		oaOfThisOrder := make([]*models.OrderAssetsMd, numOfAssets)
		orderId := fromOrderId + i

		// gen order assets
		for j := 0; j < numOfAssets; j++ {
			assetsJsonB[j] = &models.Assets{
				Address:  utils.RandVerifiedContract(),
				Erc:      int8(utils.RandErcType()),
				Id:       rand.Int63(),
				Quantity: utils.RandQuantity(),
			}
		}

		for k := 0; k < numOfAssets; k++ {
			asset := &models.OrderAssetsMd{
				OrderId:  orderId,
				Address:  assetsJsonB[k].Address,
				Erc:      assetsJsonB[k].Erc,
				Id:       assetsJsonB[k].Id,
				Quantity: int64(assetsJsonB[k].Quantity),
			}
			oaOfThisOrder[k] = asset
		}
		orderAssets = append(orderAssets, oaOfThisOrder...)
		assetsJsonStr, _ := json.Marshal(assetsJsonB)

		startedAt := time.Now().UTC().Unix()
		endAt := int64(0)
		// have matched order or not
		haveMatchedOrder := utils.RandMatchedOrder()
		if haveMatchedOrder {
			endAt = utils.RandEndAt(startedAt)
		}

		orderInfo := &models.OrdersMd{
			Id:                  orderId,
			Maker:               utils.RandAddress(constants.ADDRESS_LEN),
			Kind:                int8(rand.Intn(8)),
			Assets:              string(assetsJsonStr),
			ExpiredAt:           startedAt + constants.DURATION_60_DAYS_SECOND,
			TokenPayment:        utils.RandTokenPayment(),
			StartedAt:           startedAt,
			BasePrice:           utils.RandPrice(),
			EndedAt:             endAt,
			EndedPrice:          utils.RandPrice(),
			ExpectedState:       "",
			Nonce:               int64(0),
			MarketFeePercentage: constants.MARKET_FEE_PERCENTAGE,
			Signature:           utils.RandStringRunes(42),
			Hash:                utils.RandStringRunes(50),
		}
		orders[i] = orderInfo

		// if have order matched => create record on matched_orders
		if haveMatchedOrder {
			matchedOrders = append(matchedOrders, &models.MatchedOrdersMd{
				OrderId:     orderId,
				Price:       int64(orderInfo.EndedPrice),
				Matcher:     utils.RandAddress(constants.ADDRESS_LEN),
				BlockNumber: orderId + 100000,
				TxHash:      utils.RandStringRunes(50),
			})
		}
	}

	// insert orders
	err := s.ordersStore.InsertMany(orders)
	if err != nil {
		fmt.Printf("error while inserting orders: %s", err.Error())
		return err
	}

	///////////////
	errChan := make(chan error)
	wgDone := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)

	// insert orderAssets
	go func() {
		err := s.orderAssetsStore.InsertMany(orderAssets)
		if err != nil {
			fmt.Printf("error while inserting orderAssets: %s", err.Error())
			errChan <- err
		}

		wg.Done()
	}()

	// insert matchedOrders
	go func() {
		err := s.matchedOrdersStore.InsertMany(matchedOrders)
		if err != nil {
			fmt.Printf("error while inserting matchedOrders: %s", err.Error())
			errChan <- err
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(wgDone)
	}()

	select {
	case <-wgDone:
		return nil
	case haveErr := <-errChan:
		return haveErr
	}
}

func sendLargeRequest() {
	AsyncHTTP()
}

func sendAxieDetailRequest(axieId int64, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	url := "https://testnet-graphql.skymavis.one/graphql"
	// fmt.Println("URL:>", url)

	r := fmt.Sprintf(`{"operationName":"GetAxieDetail","variables":{"axieId":"%d"},"query":"query GetAxieDetail($axieId: ID!) {\n  axie(axieId: $axieId) {\n    ...AxieDetail\n    __typename\n  }\n}\n\nfragment AxieDetail on Axie {\n  id\n  image\n  class\n  chain\n  name\n  genes\n  newGenes\n  owner\n  birthDate\n  bodyShape\n  class\n  sireId\n  sireClass\n  matronId\n  matronClass\n  stage\n  title\n  breedCount\n  level\n  figure {\n    atlas\n    model\n    image\n    __typename\n  }\n  parts {\n    ...AxiePart\n    __typename\n  }\n  stats {\n    ...AxieStats\n    __typename\n  }\n  order {\n    ...OrderInfo\n    __typename\n  }\n  ownerProfile {\n    name\n    __typename\n  }\n  battleInfo {\n    ...AxieBattleInfo\n    __typename\n  }\n  children {\n    id\n    name\n    class\n    image\n    title\n    stage\n    __typename\n  }\n  potentialPoints {\n    beast\n    aquatic\n    plant\n    bug\n    bird\n    reptile\n    mech\n    dawn\n    dusk\n    __typename\n  }\n  __typename\n}\n\nfragment AxieBattleInfo on AxieBattleInfo {\n  banned\n  banUntil\n  level\n  __typename\n}\n\nfragment AxiePart on AxiePart {\n  id\n  name\n  class\n  type\n  specialGenes\n  stage\n  abilities {\n    ...AxieCardAbility\n    __typename\n  }\n  __typename\n}\n\nfragment AxieCardAbility on AxieCardAbility {\n  id\n  name\n  attack\n  defense\n  energy\n  description\n  backgroundUrl\n  effectIconUrl\n  __typename\n}\n\nfragment AxieStats on AxieStats {\n  hp\n  speed\n  skill\n  morale\n  __typename\n}\n\nfragment OrderInfo on Order {\n  id\n  maker\n  kind\n  assets {\n    ...AssetInfo\n    __typename\n  }\n  expiredAt\n  paymentToken\n  startedAt\n  basePrice\n  endedAt\n  endedPrice\n  expectedState\n  nonce\n  marketFeePercentage\n  signature\n  hash\n  duration\n  timeLeft\n  currentPrice\n  suggestedPrice\n  currentPriceUsd\n  __typename\n}\n\nfragment AssetInfo on Asset {\n  erc\n  address\n  id\n  quantity\n  __typename\n}\n"}`, axieId)

	var jsonStr = []byte(r)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authority", "testnet-graphql.skymavis.one")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error: ", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("err handle it", err)
	}
	fmt.Println("response Body:", string(b))
	ch <- string(b)
}

func sendTransferRequest(axieId int64, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	url := "https://testnet-graphql.skymavis.one/graphql"
	// url := "http://localhost:4201"
	// fmt.Println("URL:>", url)

	r := fmt.Sprintf(`{"operationName":"GetAxieTransferHistory","variables":{"axieId":"%d","from":0,"size":5},"query":"query GetAxieTransferHistory($axieId: ID!, $from: Int!, $size: Int!) {\n  axie(axieId: $axieId) {\n    id\n    transferHistory(from: $from, size: $size) {\n      ...TransferRecords\n      __typename\n    }\n    ethereumTransferHistory(from: $from, size: $size) {\n      ...TransferRecords\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment TransferRecords on TransferRecords {\n  total\n  results {\n    from\n    to\n    timestamp\n    txHash\n    withPrice\n    __typename\n  }\n  __typename\n}\n"}`, axieId)
	// var jsonStr = []byte(`{"operationName":"GetAxieTransferHistory","variables":{"axieId":"58702","from":0,"size":5},"query":"query GetAxieTransferHistory($axieId: ID!, $from: Int!, $size: Int!) {\n  axie(axieId: $axieId) {\n    id\n    transferHistory(from: $from, size: $size) {\n      ...TransferRecords\n      __typename\n    }\n    ethereumTransferHistory(from: $from, size: $size) {\n      ...TransferRecords\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment TransferRecords on TransferRecords {\n  total\n  results {\n    from\n    to\n    timestamp\n    txHash\n    withPrice\n    __typename\n  }\n  __typename\n}\n"}`)
	var jsonStr = []byte(r)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authority", "testnet-graphql.skymavis.one")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error: ", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("err handle it")
	}
	fmt.Println("response Body:", string(b))
	ch <- string(b)
}

func AsyncHTTP() ([]string, error) {
	ch := make(chan string)
	var responses []string
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go sendTransferRequest(int64(1136), ch, &wg)
		// go sendAxieDetailRequest(int64(886), ch, &wg)
	}

	time.Sleep(2 * time.Second)
	// close the channel in the background
	go func() {
		wg.Wait()
		close(ch)
	}()
	// read from channel as they come in until its closed
	for res := range ch {
		responses = append(responses, res)
	}

	return responses, nil
}
