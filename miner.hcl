donate = 1

cpu "Intel i5" {
  threads = 1
  coin = "demo"
}

gpu "RX 560" {
  index = 0
  coin = "eth"
}

coin "demo" {

}

coin "xmr" {
  pool {
    url = "stratum+tcp://xmr.poolmining.org:3032",
    user = "46DTAEGoGgc575EK7rLmPZFgbXTXjNzqrT4fjtCxBFZSQr5ScJFHyEScZ8WaPCEsedEFFLma6tpLwdCuyqe6UYpzK1h3TBr",
    pass = "x",
  }
}

coin "eth" {
  pool {
    url = "stratum+tcp://eth.poolmining.org:3072"
    user = "0x25ae2cbddE36CfC9D959a4d1f76964EaE7517748"
  }
}

coin "ubq" {
  pool {
    url = "stratum+tcp://ubiq.coinminer.space:8008"
    user = "0xE0F39621764F2540cd2bC3017DA041B1E4eEDCc2"
    pass = "x"
    email = "asdasd"
  }
}
