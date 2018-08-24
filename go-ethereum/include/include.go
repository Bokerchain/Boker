package include

import (
	"errors"
	"math/big"
	"time"

	"github.com/boker/go-ethereum/common"
)

const (
	ExtraVanity        = 32           //扩展字段的前缀字节数量
	ExtraSeal          = 65           //扩展字段的后缀字节数量
	InmemorySignatures = 4096         //保留在内存中的最近块签名的数量
	ProducerInterval   = int64(10)    //打包时间间隔（秒）
	TokenNoderInterval = int64(300)   //分配通证时间间隔(秒)
	EpochInterval      = int64(86400) //一个周期的时间（86400秒 = 1天）
	MaxValidatorSize   = 1            //DPOS的验证者数量
	SafeSize           = 1            //安全的验证者数量
	ConsensusSize      = 1            //共识确认验证者数量
)

const (
	AssignTokenInterval = time.Second //分配通证时间间隔(秒)
	VotesInterval       = time.Second //投票时间间隔(秒)
)

var (
	BobbyUnit          *big.Int = big.NewInt(1e+17) //Bobby的单位
	BobbyMultiple      *big.Int = big.NewInt(220)   //倍数
	TransferUnit       *big.Int = big.NewInt(1e+17) //转账单位(这个数值仅用于每次给指定账号，方便指定账号给用户分配通证)
	TransferMultiple   *big.Int = big.NewInt(330)   //转账倍数
	TimeOfFirstBlock            = int64(0)          //创世区块的时间偏移量
	ConfirmedBlockHead          = []byte("confirmed-block-head")
)

var (
	RegisterCandidateMethod = "registerCandidate" //候选人注册方法名
	VoteCandidateMethod     = "voteCandidate"     //投票方法名
	RotateVoteMethod        = "rotateVote"        //转换投票方法名
	AssignTokenMethod       = "assignToken"       //分配通证
	TickCandidateMethod     = "tickVote"          //投票时钟
	GetCandidateMethod      = "getCandidates"     //获取候选人结果
)

var (
	ErrInvalidMintBlockTime = errors.New("invalid time to mint the block")  //不正确的出块时间
	ErrInvalidProducer      = errors.New("invalid current producer")        //出块节点出错
	ErrInvalidCoinbase      = errors.New("invalid current mining coinbase") //当前挖矿账号错误
	ErrInvalidSystem        = errors.New("invalid current system votes")    //
	ErrInvalidTokenNoder    = errors.New("invalid current token noder")
)

type BokerConfig struct {
	Address common.Address
}
