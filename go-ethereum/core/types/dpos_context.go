package types

import (
	_ "bytes"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"time"

	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/crypto/sha3"
	"github.com/boker/go-ethereum/ethdb"
	"github.com/boker/go-ethereum/include"
	"github.com/boker/go-ethereum/rlp"
	"github.com/boker/go-ethereum/trie"
)

type DposContext struct {
	epochTrie     *trie.Trie //记录每个周期的验证人Hash树
	validatorTrie *trie.Trie //验证人以及对应投票人Hash树
	blockCntTrie  *trie.Trie //记录验证人在周期内的出块数目Hash树
	db            ethdb.Database
}

var (
	epochPrefix     = []byte("epoch-")     //存放周期信息
	validatorPrefix = []byte("validator-") //存放验证者投票信息
	blockCntPrefix  = []byte("blockCnt-")  //存放投票数量
)

func NewEpochTrie(root common.Hash, db ethdb.Database) (*trie.Trie, error) {
	return trie.NewTrieWithPrefix(root, epochPrefix, db)
}

func NewValidatorTrie(root common.Hash, db ethdb.Database) (*trie.Trie, error) {
	return trie.NewTrieWithPrefix(root, validatorPrefix, db)
}

func NewBlockCntTrie(root common.Hash, db ethdb.Database) (*trie.Trie, error) {
	return trie.NewTrieWithPrefix(root, blockCntPrefix, db)
}

func NewDposContext(db ethdb.Database) (*DposContext, error) {

	epochTrie, err := NewEpochTrie(common.Hash{}, db)
	if err != nil {
		return nil, err
	}

	validatorTrie, err := NewValidatorTrie(common.Hash{}, db)
	if err != nil {
		return nil, err
	}

	blockCntTrie, err := NewBlockCntTrie(common.Hash{}, db)
	if err != nil {
		return nil, err
	}
	return &DposContext{
		epochTrie:     epochTrie,
		validatorTrie: validatorTrie,
		blockCntTrie:  blockCntTrie,
		db:            db,
	}, nil
}

func NewDposContextFromProto(db ethdb.Database, ctxProto *DposContextProto) (*DposContext, error) {

	epochTrie, err := NewEpochTrie(ctxProto.EpochHash, db)
	if err != nil {
		return nil, err
	}

	validatorTrie, err := NewValidatorTrie(ctxProto.ValidatorHash, db)
	if err != nil {
		return nil, err
	}

	blockCntTrie, err := NewBlockCntTrie(ctxProto.BlockCntHash, db)
	if err != nil {
		return nil, err
	}

	return &DposContext{
		epochTrie:     epochTrie,
		validatorTrie: validatorTrie,
		blockCntTrie:  blockCntTrie,
		db:            db,
	}, nil

}

func (d *DposContext) Copy() *DposContext {

	epochTrie := *d.epochTrie
	validatorTrie := *d.validatorTrie
	blockCntTrie := *d.blockCntTrie
	return &DposContext{
		epochTrie:     &epochTrie,
		validatorTrie: &validatorTrie,
		blockCntTrie:  &blockCntTrie,
	}
}

func (d *DposContext) Root() (h common.Hash) {

	hw := sha3.NewKeccak256()
	rlp.Encode(hw, d.epochTrie.Hash())
	rlp.Encode(hw, d.validatorTrie.Hash())
	rlp.Encode(hw, d.blockCntTrie.Hash())
	hw.Sum(h[:0])
	return h
}

func (d *DposContext) Snapshot() *DposContext {
	return d.Copy()
}

func (d *DposContext) RevertToSnapShot(snapshot *DposContext) {

	d.epochTrie = snapshot.epochTrie
	d.validatorTrie = snapshot.validatorTrie
	d.blockCntTrie = snapshot.blockCntTrie
}

func (d *DposContext) FromProto(dcp *DposContextProto) error {

	var err error
	d.epochTrie, err = NewEpochTrie(dcp.EpochHash, d.db)
	if err != nil {
		return err
	}

	d.validatorTrie, err = NewValidatorTrie(dcp.ValidatorHash, d.db)
	if err != nil {
		return err
	}

	d.blockCntTrie, err = NewBlockCntTrie(dcp.BlockCntHash, d.db)
	return err
}

type DposContextProto struct {
	EpochHash     common.Hash `json:"epochRoot"        gencodec:"required"`
	ValidatorHash common.Hash `json:"validatorRoot"     gencodec:"required"`
	BlockCntHash  common.Hash `json:"blockCntRoot"      gencodec:"required"`
}

func (d *DposContext) ToProto() *DposContextProto {

	return &DposContextProto{
		EpochHash:     d.epochTrie.Hash(),
		ValidatorHash: d.validatorTrie.Hash(),
		BlockCntHash:  d.blockCntTrie.Hash(),
	}
}

func (p *DposContextProto) Root() (h common.Hash) {

	hw := sha3.NewKeccak256()
	rlp.Encode(hw, p.EpochHash)
	rlp.Encode(hw, p.ValidatorHash)
	rlp.Encode(hw, p.BlockCntHash)
	hw.Sum(h[:0])
	return h
}

func (d *DposContext) CommitTo(dbw trie.DatabaseWriter) (*DposContextProto, error) {

	epochRoot, err := d.epochTrie.CommitTo(dbw)
	if err != nil {
		return nil, err
	}

	validatorRoot, err := d.validatorTrie.CommitTo(dbw)
	if err != nil {
		return nil, err
	}

	blockCntRoot, err := d.blockCntTrie.CommitTo(dbw)
	if err != nil {
		return nil, err
	}
	return &DposContextProto{
		EpochHash:     epochRoot,
		ValidatorHash: validatorRoot,
		BlockCntHash:  blockCntRoot,
	}, nil
}

func (d *DposContext) EpochTrie() *trie.Trie             { return d.epochTrie }
func (d *DposContext) ValidatorTrie() *trie.Trie         { return d.validatorTrie }
func (d *DposContext) BlockCntTrie() *trie.Trie          { return d.blockCntTrie }
func (d *DposContext) DB() ethdb.Database                { return d.db }
func (d *DposContext) SetEpoch(epoch *trie.Trie)         { d.epochTrie = epoch }
func (d *DposContext) SetValidator(validator *trie.Trie) { d.validatorTrie = validator }
func (d *DposContext) SetMintCnt(blockCnt *trie.Trie)    { d.blockCntTrie = blockCnt }

func (dc *DposContext) GetEpochTrie() ([]common.Address, error) {

	var validators []common.Address
	key := []byte("validators")
	validatorsRLP := dc.epochTrie.Get(key)
	if err := rlp.DecodeBytes(validatorsRLP, &validators); err != nil {
		return nil, fmt.Errorf("failed to decode validators: %s", err)
	}
	return validators, nil
}

func (dc *DposContext) SetEpochTrie(validators []common.Address) error {

	key := []byte("validators")
	validatorsRLP, err := rlp.EncodeToBytes(validators)
	if err != nil {
		return fmt.Errorf("failed to encode validators to rlp bytes: %s", err)
	}
	dc.epochTrie.Update(key, validatorsRLP)

	_, errCommit := dc.epochTrie.CommitTo(dc.db)
	return errCommit
}

func (dc *DposContext) CleanEpochTrie() error {

	err := dc.epochTrie.TryDelete([]byte("validators"))
	if err != nil {
		return err
	}
	_, errCommit := dc.epochTrie.CommitTo(dc.db)
	return errCommit
}

func (dc *DposContext) DeleteEpochTrie(validator common.Address) error {

	validators, err := dc.GetEpochTrie()
	if err != nil {
		return errors.New("failed to get epoch")
	}
	for index, producer := range validators {
		if producer == validator {
			validators = append(validators[:index], validators[index+1:]...)
		}
	}
	return dc.SetEpochTrie(validators)
}

func (dc *DposContext) DeleteValidator(validator common.Address) error {

	err := dc.validatorTrie.TryDelete(validator.Bytes())
	if err != nil {
		if _, ok := err.(*trie.MissingNodeError); !ok {
			return err
		}
	}
	_, errCommit := dc.validatorTrie.CommitTo(dc.db)
	return errCommit
}

func (dc *DposContext) CleanValidators() error {

	validators, err := dc.GetEpochTrie()
	if err != nil {
		return errors.New("failed to clean validator")
	}
	for _, validator := range validators {

		err := dc.validatorTrie.TryDelete(validator.Bytes())
		if err != nil {
			if _, ok := err.(*trie.MissingNodeError); !ok {
				return err
			}
		}
	}
	_, errCommit := dc.validatorTrie.CommitTo(dc.db)
	return errCommit
}

func (dc *DposContext) InsertValidator(validator common.Address, votes *big.Int) error {

	//得到当前所有的验证者
	validators, err := dc.GetEpochTrie()
	if err != nil {
		return err
	}

	//排序
	producers := sortableAddresses{}
	for _, validator := range validators {

		cnt, err := dc.GetValidatorCnt(validator)
		if err != nil {
			return err
		}
		producers = append(producers, &sortableAddress{validator, cnt})
	}
	sort.Sort(producers)

	if len(validators) >= include.MaxValidatorSize {

		//这里需要删除最小Votes数量的验证者
		if producers[len(producers)-1].weight.Cmp(votes) >= 0 {
			return errors.New("Votes Less than Validator Min")
		} else {

			//将最小踢出，并将新的加入
			dc.DeleteValidator(producers[len(producers)-1].address)

			//加入
			validators = append(validators, validator)
			dc.SetEpochTrie(validators)
			key := validator.Bytes()
			validatorRLP, err := rlp.EncodeToBytes(votes.String())
			if err != nil {
				return fmt.Errorf("failed to encode validators to rlp bytes: %s", err)
			}
			if err := dc.validatorTrie.TryUpdate(key, validatorRLP); err != nil {
				return err
			}
			_, errCommit := dc.validatorTrie.CommitTo(dc.db)
			return errCommit
		}
	} else {

		//加入
		validators = append(validators, validator)
		dc.SetEpochTrie(validators)
		key := validator.Bytes()
		validatorRLP, err := rlp.EncodeToBytes(votes.String())
		if err != nil {
			return fmt.Errorf("failed to encode validators to rlp bytes: %s", err)
		}
		if err := dc.validatorTrie.TryUpdate(key, validatorRLP); err != nil {
			return err
		}
		_, errCommit := dc.validatorTrie.CommitTo(dc.db)
		return errCommit
	}
	return nil
}

func (dc *DposContext) GetValidatorCnt(validator common.Address) (*big.Int, error) {

	//根据地址获取数据
	key := validator.Bytes()
	v, err := dc.epochTrie.TryGet(key)
	if err != nil {
		return big.NewInt(0), err
	}

	//转换成交易类型
	cnt, err := strconv.Atoi(string(v[:]))
	return big.NewInt(int64(cnt)), nil
}

func (dc *DposContext) SetValidatorVotes(validators []common.Address, votes []*big.Int) error {

	//踢出所有的验证者
	producers, err := dc.GetEpochTrie()
	if err != nil {
		return fmt.Errorf("failed to get validator: %s", err)
	}
	for _, validator := range producers {

		err := dc.validatorTrie.TryDelete(validator.Bytes())
		if err != nil {
			if _, ok := err.(*trie.MissingNodeError); !ok {
				return err
			}
		}
	}
	_, err = dc.validatorTrie.CommitTo(dc.db)
	if err != nil {
		return err
	}

	//重建验证者和投票信息
	key := []byte("validators")
	validatorsRLP, err := rlp.EncodeToBytes(validators)
	if err != nil {
		return fmt.Errorf("failed to encode validators to rlp bytes: %s", err)
	}
	dc.epochTrie.Update(key, validatorsRLP)
	_, err = dc.epochTrie.CommitTo(dc.db)
	if err != nil {
		return err
	}

	for index, validator := range validators {

		cnt := votes[index].Int64()
		if err = dc.validatorTrie.TryUpdate(validator.Bytes(), []byte(strconv.Itoa(int(cnt)))); err != nil {

			return fmt.Errorf("failed to TryUpdate validator: %s", err)
		}
	}
	_, err = dc.validatorTrie.CommitTo(dc.db)
	if err != nil {
		return err
	}

	return nil
}

func (dc *DposContext) IsValidator(address common.Address) bool {

	validators, err := dc.GetEpochTrie()
	if err != nil {
		return false
	}
	for _, v := range validators {
		if address == v {
			return true
		}
	}
	return false
}

func (dc *DposContext) GetCurrentProducer() (common.Address, error) {

	producers, err := dc.GetEpochTrie()
	if err != nil {
		return common.Address{}, errors.New("failed to GetValidators")
	}
	producerSize := len(producers)
	if producerSize == 0 {
		return common.Address{}, errors.New("failed to producers length is zero")
	}

	offset := time.Now().Unix() % include.EpochInterval
	offset /= include.ProducerInterval

	offset %= int64(producerSize)
	return producers[offset], nil
}

func (dc *DposContext) GetCurrentTokenNoder() (common.Address, error) {

	producers, err := dc.GetEpochTrie()
	if err != nil {
		return common.Address{}, errors.New("failed to GetValidators")
	}
	producerSize := len(producers)
	if producerSize == 0 {
		return common.Address{}, errors.New("failed to producers length is zero")
	}

	offset := time.Now().Unix() % include.EpochInterval
	offset /= include.TokenNoderInterval

	offset %= int64(producerSize)
	return producers[offset], nil
}

//根据时间得到当时的打包节点
func (dc *DposContext) GetProducer(now int64) (producer common.Address, err error) {

	producer = common.Address{}
	offset := now % include.EpochInterval
	if offset%include.ProducerInterval != 0 {
		return common.Address{}, include.ErrInvalidMintBlockTime
	}
	offset /= include.ProducerInterval

	//得到验证者数组
	producers, err := dc.GetEpochTrie()
	if err != nil {
		return common.Address{}, err
	}

	//得到验证者数量并判断验证者数量是否为0
	producerSize := len(producers)
	if producerSize == 0 {
		return common.Address{}, errors.New("failed to lookup producer")
	}

	//根据当前的移位偏移量确定当前应该出块的验证者
	offset %= int64(producerSize)
	return producers[offset], nil
}

func (dc *DposContext) GetTokenNoder(now int64) (tokennoder common.Address, err error) {

	tokennoder = common.Address{}
	offset := now % include.EpochInterval
	if offset%include.TokenNoderInterval != 0 {
		return common.Address{}, include.ErrInvalidMintBlockTime
	}
	offset /= include.TokenNoderInterval

	tokennoders, err := dc.GetEpochTrie()
	if err != nil {
		return common.Address{}, err
	}

	tokennoderSize := len(tokennoders)
	if tokennoderSize == 0 {
		return common.Address{}, errors.New("failed to lookup token node")
	}
	offset %= int64(tokennoderSize)
	return tokennoders[offset], nil
}

type sortableAddress struct {
	address common.Address
	weight  *big.Int
}
type sortableAddresses []*sortableAddress

func (p sortableAddresses) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p sortableAddresses) Len() int      { return len(p) }
func (p sortableAddresses) Less(i, j int) bool {
	if p[i].weight.Cmp(p[j].weight) < 0 {
		return false
	} else if p[i].weight.Cmp(p[j].weight) > 0 {
		return true
	} else {
		return p[i].address.String() < p[j].address.String()
	}
}
