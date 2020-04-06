package block

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"github.com/corgi-kx/blockchain_golang/util"
	log "github.com/corgi-kx/logcustom"
	"math/big"
)

//工作量证明(pow)结构体
type proofOfWork struct {
	*Block
	//难度
	Target *big.Int
}

// TODO 获取POW实例
func NewProofOfWork(block *Block) *proofOfWork {
	target := big.NewInt(1)

	//返回一个大数(1 << 256-TargetBits)
	target.Lsh(target, 256-TargetBits)
	pow := &proofOfWork{block, target}
	return pow
}

var MineFlag = false

//var MineFlag = make(chan bool,1)

var MineReturnBH BlockHeader
var GlobalErr error

//用户提交的区块会进入该通道，只接受第一个合法区块
var MinedChan = make(chan *BlockHeader, 1)

//进行hash运算,获取到当前区块的hash值
func (p *proofOfWork) run(wsend WebsocketSender) (*BlockHeader, error) {
	//MinedChain =
	if len(MinedChan) == 1 {
		<-MinedChan
	}
	//这里可以定时发送，每个新连接的用户都有机会挖矿
	wsend.SendBlockHeaderToUser(p.BBlockHeader)

	MineReturnBlockHeader := <-MinedChan
	//其他云计算节点已经挖到区块了
	if MineReturnBlockHeader == nil {
		return nil, errors.New("其他云计算节点已经挖到区块了")
	}
	//返回应该需要返回一个随机数、一个幻方、区块hash、用户交易
	return MineReturnBlockHeader, nil
}

//func (p *proofOfWork) run(wsend WebsocketSender) (BlockHeader, error) {
//	//MinedChan := make(chan BlockHeader,1)
//
//	//var nonce int64 = 0
//	//var hashByte [32]byte
//	//var ts Transaction
//
//	//这里可以定时发送，每个新连接的用户都有机会挖矿
//	wsend.SendBlockHeaderToUser(*p.BlockHeader)
//
//	//注释后可以禁止该节点挖矿
//	//go asyncMine(p)
//
//	//TODO 启动一个计时器来检测当前是否已经出块,每秒检测一次
//	//<- MineFlag
//	ticker1 := time.NewTicker(500000 * time.Microsecond)
//	func(t *time.Ticker) {
//		for {
//			<-t.C
//			if MineFlag == true {
//				MineFlag = false
//				break
//			}
//		}
//	}(ticker1)
//
//	//返回应该需要返回一个随机数、一个幻方、区块hash、用户交易
//	return MineReturnBH, nil
//}

//异步挖矿(可以选择是否启动云计算节点挖矿)
//func asyncMine(p *proofOfWork) {
//	//先给自己添加奖励等
//	publicKeyHash := getPublicKeyHashFromAddress(ThisNodeAddr)
//	txo := TXOutput{TokenRewardNum, publicKeyHash}
//	ts := Transaction{nil, nil, []TXOutput{txo}}
//	ts.hash()
//	p.BlockHeader.TransactionToUser = ts
//
//	var nonce int64 = 0
//	var hashByte [32]byte
//	var hashInt big.Int
//	log.Info("准备挖矿...")
//
//	//开启一个计数器,每隔五秒打印一下当前挖矿,用来直观展现挖矿情况
//	times := 0
//	ticker1 := time.NewTicker(5 * time.Second)
//	go func(t *time.Ticker) {
//		for {
//			<-t.C
//			times += 5
//			log.Infof("正在挖矿,挖矿区块高度为%d,已经运行%ds,nonce值:%d,当前hash:%x", p.Height, times, nonce, hashByte)
//		}
//	}(ticker1)
//
//	for nonce < maxInt {
//		//检测网络上其他节点是否已经挖出了区块
//		if p.Height <= NewestBlockHeight {
//			//结束计数器
//			ticker1.Stop()
//			MineReturnBH = BlockHeader{}
//			GlobalErr = errors.New("检测到当前节点已接收到最新区块，所以终止此块的挖矿操作")
//			MineFlag = true
//			//MineFlag <- true
//			return
//		}
//
//		//TODO 假设挖出了随机幻方的第一个数字，随机幻方是有约束的
//		randomMatrix := RandomMatrix{[10][10]int64{
//			[10]int64{nonce, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//			[10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
//		}}
//
//		data := p.jointData(randomMatrix)
//
//		hashByte = sha256.Sum256(data)
//		//将hash值转换为大数字
//		hashInt.SetBytes(hashByte[:])
//		//如果hash后的data值小于设置的挖矿难度大数字,则代表挖矿成功!
//		if hashInt.Cmp(p.Target) == -1 {
//			//TODO
//			p.RandomMatrix = randomMatrix
//			break
//		} else {
//			//原来是nonce++，现在是产生随机数
//			bigInt, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
//			if err != nil {
//				log.Panic("随机数错误:", err)
//			}
//			nonce = bigInt.Int64()
//		}
//	}
//	//结束计数器
//	ticker1.Stop()
//
//	p.BlockHeader.Hash = hashByte[:]
//	log.Infof("本节点已成功挖到区块!!!,高度为:%d,nonce值为:%d,区块hash为: %x", p.Height, nonce, hashByte)
//
//	MineReturnBH = *p.BlockHeader
//	GlobalErr = nil
//	MineFlag = true
//	//MineFlag <- true
//	return
//}

//检验区块是否有效
func (p *proofOfWork) Verify() bool {
	//TODO 验证MerkelRootWHash
	ts := p.BBlockHeader.TransactionToUser
	merkelRootWHash := sha256.Sum256(bytes.Join([][]byte{ts.getTransBytes(), p.BBlockHeader.MerkelRootHash}, []byte("")))
	//此处100 为W重
	for i := 0; i < WNum; i++ {
		merkelRootWHash = sha256.Sum256(merkelRootWHash[:])
	}
	var m1 big.Int
	var m2 big.Int
	m1.SetBytes(merkelRootWHash[:])
	m2.SetBytes(p.BBlockHeader.MerkelRootWHash[:])
	if m1.Cmp(&m2) != 0 {
		log.Debug(m1.String())
		log.Debug(m2.String())
		return false
	}
	//TODO 验证RandomMatrix是否符合要求

	//TODO 验证用户奖励金额
	if p.BBlockHeader.TransactionToUser.Vout[0].Value != TokenRewardNum {
		return false
	}

	target := big.NewInt(1)
	target.Lsh(target, 256-TargetBits)
	data := p.jointData(p.BBlockHeader.RandomMatrix)
	hash := sha256.Sum256(data)
	var hashInt big.Int
	hashInt.SetBytes(hash[:])
	if hashInt.Cmp(target) == -1 {
		return true
	}
	return false
}

// TODO 将 上一区块hash、数据、时间戳、难度位数、随机数 拼接成字节数组
func (p *proofOfWork) jointData(randomMatrix RandomMatrix) (data []byte) {
	preHash := p.BBlockHeader.PreHash
	preRandomMatrixByte := RandomMatrixToBytes(p.BBlockHeader.PreRandomMatrix)
	merkelRootHash := generateMerkelRoot(p.Transactions) //p.BBlockHeader.MerkelRootHash
	merkelRootWHash := p.BBlockHeader.MerkelRootWHash
	merkelRootWSignature := p.BBlockHeader.MerkelRootWSignature
	cAByte := CAToBytes(p.BBlockHeader.CA)

	transactionToUserByte := p.BBlockHeader.TransactionToUser.getTransBytes()

	timeStampByte := util.Int64ToBytes(p.BBlockHeader.TimeStamp)
	heightByte := util.Int64ToBytes(int64(p.BBlockHeader.Height))
	randomMatrixByte := RandomMatrixToBytes(randomMatrix)
	targetBitsByte := util.Int64ToBytes(int64(TargetBits))

	data = bytes.Join([][]byte{
		preHash,
		preRandomMatrixByte,
		merkelRootHash,
		merkelRootWHash,
		merkelRootWSignature,
		cAByte,
		transactionToUserByte,
		timeStampByte,
		heightByte,
		randomMatrixByte,
		targetBitsByte},
		[]byte(""))
	return data
}
