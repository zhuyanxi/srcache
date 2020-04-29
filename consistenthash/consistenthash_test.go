package consistenthash

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestConsistentHash(t *testing.T) {
	// hash256 := sha256.Sum256([]byte("123"))
	// fmt.Println(hash256)
	// hash256Str := hex.EncodeToString(hash256[:])
	// fmt.Println(hash256Str)
	// fmt.Printf("%x\n", hash256[:])
	// hash256Decode, _ := hex.DecodeString(hash256Str)
	// fmt.Println(hash256Decode)
	// hashBig := binary.BigEndian.Uint32(hash256[:])
	// fmt.Println(hashBig)
	// hashSmall := binary.LittleEndian.Uint32(hash256[:])
	// fmt.Println(hashSmall)

	// hash := NewConsistentHash(3, func(data []byte) int {
	// 	return int(crc32.ChecksumIEEE(data))
	// })
	// hash := NewConsistentHash(4, func(data []byte) int {
	// 	h1 := sha256.Sum256(data)
	// 	// hash3 := crc32.ChecksumIEEE(data)
	// 	hash32 := binary.BigEndian.Uint32(h1[:])
	// 	return int(hash32)
	// })
	hash := NewConsistentHash(64, nil)

	nodes := []Node{
		{Addr: "localhost:3001"},
		{Addr: "localhost:3002"},
		{Addr: "localhost:3003"},
	}

	hash.Add(nodes...)

	m := hash.ring
	for _, k := range hash.sortedKeys {
		logrus.Infof("k is %d, v is %s\n", k, m[k].Addr)
	}

	ipStr := "151.176.74.57,79.182.220.16,40.118.139.44,136.59.35.199,97.213.58.164,23.132.99.218,77.208.180.89,60.124.204.66,152.178.71.149,207.72.57.63,188.245.103.75,98.23.154.47,52.226.51.118,177.126.189.18,142.98.181.74,46.255.137.244,103.29.235.36,130.209.253.33,187.199.146.214,20.92.1.80,239.77.206.158,130.32.195.89,17.203.166.210,98.26.106.186,188.90.43.182,62.67.28.22,33.144.98.81,25.139.244.151,6.171.74.28,94.154.8.6,251.101.205.185,193.179.30.22,72.149.208.112,212.161.202.90,151.144.119.244,81.206.131.235,182.42.192.47,123.245.236.92,211.97.223.54,35.55.15.88,191.211.240.76,143.246.161.156,78.238.103.40,64.83.252.100,214.166.38.109,116.125.192.145,50.73.178.50,149.229.25.111,115.255.147.206,127.232.176.70,238.160.202.0,194.24.189.135,74.13.127.23,141.13.60.48,229.159.38.136,128.148.240.41,215.196.242.70,51.248.79.233,113.146.164.28,30.178.62.198,58.191.241.144,82.232.184.94,27.65.192.203,126.187.254.67,4.34.18.245,63.124.9.59,130.49.70.178,165.10.167.34,43.17.177.136,10.251.186.129,49.198.100.160,44.83.84.216,149.59.207.237,109.208.246.223,183.132.106.203,140.48.88.94,120.38.139.201,188.115.64.49,76.56.138.244,125.143.106.66,219.206.207.159,91.127.134.145,255.201.70.80,99.192.74.107,35.74.50.117,68.77.124.75,111.165.94.223,86.119.158.255,27.186.35.185,165.183.1.112,103.44.61.252,83.116.233.174,20.44.147.155,189.23.67.60,186.161.130.150,114.224.139.201,173.58.12.143,219.40.66.75,86.151.92.246,236.57.216.213"
	ips := strings.Split(ipStr, ",")
	// for i := 0; i < 10000; i++ {
	// 	a := faker.IPv4()
	// 	ips = append(ips, a)
	// }

	var num1, num2, num3 int
	for _, ip := range ips {
		nn := hash.Get(ip)
		if strings.Contains(nn.Addr, "3001") {
			num1++
		} else if strings.Contains(nn.Addr, "3002") {
			num2++
		} else {
			num3++
		}
		logrus.Infof("Value %s will be stored on node %s", ip, nn.Addr)
	}
	logrus.Infof("Num1: %d; Num2: %d; Num3: %d", num1, num2, num3)

	// value1 := "10.192.168.10"
	// n1 := hash.Get(value1)
	// logrus.Infof("Value %s will be stored on node %s", value1, n1.Addr)

	// nodes2 := []Node{
	// 	{Addr: "8"},
	// }
	// hash.Add(nodes2...)
	// m2 := hash.ring
	// for _, k := range hash.sortedKeys {
	// 	logrus.Infof("k is %d, v is %s\n", k, m2[k].Addr)
	// }
	// value2 := "10.192.168.10"
	// n2 := hash.Get(value2)
	// logrus.Infof("Value %s will be stored on node %s", value2, n2.Addr)
}
