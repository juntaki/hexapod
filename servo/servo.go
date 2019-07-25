package main

import (
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
)

/*
・ＰＷＭサイクル：２０ｍＳ
・制御パルス：０．５ｍｓ～２．４ｍｓ
・制御角：±約９０°（１８０°）
・配線：茶＝ＧＮＤ、赤＝電源［＋］、橙＝制御信号　［ＪＲタイプ］
・トルク：２．５ｋｇｆ・ｃｍ
・動作速度：０．１秒／６０度
・動作電圧：４．８Ｖ
・温度範囲：０℃～５５℃
・外形寸法：２３ｘ１２．２ｘ２７ｍｍ
・重量：９ｇ
*/

const (
	freq = 50.0          // 50 Hz
	ms   = 1000.0 / freq // 20 ms
	min  = 0.5 / ms * 4096
	max  = 2.4 / ms * 4096
	zero = (min + max) / 2          // zero point
	unit = float64(max-min) / 180.0 // 1 degree
)

func initializeDriver() *i2c.PCA9685Driver {
	log.Println("initialize")
	adaptor := raspi.NewAdaptor()
	dri := i2c.NewPCA9685Driver(adaptor)
	log.Println("initialized")
	err := dri.Start()
	if err != nil {
		panic(err)
	}
	log.Println("freq", freq)
	err = dri.SetPWMFreq(freq)
	if err != nil {
		panic(err)
	}
	return dri
}

func NewServo(
	driver *i2c.PCA9685Driver,
	number int,
) *Servo {
	return &Servo{
		driver: driver,
		number: number,
		current: 0,
	}
}

type Servo struct {
	driver *i2c.PCA9685Driver
	number int
	current uint16
}

func (s *Servo) setDegree(deg float64) {
	off := zero + (unit * deg)
	if s.current > 0 && s.current == uint16(off) {
		return
	}
	log.Printf("servo(%d) : degree = %f, off = %d\n", s.number, deg, uint16(off))
	s.current = uint16(off)
	s.driver.SetPWM(s.number, uint16(0), uint16(off))
	return
}
