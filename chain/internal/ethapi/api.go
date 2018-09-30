package ethapi

import (
	"context"
	_ "crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/boker/chain/accounts"
	"github.com/boker/chain/accounts/keystore"
	"github.com/boker/chain/boker/api"
	"github.com/boker/chain/boker/protocol"
	"github.com/boker/chain/common"
	"github.com/boker/chain/common/hexutil"
	"github.com/boker/chain/common/math"
	"github.com/boker/chain/core"
	"github.com/boker/chain/core/types"
	"github.com/boker/chain/core/vm"
	"github.com/boker/chain/crypto"
	"github.com/boker/chain/log"
	"github.com/boker/chain/p2p"
	"github.com/boker/chain/params"
	"github.com/boker/chain/rlp"
	"github.com/boker/chain/rpc"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	defaultGas      = 90000
	defaultGasPrice = 50 * params.Shannon
)

//提供访问以太坊相关信息的API。它仅提供对公共数据进行操作的方法，任何人都可以免费使用
type PublicEthereumAPI struct {
	b Backend
}

func NewPublicEthereumAPI(b Backend) *PublicEthereumAPI {
	return &PublicEthereumAPI{
		b: b,
	}
}

//返回Gas的建议价格
func (s *PublicEthereumAPI) GasPrice(ctx context.Context) (*big.Int, error) {
	return s.b.SuggestPrice(ctx)
}

//返回此节点支持的当前以太坊协议版本
func (s *PublicEthereumAPI) ProtocolVersion() hexutil.Uint {
	return hexutil.Uint(s.b.ProtocolVersion())
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronise from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (s *PublicEthereumAPI) Syncing() (interface{}, error) {
	progress := s.b.Downloader().Progress()

	// Return not syncing if the synchronisation already completed
	if progress.CurrentBlock >= progress.HighestBlock {
		return false, nil
	}
	// Otherwise gather the block sync stats
	return map[string]interface{}{
		"startingBlock": hexutil.Uint64(progress.StartingBlock),
		"currentBlock":  hexutil.Uint64(progress.CurrentBlock),
		"highestBlock":  hexutil.Uint64(progress.HighestBlock),
		"pulledStates":  hexutil.Uint64(progress.PulledStates),
		"knownStates":   hexutil.Uint64(progress.KnownStates),
	}, nil
}

//为交易池提供API接口， 它仅对非机密数据进行操作。
type PublicTxPoolAPI struct {
	b Backend
}

func NewPublicTxPoolAPI(b Backend) *PublicTxPoolAPI {
	return &PublicTxPoolAPI{
		b: b,
	}
}

//内容返回交易池中包含的交易
func (s *PublicTxPoolAPI) Content() map[string]map[string]map[string]*RPCTransaction {
	content := map[string]map[string]map[string]*RPCTransaction{
		"pending": make(map[string]map[string]*RPCTransaction),
		"queued":  make(map[string]map[string]*RPCTransaction),
	}
	pending, queue := s.b.TxPoolContent()

	// Flatten the pending transactions
	for account, txs := range pending {
		dump := make(map[string]*RPCTransaction)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = newRPCPendingTransaction(tx)
		}
		content["pending"][account.Hex()] = dump
	}
	// Flatten the queued transactions
	for account, txs := range queue {
		dump := make(map[string]*RPCTransaction)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = newRPCPendingTransaction(tx)
		}
		content["queued"][account.Hex()] = dump
	}
	return content
}

//返回交易池中挂起和排队的交易数量。
func (s *PublicTxPoolAPI) Status() map[string]hexutil.Uint {
	pending, queue := s.b.Stats()
	return map[string]hexutil.Uint{
		"pending": hexutil.Uint(pending),
		"queued":  hexutil.Uint(queue),
	}
}

//检索交易池的内容并将其展平为一个易于检查的清单
func (s *PublicTxPoolAPI) Inspect() map[string]map[string]map[string]string {
	content := map[string]map[string]map[string]string{
		"pending": make(map[string]map[string]string),
		"queued":  make(map[string]map[string]string),
	}
	pending, queue := s.b.TxPoolContent()

	// Define a formatter to flatten a transaction into a string
	var format = func(tx *types.Transaction) string {
		if to := tx.To(); to != nil {
			return fmt.Sprintf("%s: %v wei + %v gas × %v wei", tx.To().Hex(), tx.Value(), tx.Gas(), tx.GasPrice())
		}
		return fmt.Sprintf("contract creation: %v wei + %v gas × %v wei", tx.Value(), tx.Gas(), tx.GasPrice())
	}
	// Flatten the pending transactions
	for account, txs := range pending {
		dump := make(map[string]string)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = format(tx)
		}
		content["pending"][account.Hex()] = dump
	}
	// Flatten the queued transactions
	for account, txs := range queue {
		dump := make(map[string]string)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = format(tx)
		}
		content["queued"][account.Hex()] = dump
	}
	return content
}

//提供了访问此节点管理的帐户的API接口，它仅提供可以检索帐户的方法
type PublicAccountAPI struct {
	am *accounts.Manager
}

func NewPublicAccountAPI(am *accounts.Manager) *PublicAccountAPI {
	return &PublicAccountAPI{am: am}
}

//返回此节点管理的帐户集合
func (s *PublicAccountAPI) Accounts() []common.Address {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty
	for _, wallet := range s.am.Wallets() {
		for _, account := range wallet.Accounts() {
			addresses = append(addresses, account.Address)
		}
	}
	return addresses
}

//提供访问此节点管理的帐户的API接口，它提供了创建（un）锁定列表帐户的方法。 有些方法接受密码，因此默认情况下被视为私有。
type PrivateAccountAPI struct {
	am        *accounts.Manager
	nonceLock *AddrLocker
	b         Backend
}

func NewPrivateAccountAPI(b Backend, nonceLock *AddrLocker) *PrivateAccountAPI {
	return &PrivateAccountAPI{
		am:        b.AccountManager(),
		nonceLock: nonceLock,
		b:         b,
	}
}

//返回此节点管理的帐户的地址列表
func (s *PrivateAccountAPI) ListAccounts() []common.Address {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty
	for _, wallet := range s.am.Wallets() {
		for _, account := range wallet.Accounts() {
			addresses = append(addresses, account.Address)
		}
	}
	return addresses
}

// rawWallet is a JSON representation of an accounts.Wallet interface, with its
// data contents extracted into plain fields.
type rawWallet struct {
	URL      string             `json:"url"`
	Status   string             `json:"status"`
	Failure  string             `json:"failure,omitempty"`
	Accounts []accounts.Account `json:"accounts,omitempty"`
}

// ListWallets will return a list of wallets this node manages.
func (s *PrivateAccountAPI) ListWallets() []rawWallet {
	wallets := make([]rawWallet, 0) // return [] instead of nil if empty
	for _, wallet := range s.am.Wallets() {
		status, failure := wallet.Status()

		raw := rawWallet{
			URL:      wallet.URL().String(),
			Status:   status,
			Accounts: wallet.Accounts(),
		}
		if failure != nil {
			raw.Failure = failure.Error()
		}
		wallets = append(wallets, raw)
	}
	return wallets
}

// OpenWallet initiates a hardware wallet opening procedure, establishing a USB
// connection and attempting to authenticate via the provided passphrase. Note,
// the method may return an extra challenge requiring a second open (e.g. the
// Trezor PIN matrix challenge).
func (s *PrivateAccountAPI) OpenWallet(url string, passphrase *string) error {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return err
	}
	pass := ""
	if passphrase != nil {
		pass = *passphrase
	}
	return wallet.Open(pass)
}

// DeriveAccount requests a HD wallet to derive a new account, optionally pinning
// it for later reuse.
func (s *PrivateAccountAPI) DeriveAccount(url string, path string, pin *bool) (accounts.Account, error) {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return accounts.Account{}, err
	}
	derivPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return accounts.Account{}, err
	}
	if pin == nil {
		pin = new(bool)
	}
	return wallet.Derive(derivPath, *pin)
}

// NewAccount will create a new account and returns the address for the new account.
func (s *PrivateAccountAPI) NewAccount(password string) (common.Address, error) {
	acc, err := fetchKeystore(s.am).NewAccount(password)
	if err == nil {
		return acc.Address, nil
	}
	return common.Address{}, err
}

// fetchKeystore retrives the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) *keystore.KeyStore {
	return am.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
}

// ImportRawKey stores the given hex encoded ECDSA key into the key directory,
// encrypting it with the passphrase.
func (s *PrivateAccountAPI) ImportRawKey(privkey string, password string) (common.Address, error) {
	key, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}
	acc, err := fetchKeystore(s.am).ImportECDSA(key, password)
	return acc.Address, err
}

// UnlockAccount will unlock the account associated with the given address with
// the given password for duration seconds. If duration is nil it will use a
// default of 300 seconds. It returns an indication if the account was unlocked.
func (s *PrivateAccountAPI) UnlockAccount(addr common.Address, password string, duration *uint64) (bool, error) {
	const max = uint64(time.Duration(math.MaxInt64) / time.Second)
	var d time.Duration
	if duration == nil {
		d = 300 * time.Second
	} else if *duration > max {
		return false, errors.New("unlock duration too large")
	} else {
		d = time.Duration(*duration) * time.Second
	}
	err := fetchKeystore(s.am).TimedUnlock(accounts.Account{Address: addr}, password, d)

	//在这里默认设置用户的Password信息放入到配置文件中
	s.b.SetPassword(password)
	log.Info("Set Coinbase Password", "Account", addr, "Password", password)

	//将解锁账号设置为Coinbase
	s.b.SetCoinbase(addr)

	return err == nil, err
}

// LockAccount will lock the account associated with the given address when it's unlocked.
func (s *PrivateAccountAPI) LockAccount(addr common.Address) bool {
	return fetchKeystore(s.am).Lock(addr) == nil
}

//将根据给定的参数创建一个交易，尝试使用与args.To关联的键对其进行签名。 如果给定的passwd不是能够解密失败的密钥。
func (s *PrivateAccountAPI) SendTransaction(ctx context.Context, args SendTxArgs, passwd string) (common.Hash, error) {

	log.Info("****PrivateAccountAPI SendTransaction****", "passwd", passwd)

	//查找包含所请求签名者的钱包
	account := accounts.Account{Address: args.From}

	//根据帐号得到钱包信息
	wallet, err := s.am.Find(account)
	if err != nil {
		return common.Hash{}, err
	}

	//
	if args.Nonce == nil {
		//保持帐号的互斥围绕签名从而可以防止并发分配以及多个帐户相同的随机数。
		s.nonceLock.LockAddr(args.From)
		defer s.nonceLock.UnlockAddr(args.From)
	}

	//设置一些默认值
	if err := args.SetDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}

	//由于这里只会是RPC调用到，因此这里只对普通交易进行封装，非普通交易则封装失败
	tx, transErr := args.ToTransaction()
	if transErr != nil {
		return common.Hash{}, transErr
	}

	//获取区块链的配置，检查是否是EIP155的区块号(https://github.com/ethereum/eips/issues/155)
	var chainID *big.Int
	if config := s.b.ChainConfig(); config.IsEIP155(s.b.CurrentBlock().Number()) {
		chainID = config.ChainId
	}

	//对该笔交易签名来确保该笔交易的真实有效性
	signed, err := wallet.SignTxWithPassphrase(account, passwd, tx, chainID)
	if err != nil {
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.b, signed)
}

// signHash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calulcated as
//   keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

// Sign calculates an Ethereum ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message))
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The key used to calculate the signature is decrypted with the given password.
//
// https://github.com/boker/chain/wiki/Management-APIs#personal_sign
func (s *PrivateAccountAPI) Sign(ctx context.Context, data hexutil.Bytes, addr common.Address, passwd string) (hexutil.Bytes, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Assemble sign the data with the wallet
	signature, err := wallet.SignHashWithPassphrase(account, passwd, signHash(data))
	if err != nil {
		return nil, err
	}
	signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

// EcRecover returns the address for the account that was used to create the signature.
// Note, this function is compatible with eth_sign and personal_sign. As such it recovers
// the address of:
// hash = keccak256("\x19Ethereum Signed Message:\n"${message length}${message})
// addr = ecrecover(hash, signature)
//
// Note, the signature must conform to the secp256k1 curve R, S and V values, where
// the V value must be be 27 or 28 for legacy reasons.
//
// https://github.com/boker/chain/wiki/Management-APIs#personal_ecRecover
func (s *PrivateAccountAPI) EcRecover(ctx context.Context, data, sig hexutil.Bytes) (common.Address, error) {
	if len(sig) != 65 {
		return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sig[64] != 27 && sig[64] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[64] -= 27 // Transform yellow paper V from 27/28 to 0/1

	rpk, err := crypto.Ecrecover(signHash(data), sig)
	if err != nil {
		return common.Address{}, err
	}
	pubKey := crypto.ToECDSAPub(rpk)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr, nil
}

// SignAndSendTransaction was renamed to SendTransaction. This method is deprecated
// and will be removed in the future. It primary goal is to give clients time to update.
func (s *PrivateAccountAPI) SignAndSendTransaction(ctx context.Context, args SendTxArgs, passwd string) (common.Hash, error) {

	log.Info("****PrivateAccountAPI SignAndSendTransaction****", "passwd", passwd)
	return s.SendTransaction(ctx, args, passwd)
}

//提供了一个API来访问以太坊区块链,它仅提供对公共数据进行操作的方法，任何人都可以免费使用。
type PublicBlockChainAPI struct {
	b Backend
}

func NewPublicBlockChainAPI(b Backend) *PublicBlockChainAPI {
	return &PublicBlockChainAPI{b}
}

//得到当前的区块序号
func (s *PublicBlockChainAPI) BlockNumber() *big.Int {
	header, _ := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber) // latest header should always be available
	return header.Number
}

//GetBalance返回给定地址的wei数量给定块号。 rpc.LatestBlockNumber和rpc.PendingBlockNumber元块号也是允许的。
func (s *PublicBlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (*big.Int, error) {

	state, _, err := s.b.StateAndHeaderByNumber(ctx, blockNr)
	if state == nil || err != nil {
		return nil, err
	}
	b := state.GetBalance(address)
	return b, state.Error()
}

//返回请求的块，当blockNr为-1时，返回链头。 当fullTx为真时全部完整详细地返回块中的交易，否则仅返回交易哈希。
func (s *PublicBlockChainAPI) GetBlockByNumber(ctx context.Context, blockNr rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {

	block, err := s.b.BlockByNumber(ctx, blockNr)
	if block != nil {
		response, err := s.rpcOutputBlock(block, true, fullTx)
		if err == nil && blockNr == rpc.PendingBlockNumber {
			// Pending blocks need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}
		return response, err
	}
	return nil, err
}

//返回请求的块，当fullTx为true时，块中的所有交易都将完整返回，否则只返回交易哈希
func (s *PublicBlockChainAPI) GetBlockByHash(ctx context.Context, blockHash common.Hash, fullTx bool) (map[string]interface{}, error) {
	block, err := s.b.GetBlock(ctx, blockHash)
	if block != nil {
		return s.rpcOutputBlock(block, true, fullTx)
	}
	return nil, err
}

//返回请求的块，当fullTx为true时，块中的所有交易都将完整返回，否则只返回交易哈希
func (s *PublicBlockChainAPI) GetUncleByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) (map[string]interface{}, error) {
	block, err := s.b.BlockByNumber(ctx, blockNr)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			log.Debug("Requested uncle not found", "number", blockNr, "hash", block.Hash(), "index", index)
			return nil, nil
		}
		block = types.NewBlockWithHeader(uncles[index])
		return s.rpcOutputBlock(block, false, false)
	}
	return nil, err
}

//返回给定块哈希和索引的uncle块，当fullTx为true时完整详细地返回块中的所有交易，否则仅返回交易哈希
func (s *PublicBlockChainAPI) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (map[string]interface{}, error) {
	block, err := s.b.GetBlock(ctx, blockHash)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			log.Debug("Requested uncle not found", "number", block.Number(), "hash", blockHash, "index", index)
			return nil, nil
		}
		block = types.NewBlockWithHeader(uncles[index])
		return s.rpcOutputBlock(block, false, false)
	}
	return nil, err
}

//返回给定块号的块中的叔号数
func (s *PublicBlockChainAPI) GetUncleCountByBlockNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}

//返回给定块散列的块中的叔号数
func (s *PublicBlockChainAPI) GetUncleCountByBlockHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.GetBlock(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}

//返回存储在给定块号的状态下给定地址的代码
func (s *PublicBlockChainAPI) GetCode(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumber(ctx, blockNr)
	if state == nil || err != nil {
		return nil, err
	}
	code := state.GetCode(address)
	return code, state.Error()
}

//从给定地址，key和的状态返回存储块号 rpc.LatestBlockNumber和rpc.PendingBlockNumber元块也允许使用数字。
func (s *PublicBlockChainAPI) GetStorageAt(ctx context.Context, address common.Address, key string, blockNr rpc.BlockNumber) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumber(ctx, blockNr)
	if state == nil || err != nil {
		return nil, err
	}
	res := state.GetState(address, common.HexToHash(key))
	return res[:], state.Error()
}

//****播客链新增处理****

//得到最后一次的出块节点
func (s *PublicBlockChainAPI) GetLastProducer(ctx context.Context) (common.Address, error) {

	block := s.b.CurrentBlock()
	if block != nil {
		return block.DposContext.GetLastProducer(-1)
	}
	return common.Address{}, errors.New("failed get last producer")
}

//得到最后一次的分币节点
func (s *PublicBlockChainAPI) GetLastTokenNoder(ctx context.Context) (common.Address, error) {

	block := s.b.CurrentBlock()
	if block != nil {
		return block.DposContext.GetLastTokenNoder(-1)
	}
	return common.Address{}, errors.New("failed get last token noder")
}

//得到下一次的出块节点
func (s *PublicBlockChainAPI) GetNextProducer(ctx context.Context) (common.Address, error) {

	block := s.b.CurrentBlock()
	if block != nil {
		return block.DposContext.GetCurrentProducer()
	}
	return common.Address{}, errors.New("failed get next producer")
}

//得到下一次的分币节点
func (s *PublicBlockChainAPI) GetNextTokenNoder(ctx context.Context) (common.Address, error) {

	block := s.b.CurrentBlock()
	if block != nil {
		return block.DposContext.GetCurrentTokenNoder()
	}
	return common.Address{}, errors.New("failed get next token noder")
}

//播客链新增函数处理，设置当前基础合约
func (s *PublicBlockChainAPI) SetBaseContracts(ctx context.Context, address common.Address, contractType protocol.ContractType, abiJson string) error {

	//log.Info("****SetBaseContracts****", "address", address.String(), "contractType", contractType, "abiJson", abiJson)

	//检测节点信息
	//log.Info("SetBaseContracts checkContract")
	if err := s.checkContract(); err != nil {
		log.Error("SetBaseContracts checkContract", "err", err)
		return err
	}

	//检测交易类型
	//log.Info("SetBaseContracts isExitsTxType")
	if err := s.isExitsTxType(contractType, protocol.ContractVote, protocol.ContractAssignToken); err != nil {
		log.Error("SetBaseContracts isExitsTxType", "err", err)
		return err
	}

	//类型转换
	//log.Info("SetBaseContracts txTypeRotate")
	txType := s.txTypeRotate(contractType)

	//产生一个交易
	//log.Info("SetBaseContracts SubmitBokerTransaction", "txType", txType)
	return s.b.Boker().SubmitBokerTransaction(ctx, txType, address, abiJson)
}

//播客链新增函数处理，取消一个基础合约
func (s *PublicBlockChainAPI) CancelBaseContracts(ctx context.Context, address common.Address, contractType protocol.ContractType) error {

	log.Info("****CancelBaseContracts****", "address", address.String(), "contractType", contractType)

	//检测节点信息
	log.Info("CancelBaseContracts checkContract")
	if err := s.checkContract(); err != nil {
		log.Error("CancelBaseContracts checkContract", "err", err)
		return err
	}

	//检测交易类型
	log.Info("CancelBaseContracts isExitsTxType")
	if err := s.isExitsTxType(contractType, protocol.UnContractVote, protocol.UnContractAssignToken); err != nil {
		log.Error("CancelBaseContracts isExitsTxType", "err", err)
		return err
	}

	//类型转换
	log.Info("CancelBaseContracts txTypeRotate")
	txType := s.txTypeRotate(contractType)

	//产生一个交易
	log.Info("CancelBaseContracts SubmitBokerTransaction", "txType", txType)
	return s.b.Boker().SubmitBokerTransaction(ctx, txType, address, "")
}

func (s *PublicBlockChainAPI) txTypeRotate(contractType protocol.ContractType) protocol.TxType {

	if contractType == protocol.ContractVote {
		return protocol.SetVote
	} else if contractType == protocol.ContractAssignToken {
		return protocol.SetAssignToken
	} else if contractType == protocol.UnContractVote {
		return protocol.CancelVote
	} else if contractType == protocol.UnContractAssignToken {
		return protocol.CanclAssignToken
	}
	return protocol.Binary
}

func (s *PublicBlockChainAPI) isExitsTxType(contractType protocol.ContractType, needTypes ...protocol.ContractType) error {

	var judge bool = false
	for _, needType := range needTypes {

		if needType == contractType {
			judge = true
			break
		}
	}
	if !judge {
		return errors.New("CoinBase Not`s Set Base Contracts Account")
	}
	return nil
}

func (s *PublicBlockChainAPI) checkContract() error {

	block := s.b.CurrentBlock()
	if block == nil {

		errors.New("failed baseContractsDeal")
	}

	//获取当前Coinbase
	coinbase, err := s.b.Coinbase()
	if err != nil {
		return err
	}

	//判断账号是否是验证人账号
	if !block.DposContext.IsValidator(coinbase) {
		return errors.New("Current coinbase Not`s Validator")
	}

	//判断账号是否是指定账号
	/*txLevel, err := s.b.Boker().GetAccount(coinbase)
	if err != nil {
		return err
	}

	if !bokerapi.ExistsTxType(protocol.SetVote, txLevel) &&
		!bokerapi.ExistsTxType(protocol.CancelVote, txLevel) &&
		!bokerapi.ExistsTxType(protocol.SetAssignToken, txLevel) &&
		!bokerapi.ExistsTxType(protocol.CanclAssignToken, txLevel) {

		return errors.New("CoinBase Not`s Set Base Contracts Account")
	}*/

	return nil
}

func (s *PublicBlockChainAPI) baseContractsDeal(ctx context.Context, address common.Address, abiJson string, contractType protocol.ContractType, needTypes ...protocol.ContractType) error {

	//判断是否属于要求的类型
	var judge bool = false
	for _, needType := range needTypes {

		if needType == contractType {
			judge = true
			break
		}
	}
	if !judge {
		return errors.New("CoinBase Not`s Set Base Contracts Account")
	}

	block := s.b.CurrentBlock()
	if block != nil {

		//获取当前Coinbase
		coinbase, err := s.b.Coinbase()
		if err != nil {
			return err
		}

		//判断账号是否是验证人账号
		if !block.DposContext.IsValidator(coinbase) {
			return errors.New("Current coinbase Not`s Validator")
		}

		//判断账号是否是指定账号
		txLevel, err := s.b.Boker().GetAccount(coinbase)
		if err != nil {
			return err
		}

		if !bokerapi.ExistsTxType(protocol.SetVote, txLevel) &&
			!bokerapi.ExistsTxType(protocol.CancelVote, txLevel) &&
			!bokerapi.ExistsTxType(protocol.SetAssignToken, txLevel) &&
			!bokerapi.ExistsTxType(protocol.CanclAssignToken, txLevel) {

			return errors.New("CoinBase Not`s Set Base Contracts Account")
		}

		//产生一个设置基础合约的交易
		if contractType == protocol.ContractVote {
			s.b.Boker().SubmitBokerTransaction(ctx, protocol.SetVote, address, abiJson)
		} else if contractType == protocol.ContractAssignToken {
			s.b.Boker().SubmitBokerTransaction(ctx, protocol.SetAssignToken, address, abiJson)
		} else if contractType == protocol.UnContractVote {
			s.b.Boker().SubmitBokerTransaction(ctx, protocol.CancelVote, address, abiJson)
		} else if contractType == protocol.UnContractAssignToken {
			s.b.Boker().SubmitBokerTransaction(ctx, protocol.CanclAssignToken, address, abiJson)
		}
		return nil
	}
	return errors.New("failed baseContractsDeal")
}

//播客链新增函数处理，添加一个验证者信息
func (s *PublicBlockChainAPI) AddValidator(ctx context.Context, address common.Address, votes *big.Int) error {

	//log.Info("****AddValidator****", "address", address.String(), "votes", votes.Int64())

	block, err := s.b.BlockByNumber(ctx, 0)
	if err != nil {
		return err
	}

	if block != nil {

		//获取当前Coinbase
		coinbase, err := s.b.Coinbase()
		if err != nil {
			return err
		}

		if s.b.Boker() == nil {
			log.Error("AddValidator error boker is nil")
			return nil
		}

		txLevel, err := s.b.Boker().GetAccount(coinbase)
		if err != nil {
			return err
		}

		//判断此账号是否具有设置验证者权限
		if !bokerapi.ExistsTxType(protocol.SetValidator, txLevel) {
			return errors.New("CoinBase Not`s Add Validator")
		}

		//判断此账号是否已经是验证者
		if block.DposContext.IsValidator(address) {
			return errors.New("Account has Validator")
		}

		//判断当前是否验证者已满
		if block.DposContext.IsValidatorFull() {
			return errors.New("Validator has Full")
		}

		//产生一个设置验证者的交易
		s.b.Boker().SubmitBokerTransaction(ctx, protocol.SetValidator, address, "")

		return nil
	}
	return errors.New("failed AddValidator")
}

//播客链新增函数处理，添加一个验证者信息
func (s *PublicBlockChainAPI) DecodeAbi(ctx context.Context, abiJson string, method string, payload string) error {

	log.Info("****DecodeAbi****", "abiJson", abiJson, "method", method, "payload", payload)

	_, err := protocol.DecodeAbi(abiJson, method, payload)
	if err != nil {
		return err
	}
	return nil

}

// CallArgs represents the arguments for a call.
type CallArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      hexutil.Big     `json:"gas"`
	GasPrice hexutil.Big     `json:"gasPrice"`
	Value    hexutil.Big     `json:"value"`
	Data     hexutil.Bytes   `json:"data"`
	Extra    hexutil.Bytes   `json:"extra"`
	TxType   protocol.TxType `json:"txType"`
}

func (s *PublicBlockChainAPI) doCall(ctx context.Context, args CallArgs, blockNr rpc.BlockNumber, vmCfg vm.Config) ([]byte, *big.Int, bool, error) {
	defer func(start time.Time) { log.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	log.Info("****doCall****")

	state, header, err := s.b.StateAndHeaderByNumber(ctx, blockNr)
	if state == nil || err != nil {
		return nil, common.Big0, false, err
	}
	// Set sender address or use a default if none specified
	addr := args.From
	if addr == (common.Address{}) {
		if wallets := s.b.AccountManager().Wallets(); len(wallets) > 0 {
			if accounts := wallets[0].Accounts(); len(accounts) > 0 {
				addr = accounts[0].Address
			}
		}
	}
	// Set default gas & gas price if none were set
	gas, gasPrice := args.Gas.ToInt(), args.GasPrice.ToInt()
	if gas.Sign() == 0 {
		gas = big.NewInt(50000000)
	}
	if gasPrice.Sign() == 0 {
		gasPrice = new(big.Int).SetUint64(defaultGasPrice)
	}

	// Create new call message
	msg := types.NewMessage(addr, args.To, 0, args.Value.ToInt(), gas, gasPrice, args.Data, args.Extra, false, args.TxType)

	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var cancel context.CancelFunc
	if vmCfg.DisableGasMetering {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer func() { cancel() }()

	// Get a new instance of the EVM.
	evm, vmError, err := s.b.GetEVM(ctx, msg, state, header, vmCfg)
	if err != nil {
		return nil, common.Big0, false, err
	}
	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		evm.Cancel()
	}()

	// Setup the gas pool (also for unmetered requests)
	// and apply the message.
	gp := new(core.GasPool).AddGas(math.MaxBig256)

	res, _, gas, failed, err := core.BinaryMessage(evm, msg, gp, s.b.Boker())
	if err := vmError(); err != nil {

		log.Error("doCall", "err", err)
		return nil, common.Big0, false, err
	}

	log.Info("doCall", "res", res, "resLength", len(res))
	return res, gas, failed, err
}

// Call executes the given transaction on the state for the given block number.
// It doesn't make and changes in the state/blockchain and is useful to execute and retrieve values.
func (s *PublicBlockChainAPI) Call(ctx context.Context, args CallArgs, blockNr rpc.BlockNumber) (hexutil.Bytes, error) {

	result, _, _, err := s.doCall(ctx, args, blockNr, vm.Config{DisableGasMetering: true})

	log.Info("****Call****", "result", result)
	return (hexutil.Bytes)(result), err
}

// EstimateGas returns an estimate of the amount of gas needed to execute the
// given transaction against the current pending block.
func (s *PublicBlockChainAPI) EstimateGas(ctx context.Context, args CallArgs) (*hexutil.Big, error) {
	// Determine the lowest and highest possible gas limits to binary search in between
	var (
		lo  uint64 = params.TxGas - 1
		hi  uint64
		cap uint64
	)
	if (*big.Int)(&args.Gas).Uint64() >= params.TxGas {
		hi = (*big.Int)(&args.Gas).Uint64()
	} else {
		// Retrieve the current pending block to act as the gas ceiling
		block, err := s.b.BlockByNumber(ctx, rpc.PendingBlockNumber)
		if err != nil {
			return nil, err
		}
		hi = block.GasLimit().Uint64()
	}
	cap = hi

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) bool {
		(*big.Int)(&args.Gas).SetUint64(gas)
		_, _, failed, err := s.doCall(ctx, args, rpc.PendingBlockNumber, vm.Config{})
		if err != nil || failed {
			return false
		}
		return true
	}
	// Execute the binary search and hone in on an executable gas limit
	for lo+1 < hi {
		mid := (hi + lo) / 2
		if !executable(mid) {
			lo = mid
		} else {
			hi = mid
		}
	}
	// Reject the transaction as invalid if it still fails at the highest allowance
	if hi == cap {
		if !executable(hi) {
			return nil, fmt.Errorf("gas required exceeds allowance or always failing transaction")
		}
	}
	return (*hexutil.Big)(new(big.Int).SetUint64(hi)), nil
}

// ExecutionResult groups all structured logs emitted by the EVM
// while replaying a transaction in debug mode as well as transaction
// execution status, the amount of gas used and the return value
type ExecutionResult struct {
	Gas         *big.Int       `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue string         `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

// StructLogRes stores a structured log emitted by the EVM while replaying a
// transaction in debug mode
type StructLogRes struct {
	Pc      uint64             `json:"pc"`
	Op      string             `json:"op"`
	Gas     uint64             `json:"gas"`
	GasCost uint64             `json:"gasCost"`
	Depth   int                `json:"depth"`
	Error   error              `json:"error,omitempty"`
	Stack   *[]string          `json:"stack,omitempty"`
	Memory  *[]string          `json:"memory,omitempty"`
	Storage *map[string]string `json:"storage,omitempty"`
}

// formatLogs formats EVM returned structured logs for json output
func FormatLogs(logs []vm.StructLog) []StructLogRes {
	formatted := make([]StructLogRes, len(logs))
	for index, trace := range logs {
		formatted[index] = StructLogRes{
			Pc:      trace.Pc,
			Op:      trace.Op.String(),
			Gas:     trace.Gas,
			GasCost: trace.GasCost,
			Depth:   trace.Depth,
			Error:   trace.Err,
		}
		if trace.Stack != nil {
			stack := make([]string, len(trace.Stack))
			for i, stackValue := range trace.Stack {
				stack[i] = fmt.Sprintf("%x", math.PaddedBigBytes(stackValue, 32))
			}
			formatted[index].Stack = &stack
		}
		if trace.Memory != nil {
			memory := make([]string, 0, (len(trace.Memory)+31)/32)
			for i := 0; i+32 <= len(trace.Memory); i += 32 {
				memory = append(memory, fmt.Sprintf("%x", trace.Memory[i:i+32]))
			}
			formatted[index].Memory = &memory
		}
		if trace.Storage != nil {
			storage := make(map[string]string)
			for i, storageValue := range trace.Storage {
				storage[fmt.Sprintf("%x", i)] = fmt.Sprintf("%x", storageValue)
			}
			formatted[index].Storage = &storage
		}
	}
	return formatted
}

// rpcOutputBlock converts the given block to the RPC output which depends on fullTx. If inclTx is true transactions are
// returned. When fullTx is true the returned block contains full transaction details, otherwise it will only contain
// transaction hashes.
func (s *PublicBlockChainAPI) rpcOutputBlock(b *types.Block, inclTx bool, fullTx bool) (map[string]interface{}, error) {

	head := b.Header()
	fields := map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             b.Hash(),
		"parentHash":       head.ParentHash,
		"nonce":            head.Nonce,
		"mixHash":          head.MixDigest,
		"sha3Uncles":       head.UncleHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"validator":        head.Validator,
		"coinbase":         head.Coinbase,
		"difficulty":       (*hexutil.Big)(head.Difficulty),
		"totalDifficulty":  (*hexutil.Big)(s.b.GetTd(b.Hash())),
		"extraData":        hexutil.Bytes(head.Extra),
		"size":             hexutil.Uint64(uint64(b.Size().Int64())),
		"gasLimit":         (*hexutil.Big)(head.GasLimit),
		"gasUsed":          (*hexutil.Big)(head.GasUsed),
		"timestamp":        (*hexutil.Big)(head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
		"dposProto":        head.DposProto.Root().String(),
		"bokerProto":       head.BokerProto.Root().String(),
	}

	if inclTx {
		formatTx := func(tx *types.Transaction) (interface{}, error) {
			return tx.Hash(), nil
		}

		if fullTx {
			formatTx = func(tx *types.Transaction) (interface{}, error) {
				return newRPCTransactionFromBlockHash(b, tx.Hash()), nil
			}
		}

		txs := b.Transactions()
		transactions := make([]interface{}, len(txs))
		var err error
		for i, tx := range b.Transactions() {
			if transactions[i], err = formatTx(tx); err != nil {
				return nil, err
			}
		}
		fields["transactions"] = transactions
	}

	uncles := b.Uncles()
	uncleHashes := make([]common.Hash, len(uncles))
	for i, uncle := range uncles {
		uncleHashes[i] = uncle.Hash()
	}
	fields["uncles"] = uncleHashes

	return fields, nil
}

//播客链中添加获取当前候选人相关信息
type ValidatorList struct {
	Address []common.Address `json:"address"`
}

func (s *PublicBlockChainAPI) GetCurrentValidator(ctx context.Context, blockNr rpc.BlockNumber) ([]byte, error) {

	//根据块号得到相关的块信息
	block, blockErr := s.b.BlockByNumber(ctx, blockNr)
	if block != nil {

		//根据块信息得到候选人列表信息
		validatorlist := ValidatorList{}
		validators, dposErr := block.DposContext.GetEpochTrie()
		if dposErr != nil {
			return nil, dposErr
		}

		//将候选人信息放入到json中，并进行打包
		for index, _ := range validators {
			validatorlist.Address = append(validatorlist.Address, validators[index])
		}
		jsonBytes, jsonErr := json.Marshal(validatorlist)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return jsonBytes, nil
	}
	return nil, blockErr
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	Type             protocol.TxType `json:"type"`
	BlockHash        common.Hash     `json:"blockHash"`
	BlockNumber      *hexutil.Big    `json:"blockNumber"`
	From             common.Address  `json:"from"`
	Gas              *hexutil.Big    `json:"gas"`
	GasPrice         *hexutil.Big    `json:"gasPrice"`
	Hash             common.Hash     `json:"hash"`
	Input            hexutil.Bytes   `json:"input"`
	Extra            hexutil.Bytes   `json:"extra"`
	Nonce            hexutil.Uint64  `json:"nonce"`
	To               *common.Address `json:"to"`
	TransactionIndex hexutil.Uint    `json:"transactionIndex"`
	Value            *hexutil.Big    `json:"value"`
	V                *hexutil.Big    `json:"v"`
	R                *hexutil.Big    `json:"r"`
	S                *hexutil.Big    `json:"s"`
}

func newRPCTransaction(tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64) *RPCTransaction {

	from, _ := types.Sender(types.HomesteadSigner{}, tx)
	v, r, s := tx.RawSignatureValues()

	result := &RPCTransaction{
		Type:     tx.Type(),
		From:     from,
		Gas:      (*hexutil.Big)(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    hexutil.Bytes(tx.Data()),
		Extra:    hexutil.Bytes(tx.Extra()),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(v),
		R:        (*hexutil.Big)(r),
		S:        (*hexutil.Big)(s),
	}
	if blockHash != (common.Hash{}) {
		result.BlockHash = blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = hexutil.Uint(index)
	}
	return result
}

// newRPCPendingTransaction returns a pending transaction that will serialize to the RPC representation
func newRPCPendingTransaction(tx *types.Transaction) *RPCTransaction {
	return newRPCTransaction(tx, common.Hash{}, 0, 0)
}

// newRPCTransactionFromBlockIndex returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockIndex(b *types.Block, index uint64) *RPCTransaction {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	return newRPCTransaction(txs[index], b.Hash(), b.NumberU64(), index)
}

func newRPCRawTransactionFromBlockIndex(b *types.Block, index uint64) hexutil.Bytes {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	blob, _ := rlp.EncodeToBytes(txs[index])
	return blob
}

func newRPCTransactionFromBlockHash(b *types.Block, hash common.Hash) *RPCTransaction {
	for idx, tx := range b.Transactions() {
		if tx.Hash() == hash {
			return newRPCTransactionFromBlockIndex(b, uint64(idx))
		}
	}
	return nil
}

// PublicTransactionPoolAPI exposes methods for the RPC interface
type PublicTransactionPoolAPI struct {
	b         Backend
	nonceLock *AddrLocker
}

// NewPublicTransactionPoolAPI creates a new RPC service with methods specific for the transaction pool.
func NewPublicTransactionPoolAPI(b Backend, nonceLock *AddrLocker) *PublicTransactionPoolAPI {
	return &PublicTransactionPoolAPI{b, nonceLock}
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
func (s *PublicTransactionPoolAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}

// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
func (s *PublicTransactionPoolAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.GetBlock(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}

// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
func (s *PublicTransactionPoolAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func (s *PublicTransactionPoolAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.GetBlock(ctx, blockHash); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
func (s *PublicTransactionPoolAPI) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
func (s *PublicTransactionPoolAPI) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.GetBlock(ctx, blockHash); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetTransactionCount returns the number of transactions the given address has sent for the given block number
func (s *PublicTransactionPoolAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (*hexutil.Uint64, error) {
	state, _, err := s.b.StateAndHeaderByNumber(ctx, blockNr)
	if state == nil || err != nil {
		return nil, err
	}
	nonce := state.GetNonce(address)
	return (*hexutil.Uint64)(&nonce), state.Error()
}

// GetTransactionByHash returns the transaction for the given hash
func (s *PublicTransactionPoolAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) *RPCTransaction {
	// Try to return an already finalized transaction
	if tx, blockHash, blockNumber, index := core.GetTransaction(s.b.ChainDb(), hash); tx != nil {
		return newRPCTransaction(tx, blockHash, blockNumber, index)
	}
	// No finalized transaction, try to retrieve it from the pool
	if tx := s.b.GetPoolTransaction(hash); tx != nil {
		return newRPCPendingTransaction(tx)
	}
	// Transaction unknown, return as such
	return nil
}

// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
func (s *PublicTransactionPoolAPI) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	var tx *types.Transaction

	// Retrieve a finalized transaction, or a pooled otherwise
	if tx, _, _, _ = core.GetTransaction(s.b.ChainDb(), hash); tx == nil {
		if tx = s.b.GetPoolTransaction(hash); tx == nil {
			// Transaction not found anywhere, abort
			return nil, nil
		}
	}

	// Serialize to RLP and return
	return rlp.EncodeToBytes(tx)
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
func (s *PublicTransactionPoolAPI) GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error) {
	tx, blockHash, blockNumber, index := core.GetTransaction(s.b.ChainDb(), hash)
	if tx == nil {
		return nil, nil
	}
	receipt, _, _, _ := core.GetReceipt(s.b.ChainDb(), hash) // Old receipts don't have the lookup data available
	from, _ := types.Sender(types.HomesteadSigner{}, tx)

	fields := map[string]interface{}{
		"blockHash":         blockHash,
		"blockNumber":       hexutil.Uint64(blockNumber),
		"transactionHash":   hash,
		"transactionIndex":  hexutil.Uint64(index),
		"type":              tx.Type(),
		"from":              from,
		"to":                tx.To(),
		"extra":             hexutil.Bytes(tx.Extra()),
		"gasUsed":           (*hexutil.Big)(receipt.GasUsed),
		"cumulativeGasUsed": (*hexutil.Big)(receipt.CumulativeGasUsed),
		"contractAddress":   nil,
		"logs":              receipt.Logs,
		"logsBloom":         receipt.Bloom,
	}

	// Assign receipt status or post state.
	if len(receipt.PostState) > 0 {
		fields["root"] = hexutil.Bytes(receipt.PostState)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}
	if receipt.Logs == nil {
		fields["logs"] = [][]*types.Log{}
	}
	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}
	return fields, nil
}

// sign is a helper function that signs a transaction with the private key of the given address.
func (s *PublicTransactionPoolAPI) sign(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {

	if err := tx.Validate(); err != nil {
		return nil, err
	}
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}

	// Request the wallet to sign the transaction
	var chainID *big.Int
	if config := s.b.ChainConfig(); config.IsEIP155(s.b.CurrentBlock().Number()) {
		chainID = config.ChainId
	}
	return wallet.SignTx(account, tx, chainID)
}

// SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
type SendTxArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Big    `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Data     hexutil.Bytes   `json:"data"`
	//Extra    hexutil.Bytes   `json:"extra"`
	Nonce *hexutil.Uint64 `json:"nonce"`
	Type  protocol.TxType `json:"type"`
}

// prepareSendTxArgs is a helper function that fills in default values for unspecified tx fields.
func (args *SendTxArgs) SetDefaults(ctx context.Context, b Backend) error {

	//如果Gas为空，则给一个默认的Gas（defaultGas = 90000）
	if args.Gas == nil {
		args.Gas = (*hexutil.Big)(big.NewInt(defaultGas))
	}

	//如果GasPrice是空，则给一个建议的GasPrice
	if args.GasPrice == nil {
		price, err := b.SuggestPrice(ctx)
		if err != nil {
			return err
		}
		args.GasPrice = (*hexutil.Big)(price)
	}

	//如果Value为空则产生一个Value
	if args.Value == nil {
		args.Value = new(hexutil.Big)
	}

	//如果Nonce为空，则产生一个Nonce
	if args.Nonce == nil {
		nonce, err := b.GetPoolNonce(ctx, args.From)
		if err != nil {
			return err
		}
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}
	log.Info("SetDefaults", "Nonce", args.Nonce)

	return nil
}

//这里需要进行判断
func (args *SendTxArgs) ToTransaction() (*types.Transaction, error) {

	//log.Info("****SendTxArgs ToTransaction****")

	//判断交易地址是否为空
	if args.To == nil {

		if (args.Type >= protocol.SetVote) && (args.Type <= protocol.AssignReward) {

			//return types.NewBaseContractCreation(uint64(*args.Nonce), (*big.Int)(args.Value), args.Data), nil
			return nil, errors.New("base contract transaction type not found contract address")
		} else if types.IsBinary(args.Type) {

			return types.NewContractCreation(uint64(*args.Nonce), (*big.Int)(args.Value), (*big.Int)(args.Gas), (*big.Int)(args.GasPrice), args.Data), nil
		}
		return nil, errors.New("unknown transaction type")
	}

	//设置交易地址
	to := common.Address{}
	if args.To != nil {
		to = *args.To
	}
	return types.NewTransaction(args.Type, uint64(*args.Nonce), to, (*big.Int)(args.Value), (*big.Int)(args.Gas), (*big.Int)(args.GasPrice), args.Data), nil
}

//submitTransaction是一个辅助函数，它将tx提交给txPool并记录消息。
func SubmitTransaction(ctx context.Context, b Backend, tx *types.Transaction) (common.Hash, error) {

	//log.Info("****SubmitTransaction****", "gas", tx.Gas(), "gasprice", tx.GasPrice(), "hash", tx.Hash().String())

	//判断交易类型是否是限定的类型
	if err := tx.Validate(); err != nil {
		log.Error("SubmitTransaction Validate", "error", err)
		return common.Hash{}, err
	}

	//发送交易
	if err := b.SendTx(ctx, tx); err != nil {
		log.Error("SubmitTransaction SendTx", "error", err, "txType", tx.Type())
		return common.Hash{}, err
	}

	//如果to为空得到签名者，并进行签名
	if tx.To() == nil {
		from, err := types.Sender(types.HomesteadSigner{}, tx)
		if err != nil {
			log.Error("SubmitTransaction Sender", "error", err)
			return common.Hash{}, err
		}
		addr := crypto.CreateAddress(from, tx.Nonce())
		log.Info("Submitted contract creation", "fullhash", tx.Hash().Hex(), "contract", addr.Hex())
	} else {
		log.Info("Submitted transaction", "fullhash", tx.Hash().Hex(), "recipient", tx.To())
	}
	return tx.Hash(), nil
}

//用户通过JSON RPC发起eth_sendTransaction请求，最终会调用PublicTransactionPoolAPI
//SendTransaction为给定的参数创建一个交易，对其进行签名并将其提交给交易池。
func (s *PublicTransactionPoolAPI) SendTransaction(ctx context.Context, args SendTxArgs) (common.Hash, error) {

	//这里需要将bin格式的Data进行转换，得到最后一个参数

	//log.Info("****SendTransaction****", "Nonce", args.Nonce.String(), "from", args.From, "Gas", args.Gas, "GasPrice", args.GasPrice, "to", args.To, "json", args)
	account := accounts.Account{Address: args.From}
	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		log.Error("SendTransaction", "from", args.From, "error", err.Error())
		return common.Hash{}, err
	}

	if args.Nonce == nil {
		s.nonceLock.LockAddr(args.From)
		defer s.nonceLock.UnlockAddr(args.From)
	}

	// Set some sanity defaults and terminate on failure
	if err := args.SetDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}

	// Assemble the transaction and sign with the wallet
	tx, tranErr := args.ToTransaction()
	if tranErr != nil {
		return common.Hash{}, tranErr
	}

	var chainID *big.Int
	if config := s.b.ChainConfig(); config.IsEIP155(s.b.CurrentBlock().Number()) {
		chainID = config.ChainId
	}
	signed, err := wallet.SignTx(account, tx, chainID)
	if err != nil {
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.b, signed)
}

// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *PublicTransactionPoolAPI) SendRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error) {

	//log.Info("SendRawTransaction", "len", len(encodedTx), "encodedTx", encodedTx)
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {

		log.Error("SendRawTransaction", "error", err, "encodedTx", encodedTx)
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.b, tx)
}

// Sign calculates an ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message).
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The account associated with addr must be unlocked.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (s *PublicTransactionPoolAPI) Sign(addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Sign the requested hash with the wallet
	signature, err := wallet.SignHash(account, signHash(data))
	if err == nil {
		signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	}
	return signature, err
}

// SignTransactionResult represents a RLP encoded signed transaction.
type SignTransactionResult struct {
	Raw hexutil.Bytes      `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}

// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.
func (s *PublicTransactionPoolAPI) SignTransaction(ctx context.Context, args SendTxArgs) (*SignTransactionResult, error) {

	log.Info("****PublicTransactionPoolAPI SignTransaction****")

	if args.Nonce == nil {
		// Hold the addresse's mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.From)
		defer s.nonceLock.UnlockAddr(args.From)
	}
	if err := args.SetDefaults(ctx, s.b); err != nil {
		return nil, err
	}

	trans, err := args.ToTransaction()
	if err != nil {
		return nil, err
	}

	//tx, err := s.sign(args.From, args.toTransaction())
	tx, err := s.sign(args.From, trans)
	if err != nil {
		return nil, err
	}
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, tx}, nil
}

// PendingTransactions returns the transactions that are in the transaction pool and have a from address that is one of
// the accounts this node manages.
func (s *PublicTransactionPoolAPI) PendingTransactions() ([]*RPCTransaction, error) {
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return nil, err
	}

	transactions := make([]*RPCTransaction, 0, len(pending))
	for _, tx := range pending {
		var signer types.Signer = types.HomesteadSigner{}
		/*if tx.Protected() {
			signer = types.NewEIP155Signer(tx.ChainId())
		}*/
		from, _ := types.Sender(signer, tx)
		if _, err := s.b.AccountManager().Find(accounts.Account{Address: from}); err == nil {
			transactions = append(transactions, newRPCPendingTransaction(tx))
		}
	}
	return transactions, nil
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (s *PublicTransactionPoolAPI) Resend(ctx context.Context, sendArgs SendTxArgs, gasPrice, gasLimit *hexutil.Big) (common.Hash, error) {

	if sendArgs.Nonce == nil {
		return common.Hash{}, fmt.Errorf("missing transaction nonce in transaction spec")
	}
	if err := sendArgs.SetDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}
	matchTx, err := sendArgs.ToTransaction()
	if err != nil {
		return common.Hash{}, err
	}

	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return common.Hash{}, err
	}

	for _, p := range pending {
		var signer types.Signer = types.HomesteadSigner{}
		/*if p.Protected() {
			signer = types.NewEIP155Signer(p.ChainId())
		}*/
		wantSigHash := signer.Hash(matchTx)

		if pFrom, err := types.Sender(signer, p); err == nil && pFrom == sendArgs.From && signer.Hash(p) == wantSigHash {
			// Match. Re-sign and send the transaction.
			if gasPrice != nil {
				sendArgs.GasPrice = gasPrice
			}
			if gasLimit != nil {
				sendArgs.Gas = gasLimit
			}

			trans, err := sendArgs.ToTransaction()
			if err != nil {
				return common.Hash{}, err
			}

			//signedTx, err := s.sign(sendArgs.From, sendArgs.toTransaction())
			signedTx, err := s.sign(sendArgs.From, trans)
			if err != nil {
				return common.Hash{}, err
			}

			log.Info("****Resend****", "Nonce", signedTx.Nonce())
			if err = s.b.SendTx(ctx, signedTx); err != nil {
				return common.Hash{}, err
			}
			return signedTx.Hash(), nil
		}
	}

	return common.Hash{}, fmt.Errorf("Transaction %#x not found", matchTx.Hash())
}

// PublicDebugAPI is the collection of Ethereum APIs exposed over the public
// debugging endpoint.
type PublicDebugAPI struct {
	b Backend
}

// NewPublicDebugAPI creates a new API definition for the public debug methods
// of the Ethereum service.
func NewPublicDebugAPI(b Backend) *PublicDebugAPI {
	return &PublicDebugAPI{b: b}
}

// GetBlockRlp retrieves the RLP encoded for of a single block.
func (api *PublicDebugAPI) GetBlockRlp(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	encoded, err := rlp.EncodeToBytes(block)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", encoded), nil
}

// PrintBlock retrieves a block and returns its pretty printed form.
func (api *PublicDebugAPI) PrintBlock(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	return block.String(), nil
}

// PrivateDebugAPI is the collection of Ethereum APIs exposed over the private
// debugging endpoint.
type PrivateDebugAPI struct {
	b Backend
}

// NewPrivateDebugAPI creates a new API definition for the private debug methods
// of the Ethereum service.
func NewPrivateDebugAPI(b Backend) *PrivateDebugAPI {
	return &PrivateDebugAPI{b: b}
}

// ChaindbProperty returns leveldb properties of the chain database.
func (api *PrivateDebugAPI) ChaindbProperty(property string) (string, error) {
	ldb, ok := api.b.ChainDb().(interface {
		LDB() *leveldb.DB
	})
	if !ok {
		return "", fmt.Errorf("chaindbProperty does not work for memory databases")
	}
	if property == "" {
		property = "leveldb.stats"
	} else if !strings.HasPrefix(property, "leveldb.") {
		property = "leveldb." + property
	}
	return ldb.LDB().GetProperty(property)
}

func (api *PrivateDebugAPI) ChaindbCompact() error {
	ldb, ok := api.b.ChainDb().(interface {
		LDB() *leveldb.DB
	})
	if !ok {
		return fmt.Errorf("chaindbCompact does not work for memory databases")
	}
	for b := byte(0); b < 255; b++ {
		log.Info("Compacting chain database", "range", fmt.Sprintf("0x%0.2X-0x%0.2X", b, b+1))
		err := ldb.LDB().CompactRange(util.Range{Start: []byte{b}, Limit: []byte{b + 1}})
		if err != nil {
			log.Error("Database compaction failed", "err", err)
			return err
		}
	}
	return nil
}

// SetHead rewinds the head of the blockchain to a previous block.
func (api *PrivateDebugAPI) SetHead(number hexutil.Uint64) {
	api.b.SetHead(uint64(number))
}

// PublicNetAPI offers network related RPC methods
type PublicNetAPI struct {
	net            *p2p.Server
	networkVersion uint64
}

// NewPublicNetAPI creates a new net API instance.
func NewPublicNetAPI(net *p2p.Server, networkVersion uint64) *PublicNetAPI {
	return &PublicNetAPI{net, networkVersion}
}

// Listening returns an indication if the node is listening for network connections.
func (s *PublicNetAPI) Listening() bool {
	return true // always listening
}

// PeerCount returns the number of connected peers
func (s *PublicNetAPI) PeerCount() hexutil.Uint {
	return hexutil.Uint(s.net.PeerCount())
}

// Version returns the current ethereum protocol version.
func (s *PublicNetAPI) Version() string {
	return fmt.Sprintf("%d", s.networkVersion)
}
