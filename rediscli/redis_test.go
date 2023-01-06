package rediscli

import (
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestRedisCli(t *testing.T) {
	convey.Convey("TestRedisCliHttp", t, func() {
		var fac = &RedisFactory{}
		cli, err := fac.CreateRedisCli("", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("create conn pool err ", err)
			return
		}
		err = cli.Set("abc", "hahaha", 5)
		convey.So(err, convey.ShouldBeNil)
		buf, err := cli.Get("abc")
		convey.So(err, convey.ShouldBeNil)
		s := string(buf)
		convey.So(s, convey.ShouldEqual, "hahaha")
	})

}

func TestRedisSecKill(t *testing.T) {
	convey.Convey("TestRedisCliHttp", t, func() {
		var fac = &RedisFactory{}
		cli, err := fac.CreateRedisCli("", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("create conn pool err ", err)
			return
		}
		sk := CreateSecKillHandler(cli)
		err = sk.BuildData("pika", 100)
		fmt.Println(err)
		err = sk.Aquire("pika", 2)
		fmt.Println(err)
		//s := string(buf)
		//convey.So(s, convey.ShouldEqual, "hahaha")

	})
}

func TestRedisSecKillAquire(t *testing.T) {
	convey.Convey("TestRedisCliHttp", t, func() {
		var fac = &RedisFactory{}
		cli, err := fac.CreateRedisCli("", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("create conn pool err ", err)
			return
		}
		sk := CreateSecKillHandler(cli)
		fmt.Println(err)
		err = sk.Aquire("pika", 95)
		fmt.Println(err)
		//s := string(buf)
		//convey.So(s, convey.ShouldEqual, "hahaha")

	})
}

func TestIncrease(t *testing.T) {
	key := "cataabb"
	var fac = &RedisFactory{}
	cli, err := fac.CreateRedisCli("", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	nums, err := cli.IncrBy(key, 20)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("nums ", nums)
}

func TestDesc(t *testing.T) {
	key := "cataabbaaa"
	var fac = &RedisFactory{}
	cli, err := fac.CreateRedisCli("", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	nums, err := cli.DescBy(key, 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("nums ", nums)
}

func TestSetNX(t *testing.T) {
	key := "TestKeyJing3"
	value := "this is value3"
	var fac = &RedisFactory{}
	cli, err := fac.CreateRedisCli("123", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cli.SetNX(key, value, 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cli.SetNX(key, value, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("nums ", err)
}

func TestSetNX2(t *testing.T) {
	key := "TestKeyJing4"
	value := "8"
	var fac = &RedisFactory{}
	cli, err := fac.CreateRedisCli("", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("CreateRedisCli err:", err)
		return
	}
	err = cli.SetNX(key, value, 60)
	if err != nil {
		fmt.Println("SetNX err:", err)
		return
	}
	fmt.Println(value)
	//加
	by, err := cli.IncrBy(key, 1)
	fmt.Println("加1等于")
	fmt.Println(by, err)
	//减
	deby, err := cli.DescBy(key, 2)
	fmt.Println("减2等于")
	fmt.Println(deby, err)
}

func TestPushList(t *testing.T) {
	key := "list1"
	value := "8"
	var fac = &RedisFactory{}
	cli, err := fac.CreateRedisCli("", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("CreateRedisCli err:", err)
		return
	}
	num, err := cli.RPushAndSetExpire(key, value, 10)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(num)
}

func TestLpop(t *testing.T) {
	key := "list1"
	var fac = &RedisFactory{}
	cli, err := fac.CreateRedisCli("", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("CreateRedisCli err:", err)
		return
	}
	temp, err := cli.LPop(key)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(temp)
}

func TestLRange(t *testing.T) {
	key := "list1"
	var fac = &RedisFactory{}
	cli, err := fac.CreateRedisCli("", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("CreateRedisCli err:", err)
		return
	}
	temp, err := cli.LRange(key, 0, 2)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(temp)
}
