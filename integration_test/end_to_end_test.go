/*
Copyright SecureKey Technologies Inc. All Rights Reserved.


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at


      http://www.apache.org/licenses/LICENSE-2.0


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration_test

import (
	"encoding/pem"
	"fmt"
	"strconv"
	"testing"
	"time"

	fabric_sdk "github.com/hyperledger/fabric-sdk-go"
	events "github.com/hyperledger/fabric-sdk-go/events"

	config "github.com/hyperledger/fabric-sdk-go/config"
	kvs "github.com/hyperledger/fabric-sdk-go/keyvaluestore"
	msp "github.com/hyperledger/fabric-sdk-go/msp"
	"github.com/hyperledger/fabric/bccsp"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/common/util"

	"github.com/hyperledger/fabric/bccsp/sw"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var chainCodeId = "end2end"
var chainId = "testchainid"

func TestChainCodeInvoke(t *testing.T) {
	InitConfigForEndToEnd()

	eventHub := GetEventHub(t, nil)
	querychain, invokechain := GetChains(t)

	// Get Query value before invoke
	value, err := getQueryValue(t, querychain)
	if err != nil {
		t.Fatalf("getQueryValue return error: %v", err)
	}
	fmt.Printf("*** QueryValue before invoke %s\n", value)

	err = invoke(t, invokechain, eventHub)
	if err != nil {
		t.Fatalf("invoke return error: %v", err)
	}

	valueAfterInvoke, err := getQueryValue(t, querychain)
	if err != nil {
		t.Errorf("getQueryValue return error: %v", err)
		return
	}
	fmt.Printf("*** QueryValue after invoke %s\n", valueAfterInvoke)

	valueInt, _ := strconv.Atoi(value)
	valueInt = valueInt + 1
	valueAfterInvokeInt, _ := strconv.Atoi(valueAfterInvoke)
	if valueInt != valueAfterInvokeInt {
		t.Fatalf("SendTransaction didn't change the QueryValue")

	}

}

func getQueryValue(t *testing.T, chain *fabric_sdk.Chain) (string, error) {

	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, "b")

	txid := util.GenerateUUID()

	signedProposal, _, err := chain.CreateTransactionProposal(chainCodeId, chainId, args, true, txid, nil)
	if err != nil {
		return "", fmt.Errorf("SendTransactionProposal return error: %v", err)
	}
	transactionProposalResponse, err := chain.SendTransactionProposal(signedProposal, 0)
	if err != nil {
		return "", fmt.Errorf("SendTransactionProposal return error: %v", err)
	}

	for _, v := range transactionProposalResponse {
		if v.Err != nil {
			return "", fmt.Errorf("Endorser %s return error: %v", v.Endorser, v.Err)
		}
		return string(v.ProposalResponse.GetResponse().Payload), nil
	}
	return "", nil
}

func invoke(t *testing.T, chain *fabric_sdk.Chain, eventHub *events.EventHub) error {

	var args []string
	args = append(args, "invoke")
	args = append(args, "move")
	args = append(args, "a")
	args = append(args, "b")
	args = append(args, "1")

	txId := util.GenerateUUID()

	signedProposal, proposal, err := chain.CreateTransactionProposal(chainCodeId, chainId, args, true, txId, nil)
	if err != nil {
		return fmt.Errorf("SendTransactionProposal return error: %v", err)
	}
	transactionProposalResponse, err := chain.SendTransactionProposal(signedProposal, 0)
	if err != nil {
		return fmt.Errorf("SendTransactionProposal return error: %v", err)
	}

	var proposalResponses []*pb.ProposalResponse
	for _, v := range transactionProposalResponse {
		if v.Err != nil {
			return fmt.Errorf("Endorser %s return error: %v", v.Endorser, v.Err)
		}
		proposalResponses = append(proposalResponses, v.ProposalResponse)
		fmt.Printf("Endorser '%s' return ProposalResponse:%v\n", v.Endorser, v.ProposalResponse.GetResponse())
	}

	tx, err := chain.CreateTransaction(proposal, proposalResponses)
	if err != nil {
		return fmt.Errorf("CreateTransaction return error: %v", err)

	}
	transactionResponse, err := chain.SendTransaction(proposal, tx)
	if err != nil {
		return fmt.Errorf("SendTransaction return error: %v", err)

	}
	for _, v := range transactionResponse {
		if v.Err != nil {
			return fmt.Errorf("Orderer %s return error: %v", v.Orderer, v.Err)
		}
	}
	done := make(chan bool)
	eventHub.RegisterTxEvent(txId, func(txId string, err error) {
		fmt.Printf("receive success event for txid(%s)\n", txId)
		done <- true
	})

	select {
	case <-done:
	case <-time.After(time.Second * 20):
		return fmt.Errorf("Didn't receive block event for txid(%s)\n", txId)
	}
	return nil

}

func InitConfigForEndToEnd() {
	err := config.InitConfig("./test_resources/config/config_test.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetChains(t *testing.T) (*fabric_sdk.Chain, *fabric_sdk.Chain) {

	client := fabric_sdk.NewClient()
	ks := &sw.FileBasedKeyStore{}
	if err := ks.Init(nil, config.GetKeyStorePath(), false); err != nil {
		t.Fatalf("Failed initializing key store [%s]", err)
	}

	cryptoSuite, err := bccspFactory.GetBCCSP(&bccspFactory.SwOpts{Ephemeral_: true, SecLevel: config.GetSecurityLevel(),
		HashFamily: config.GetSecurityAlgorithm(), KeyStore: ks})
	if err != nil {
		t.Fatalf("Failed getting ephemeral software-based BCCSP [%s]", err)
	}
	client.SetCryptoSuite(cryptoSuite)
	stateStore, err := kvs.CreateNewFileKeyValueStore("/tmp/enroll_user")
	if err != nil {
		t.Fatalf("CreateNewFileKeyValueStore return error[%s]", err)
	}
	client.SetStateStore(stateStore)
	user, err := client.GetUserContext("testUser")
	if err != nil {
		t.Fatalf("client.GetUserContext return error: %v", err)
	}
	if user == nil {
		msps, err1 := msp.NewMSPServices(config.GetMspClientPath())
		if err1 != nil {
			t.Fatalf("NewFabricCOPServices return error: %v", err1)
		}
		key, cert, err1 := msps.Enroll("testUser", "user1")
		keyPem, _ := pem.Decode(key)
		if err1 != nil {
			t.Fatalf("Enroll return error: %v", err1)
		}
		user := fabric_sdk.NewUser("testUser")
		k, err1 := client.GetCryptoSuite().KeyImport(keyPem.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: false})
		if err1 != nil {
			t.Fatalf("KeyImport return error: %v", err1)
		}
		user.SetPrivateKey(k)
		user.SetEnrollmentCertificate(cert)
		err = client.SetUserContext(user, false)
		if err != nil {
			t.Fatalf("client.SetUserContext return error: %v", err)
		}
	}

	querychain, err := client.NewChain("querychain")
	if err != nil {
		t.Fatalf("NewChain return error: %v", err)
	}

	for _, p := range config.GetPeersConfig() {
		endorser := fabric_sdk.CreateNewPeer(fmt.Sprintf("%s:%s", p.Host, p.Port))
		querychain.AddPeer(endorser)
		break
	}

	invokechain, err := client.NewChain("invokechain")
	if err != nil {
		t.Fatalf("NewChain return error: %v", err)
	}
	orderer := fabric_sdk.CreateNewOrderer(fmt.Sprintf("%s:%s", config.GetOrdererHost(), config.GetOrdererPort()))
	invokechain.AddOrderer(orderer)

	for _, p := range config.GetPeersConfig() {
		endorser := fabric_sdk.CreateNewPeer(fmt.Sprintf("%s:%s", p.Host, p.Port))
		invokechain.AddPeer(endorser)
	}

	return querychain, invokechain

}

func GetEventHub(t *testing.T, interestedEvents []*pb.Interest) *events.EventHub {
	eventHub := events.NewEventHub()
	foundEventHub := false
	for _, p := range config.GetPeersConfig() {
		if p.EventHost != "" && p.EventPort != "" {
			eventHub.SetPeerAddr(fmt.Sprintf("%s:%s", p.EventHost, p.EventPort))
			foundEventHub = true
			break
		}
	}

	if !foundEventHub {
		t.Fatalf("No EventHub configuration found")
	}

	if interestedEvents != nil {
		eventHub.SetInterestedEvents(interestedEvents)
	}

	if err := eventHub.Connect(); err != nil {
		t.Fatalf("Failed eventHub.Connect() [%s]", err)
	}

	return eventHub
}
