package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"mp-generate-test-data/config"
	"mp-generate-test-data/internal/constants"
	"mp-generate-test-data/internal/models"
	"mp-generate-test-data/internal/stores"
	"mp-generate-test-data/internal/utils"
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

// func (s *Service) enableCors(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Headers", "*")
// 	(*w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(*w).Header().Set("Access-Control-Allow-Methods", "*")
// }

func (s *Service) Run() {
	s.gen_test_data()
}

func (s *Service) gen_test_data() {
	// get the largest id in table orders
	order, exist, err := s.ordersStore.GetLargestOrderId()
	if err != nil {
		fmt.Println("error while getting largest order id: " + err.Error())
		return
	}

	fromOrderId := int64(0)
	if exist {
		fromOrderId = order.Id + 1
	}

	for i := 0; i < constants.BATCH_NUM; i++ {
		err := s.gen_orders(fromOrderId + int64(constants.BATCH_SIZE*i))
		if err != nil {
			return
		} else {
			fmt.Printf("Done order id: %d\n", fromOrderId+int64(constants.BATCH_SIZE*i))
		}
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
			assetsJsonB = append(assetsJsonB, &models.Assets{
				Address:  utils.RandVerifiedContract(),
				Erc:      int8(utils.RandErcType()),
				Id:       rand.Int63(),
				Quantity: utils.RandQuantity(),
			})
		}

		for k := 0; k < numOfAssets; k++ {
			oaOfThisOrder = append(oaOfThisOrder, &models.OrderAssetsMd{
				OrderId:  orderId,
				Address:  assetsJsonB[k].Address,
				Erc:      assetsJsonB[k].Erc,
				Id:       assetsJsonB[k].Id,
				Quantity: int64(assetsJsonB[k].Quantity),
			})
		}
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

	// insert orderAssets
	err = s.orderAssetsStore.InsertMany(orderAssets)
	if err != nil {
		fmt.Printf("error while inserting orderAssets: %s", err.Error())
		return err
	}

	// insert matchedOrders
	err = s.matchedOrdersStore.InsertMany(matchedOrders)
	if err != nil {
		fmt.Printf("error while inserting matchedOrders: %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) gen_order_assets()
