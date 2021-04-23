package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	blockchain := Blockchain{}
	blockchain.initialize()

	nodeAddress := "randomAddress"

	r := gin.Default()

	//  localhost:PORT/mineBlock
	r.GET("/mineBlock", func(c *gin.Context) {
		prevBlock := blockchain.getPreviousBlock()
		prevProof, _ := strconv.Atoi(prevBlock["proof"])
		proof := blockchain.proofOfWork(prevProof)
		prevHash := blockchain.hash(prevBlock)
		//Reward for mining
		blockchain.addTransaction(nodeAddress, "Anup", 50)
		block := blockchain.createBlock(proof, prevHash)
		c.JSON(200, gin.H{
			"message":      "Successfully mined a block",
			"index":        block["index"],
			"timestamp":    block["timestamp"],
			"proof":        block["proof"],
			"previousHash": block["previousHash"],
			"transactions": block["transactions"],
		})
	})

	//  localhost:PORT/getChain
	r.GET("/getChain", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"chain":               blockchain.chain,
			"pendingTransactions": blockchain.transactions,
			"Length":              len(blockchain.chain),
		})
	})

	//  localhost:PORT/isChainValid
	r.GET("/isChainValid", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"isValid": blockchain.isChainValid(blockchain.chain),
		})
	})

	//  localhost:PORT/replaceChain
	r.GET("/replaceChain", func(c *gin.Context) {
		isChainReplaced := blockchain.replaceChain()
		message := ""
		if isChainReplaced {
			message = "Chain was successfully replaced by the longet one"
		} else {
			message = "Chain was already the longest one. Did not replace"
		}
		c.JSON(200, gin.H{
			"Message": message,
		})
	})

	// localhost:PORT/connectNode
	r.POST("/connectNode", func(c *gin.Context) {

		nodeAddress := c.Query("nodeAddress")
		blockchain.addNode(nodeAddress)
		c.JSON(200, gin.H{
			"AddedNodeAddress": nodeAddress,
			"Message":          "Added node to the blockchain",
			"NodeCount":        len(blockchain.nodes),
			"Nodes":            blockchain.nodes,
		})
	})
	// localhost:8080/addTransaction
	r.POST("/addTransaction", func(c *gin.Context) {
		fmt.Println("addTransaction Method")
		sender := c.Query("sender")
		receiver := c.Query("receiver")
		amount, _ := strconv.ParseFloat(c.Query("amount"), 64)
		index := blockchain.addTransaction(sender, receiver, amount)
		c.JSON(200, gin.H{
			"BlockIndex": index,
			"Message":    "Transaction will be added to mentioned block",
		})
	})

	r.Run(":8000") // listen and serve on 0.0.0.0:8080
}
