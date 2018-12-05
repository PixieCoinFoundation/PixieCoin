package eth

import (
	"appcfg"
	"constants"
	"encoding/json"
	"ethereum"
	"fmt"
	. "language"
	. "logger"
	"math/big"

	deth "db/eth"
	"db/global"
	"errors"
	. "pixie_contract/api_specification"
	"service/mails"
	"shutdown"
	"time"
)

const (
	BALANCE_FUNCTION_SHA3  = "70a08231"
	TRANSFER_FUNCTION_SHA3 = "a9059cbb"
)

var TransactionReceiptNotFoundErr = errors.New("transaction receipt not found")

type TransactionReceipt struct {
	Status string `json:"status"`
}

var contractAddress string
var client *ethereum.Client

func init() {
	if appcfg.GetServerType() == "" && appcfg.SupportEthereum() {
		test1 := GenWeiBigIntFromETHFloat(0.01)
		test2 := GenWeiBigIntFromETHFloat(1)
		test3 := GenWeiBigIntFromETHFloat(100)
		if test1.Text(16) != "2386f26fc10000" {
			panic("not match 0.01 eth")
		}

		if test2.Text(16) != "de0b6b3a7640000" {
			panic("not match 1 eth")
		}

		if test3.Text(16) != "56bc75e2d63100000" {
			panic("not match 100 eth")
		}

		ethHost := appcfg.GetString("eth_host", "127.0.0.1")
		ethPort := appcfg.GetInt("eth_port", 8545)
		contractAddress = appcfg.GetString("eth_contract_address", "")

		if ethHost == "" || ethPort <= 0 || contractAddress == "" {
			Info("params empty", ethHost, ethPort, contractAddress)
			panic("ethereum init parameter empty!")
		} else {
			client = ethereum.NewClient(ethHost, ethPort)
		}

		if appcfg.GetBool("init_test_ethereum", false) {
			testAccount := appcfg.GetString("init_test_ethereum_account", "")

			if testAccount != "" {
				if eb, e := GetETHBalance(testAccount); e != nil {
					panic(e)
				} else {
					Info("account", testAccount, "eth balance", eb)
				}

				if pb, e := GetContractTokenBalance(testAccount); e != nil {
					panic(e)
				} else {
					Info("account", testAccount, "pxc balance", pb)
				}
			}
		}

		if appcfg.GetBool("main_server", false) {
			go checkTransactionLoop()
		}
	}
}

func checkTransactionLoop() {
	for {
		checkEthereumTransactionLoop1()

		time.Sleep(30 * time.Second)
	}
}

func checkEthereumTransactionLoop1() {
	shutdown.AddOneRequest("checkBuyTransactionLoop1")
	defer shutdown.DoneOneRequest("checkBuyTransactionLoop1")

	if list, err := deth.ListTransaction(constants.TRANSACTION_STATUS_SUBMIT, 1, 100); err == nil {
		for _, l := range list {
			if success, fail, receipt, data, err := CheckTransactionReceipt(l.Hash); err == nil {
				if receipt == "" {
					Err("cannot find receipt for transaction.please check", l.Hash)
					continue
				}

				if receipt == "null" {
					Info("transaction not finish:", l.Hash)
					continue
				}

				jobToken := "pxc_buy_" + l.Hash

				if success {
					deth.UpdateTransaction(l.Hash, constants.TRANSACTION_STATUS_SUCCESS, receipt, data)

					if l.AmountType == constants.AMOUNT_PXC_TYPE && l.ClothesList != nil && len(l.ClothesList) > 0 {
						if global.DoOnceJob(jobToken) {
							clos := make([]ClothesInfo, 0)
							for _, c := range l.ClothesList {
								clos = append(clos, ClothesInfo{ClothesID: c.ClothesID, Count: c.ClothesCount})
							}

							content := fmt.Sprintf(L("pxc3"), l.ToNickname)
							mails.SendToOneC("", l.FromUsername, L("pxc1"), content, 0, 0, clos, true)
						}
					}
				} else {
					if fail {
						if global.DoOnceJob(jobToken) {
							content := fmt.Sprintf(L("pxc5"), l.ToNickname)

							if l.ClothesList != nil && len(l.ClothesList) > 0 {
								mails.SendPXCBuyFailMail(l.FromUsername, L("pxc4"), content, l.ClothesList)
							}
						}

						deth.UpdateTransaction(l.Hash, constants.TRANSACTION_STATUS_FAIL, receipt, data)
					} else {
						ErrMail("unknown eth transaction", data, receipt)
						// deth.UpdateTransaction(l.Hash, constants.TRANSACTION_STATUS_UNKNOWN, receipt, data)
					}
				}
			}
		}
	}
}

func GenWeiBigIntFromETHFloat(v float64) *big.Int {
	in := v
	div := 1
	for {
		f := big.NewFloat(in)

		if f.IsInt() {
			break
		} else {
			in = in * 10
			div = div * 10
		}
	}

	inInt := big.NewInt(int64(in))
	res := big.NewInt(0)

	res = res.Mul(inInt, big.NewInt(int64(ethereum.OneEtherInWei/div)))

	return res
}

func GenETHAccount(cli *ethereum.Client, pwd string) (account string, unlockSuccess bool, err error) {
	account, unlockSuccess, err = newAccountWithUnlock(cli, pwd)

	return
}

func newAccountWithUnlock(cli *ethereum.Client, pwd string) (na string, unlockSuccess bool, err error) {
	c := cli
	if c == nil {
		c = client
	}

	var resp json.RawMessage
	if resp, err = c.Call("personal_newAccount", pwd); err != nil {
		Err(err)
		return
	} else {
		json.Unmarshal(resp, &na)

		//unlock account
		unlockSuccess, err = UnlockAccount(c, na, pwd)
	}

	return
}

func UnlockAccount(cli *ethereum.Client, account string, pwd string) (success bool, err error) {
	c := cli
	if c == nil {
		c = client
	}

	var resp json.RawMessage
	if resp, err = c.Call("personal_unlockAccount", account, pwd); err != nil {
		Err(err)
		return
	} else {
		success = string(resp) == "true"
	}
	return
}

func GetETHBalance(account string) (ethBalance string, err error) {
	var bw *ethereum.Wei
	if bw, err = client.GetBalance(account); err != nil {
		Err(err)
		return
	}

	ethBalance = bw.EtherFloatText()
	return
}

func GetContractTokenBalance(account string) (pxcBalance string, err error) {
	bp := &ethereum.BalanceParam{
		To:   contractAddress,
		Data: getPXCBalanceData(account),
	}

	//get pxc balance
	var resp json.RawMessage
	if resp, err = client.Call("eth_call", bp, "latest"); err != nil {
		Err(err)
		return
	} else {
		var res string
		if err = json.Unmarshal(resp, &res); err != nil {
			Err(err)
			return
		}

		if res == "0x" {
			Err("contractAddress may not legal:", contractAddress)
			err = constants.GetContractBalanceErr
			return
		}

		var pb big.Int
		if err = ethereum.DecodeHex(res, &pb); err != nil {
			Err(err)
			return
		}

		pxcBalance = pb.Text(10)
	}

	return
}

func GetGasPrice(cli *ethereum.Client) (gp string, err error) {
	c := cli
	if c == nil {
		c = client
	}

	var priceInWei *ethereum.Wei
	if priceInWei, err = c.GetGasPrice(); err != nil {
		Err(err)
		return
	}

	gp = priceInWei.Int.Text(10)

	return
}

func EstimateGasForTransferETH(cli *ethereum.Client, from, to string) (gas string, err error) {
	c := cli
	if c == nil {
		c = client
	}

	var rw *ethereum.Wei
	if rw, err = c.EstimateGas(from, to); err != nil {
		Err(err)
		return
	}

	gas = rw.Int.Text(10)

	return
}

func EstimateGasForTransferContractToken(cli *ethereum.Client, ca, from, to string, amountHex string) (gas string, err error) {
	c := cli
	if c == nil {
		c = client
	}

	contractAddr := ca
	if contractAddr == "" {
		contractAddr = contractAddress
	}

	var rw *ethereum.Wei
	if rw, err = c.EstimateContractGas(from, contractAddr, getTransferPXCData(to, amountHex)); err != nil {
		Err(err)
		return
	}

	gas = rw.Int.Text(10)

	return
}

func TransferETHInWei(from, to, pwd string, amount *big.Int) (transactionHash string, err error) {
	var transaction *ethereum.Transaction
	if transaction, err = client.PersonalSendTransaction(from, to, pwd, &ethereum.Wei{*amount}); err != nil {
		if err.Error() == constants.CURRENCY_NOT_ENOUGH_ERR_MSG {
			Info("not enough wei", from, to, amount)
		} else {
			Err(err)
		}

		return
	} else {
		transactionHash = transaction.ID
	}

	return
}

func TransferContractToken(from, to, pwd string, amountHex string) (transactionHash string, err error) {
	var transaction *ethereum.TokenTransaction
	if transaction, err = client.PersonalSendTokenTransaction(from, contractAddress, pwd, getTransferPXCData(to, amountHex)); err != nil {
		if err.Error() == constants.CURRENCY_NOT_ENOUGH_ERR_MSG {
			Info("not enough contract token", from, to, amountHex, contractAddress)
		} else {
			Err(err)
		}

		return
	} else {
		transactionHash = transaction.ID
	}

	return
}

func CheckTransactionReceipt(transactionHash string) (success bool, fail bool, receipt string, data string, err error) {
	var resp json.RawMessage

	//get data
	if resp, err = client.Call("eth_getTransactionByHash", transactionHash); err != nil {
		Err(err)
		return
	} else {
		data = string(resp)
	}

	//get receipt
	if resp, err = client.Call("eth_getTransactionReceipt", transactionHash); err != nil {
		if err.Error() == constants.UNKNOWN_TRANSACTION_ERR_MSG {
			Info("transaction not found maybe pending...", transactionHash)
			err = TransactionReceiptNotFoundErr
			return
		} else {
			Err(err)
			return
		}
	} else {
		receipt = string(resp)
		if receipt == "null" || receipt == "" {
			return
		}

		var rp TransactionReceipt
		if err = json.Unmarshal(resp, &rp); err != nil {
			Err(err)
			return
		}

		success = rp.Status == "0x1"
		fail = rp.Status == "0x0"
	}
	return
}

//support functions
func getPXCBalanceData(addr string) string {
	return fmt.Sprintf("0x%s000000000000000000000000%s", BALANCE_FUNCTION_SHA3, addr[2:])
}

func getTransferPXCData(dest string, amountHex string) string {
	if len(dest) <= 2 || len(amountHex) <= 2 {
		return ""
	}

	amountHexTail := amountHex[2:]
	for len(amountHexTail) < 64 {
		amountHexTail = "0" + amountHexTail
	}

	return fmt.Sprintf("0x%s000000000000000000000000%s%s", TRANSFER_FUNCTION_SHA3, dest[2:], amountHexTail)
}
