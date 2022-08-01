package constants

var VER_CONTRACT_ADDR = []string{
	"0xa3371de234bd8791126ac4fa8b88813dfe4f86e6", // rune-charm
	"0xcaca1c072d26e46686d932686015207fbe08fdb8", // axie contract
	"0x70bd60f625f6dd082ae1f59b80dc78cfa8b47f18", // land contract
	"0x8068a2c7735060589ab03685e220b322b5ec9a71",
}

var TOKEN_PAYMENT = []string{
	"0x29c6f8349a028e1bdfc68bfa08bdee7bc5d47e16",  // Ronin WETH Contract
	"0x3c4e17b9056272ce1b49f6900d8cfd6171a1 869d", // ASX contract
	"0x82f5483623d636bc3deba8ae67e1751b6cf2bad2",  // SLP contract
	"0x04ef1d4f687bb20eedcf05c7f710c078ba39f328",  // USD Coin Contract
}

// total orders record => BATCH_SIZE * BATCH_NUM
const BATCH_SIZE = 10
const BATCH_NUM = 3

// len address of matcher
const ADDRESS_LEN = 40

const DURATION_4_DAYS_SECOND = 345600
const DURATION_60_DAYS_SECOND int64 = 5184000

// max asset in a order ( just my rule =))) )
const MAX_NUM_ASSETS = 4

////////////////////////
// erc type
// 0: Erc20
// 1: Erc721
// 2: Erc1155
var ERC_TYPE = []int{
	0, 1, 2,
}

const MAX_QUANTITY = 100

const MARKET_FEE_PERCENTAGE = 425
