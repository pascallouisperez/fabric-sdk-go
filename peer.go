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

package fabricsdk

import (
	"encoding/pem"
	"time"

	pb "github.com/hyperledger/fabric/protos/peer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	config "github.com/hyperledger/fabric-sdk-go/config"
)

// Peer ...
/**
 * The Peer class represents a peer in the target blockchain network to which
 * HFC sends endorsement proposals, transaction ordering or query requests.
 *
 * The Peer class represents the remote Peer node and its network membership materials,
 * aka the ECert used to verify signatures. Peer membership represents organizations,
 * unlike User membership which represents individuals.
 *
 * When constructed, a Peer instance can be designated as an event source, in which case
 * a “eventSourceUrl” attribute should be configured. This allows the SDK to automatically
 * attach transaction event listeners to the event stream.
 *
 * It should be noted that Peer event streams function at the Peer level and not at the
 * chain and chaincode levels.
 */
type Peer struct {
	url                   string
	grpcDialOption        []grpc.DialOption
	name                  string
	roles                 []string
	enrollmentCertificate *pem.Block
}

// CreateNewPeer ...
/**
 * Constructs a Peer given its endpoint configuration settings.
 *
 * @param {string} url The URL with format of "host:port".
 */
func CreateNewPeer(url string) *Peer {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTimeout(time.Second*3))
	if config.IsTLSEnabled() {
		creds := credentials.NewClientTLSFromCert(config.GetTLSCACertPool(), config.GetTLSServerHostOverride())
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return &Peer{url: url, grpcDialOption: opts, name: "", roles: nil}
}

// ConnectEventSource ...
/**
 * Since practically all Peers are event producers, when constructing a Peer instance,
 * an application can designate it as the event source for the application. Typically
 * only one of the Peers on a Chain needs to be the event source, because all Peers on
 * the Chain produce the same events. This method tells the SDK which Peer(s) to use as
 * the event source for the client application. It is the responsibility of the SDK to
 * manage the connection lifecycle to the Peer’s EventHub. It is the responsibility of
 * the Client Application to understand and inform the selected Peer as to which event
 * types it wants to receive and the call back functions to use.
 * @returns {Future} This gives the app a handle to attach “success” and “error” listeners
 */
func (p *Peer) ConnectEventSource() {
	//to do
}

// IsEventListened ...
/**
 * A network call that discovers if at least one listener has been connected to the target
 * Peer for a given event. This helps application instance to decide whether it needs to
 * connect to the event source in a crash recovery or multiple instance deployment.
 * @param {string} eventName required
 * @param {Chain} chain optional
 * @result {bool} Whether the said event has been listened on by some application instance on that chain.
 */
func (p *Peer) IsEventListened(event string, chain *Chain) (bool, error) {
	//to do
	return false, nil
}

// AddListener ...
/**
 * For a Peer that is connected to eventSource, the addListener registers an EventCallBack for a
 * set of event types. addListener can be invoked multiple times to support differing EventCallBack
 * functions receiving different types of events.
 *
 * Note that the parameters below are optional in certain languages, like Java, that constructs an
 * instance of a listener interface, and pass in that instance as the parameter.
 * @param {string} eventType : ie. Block, Chaincode, Transaction
 * @param  {object} eventTypeData : Object Specific for event type as necessary, currently needed
 * for “Chaincode” event type, specifying a matching pattern to the event name set in the chaincode(s)
 * being executed on the target Peer, and for “Transaction” event type, specifying the transaction ID
 * @param {struct} eventCallback Client Application class registering for the callback.
 * @returns {string} An ID reference to the event listener.
 */
func (p *Peer) AddListener(eventType string, eventTypeData interface{}, eventCallback interface{}) (string, error) {
	//to do
	return "", nil
}

// RemoveListener ...
/**
 * Unregisters a listener.
 * @param {string} eventListenerRef Reference returned by SDK for event listener.
 * @return {bool} Success / Failure status
 */
func (p *Peer) RemoveListener(eventListenerRef string) (bool, error) {
	return false, nil
	//to do
}

// GetName ...
/**
 * Get the Peer name. Required property for the instance objects.
 * @returns {string} The name of the Peer
 */
func (p *Peer) GetName() string {
	return p.name
}

// SetName ...
/**
 * Set the Peer name / id.
 * @param {string} name
 */
func (p *Peer) SetName(name string) {
	p.name = name
}

// GetRoles ...
/**
 * Get the user’s roles the Peer participates in. It’s an array of possible values
 * in “client”, and “auditor”. The member service defines two more roles reserved
 * for peer membership: “peer” and “validator”, which are not exposed to the applications.
 * @returns {[]string} The roles for this user.
 */
func (p *Peer) GetRoles() []string {
	return p.roles
}

// SetRoles ...
/**
 * Set the user’s roles the Peer participates in. See getRoles() for legitimate values.
 * @param {[]string} roles The list of roles for the user.
 */
func (p *Peer) SetRoles(roles []string) {
	p.roles = roles
}

// GetEnrollmentCertificate ...
/**
 * Returns the Peer's enrollment certificate.
 * @returns {pem.Block} Certificate in PEM format signed by the trusted CA
 */
func (p *Peer) GetEnrollmentCertificate() *pem.Block {
	return p.enrollmentCertificate
}

// SetEnrollmentCertificate ...
/**
 * Set the Peer’s enrollment certificate.
 * @param {pem.Block} enrollment Certificate in PEM format signed by the trusted CA
 */
func (p *Peer) SetEnrollmentCertificate(pem *pem.Block) {
	p.enrollmentCertificate = pem
}

// GetURL ...
/**
 * Get the Peer url. Required property for the instance objects.
 * @returns {string} The address of the Peer
 */
func (p *Peer) GetURL() string {
	return p.url
}

// SendProposal ...
/**
 * Send  the created proposal to peer for endorsement.
 */
func (p *Peer) SendProposal(signedProposal *pb.SignedProposal) (*pb.ProposalResponse, error) {
	conn, err := grpc.Dial(p.url, p.grpcDialOption...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	endorserClient := pb.NewEndorserClient(conn)
	proposalResponse, err := endorserClient.ProcessProposal(context.Background(), signedProposal)
	if err != nil {
		return nil, err
	}
	return proposalResponse, nil
}
