package main

import (
	"fmt"
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

const (
	mainAccountPublicKey string = "GC4W3QO7ZJCX3OYVH4EWGNIUYPA4BA6SQUUQKUJFVGMH5KR6K3AIQQKZ"
	mainAccountSecretKey string = "SDJFZK4VGVZ45XURGEAD3FW73S2CRTRYFFPQBXMO6X33EPTGBZDSZA4B"
	networkPassphrase    string = network.TestNetworkPassphrase
)

var (
	networkClient *horizonclient.Client = horizonclient.DefaultTestNetClient
)

func getMainAccountDetails() (horizon.Account, error) {

	kp, _ := keypair.ParseFull(mainAccountSecretKey)

	req := horizonclient.AccountRequest{
		AccountID: kp.Address(),
	}

	return networkClient.AccountDetail(req)
}

func createSignedTransaction(mainAccount horizon.Account, mainAccountKeypair *keypair.Full) (*txnbuild.Transaction, error) {

	manageDataOperation := txnbuild.ManageData{
		Name:  "Hello",
		Value: []byte("World"),
	}

	transaction, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &mainAccount,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&manageDataOperation},
		BaseFee:              txnbuild.MinBaseFee,
		Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()}, // Use a real timeout in production!
	})

	if err != nil {
		return nil, err
	}

	return transaction.Sign(networkPassphrase, mainAccountKeypair)
}

func main() {

	accountDetails, err := getMainAccountDetails()

	if err != nil {
		log.Println("error getting account details", err)
		return
	}

	accountKeypair := keypair.MustParseFull(mainAccountSecretKey)

	transaction, err := createSignedTransaction(accountDetails, accountKeypair)

	if err != nil {
		log.Println("error creating signed transaction", err)
		return
	}

	_, err = networkClient.SubmitTransaction(transaction)

	if err != nil {
		log.Println("error submitting the manage data transaction", err)
		return
	}

	fmt.Println("Successfully added metadata to an account")
}
