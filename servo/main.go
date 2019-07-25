package main

import (
	"../common"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Lmicroseconds)


	// Servo
	ss := NewHexaServo()
	go ss.Control()

	mesChan := make(chan *common.Message, 100)
	deferedMessage := make([]*common.Message, 0)

	go func() {
		currentSequence := 0

		for {
			select {
			case mes := <-mesChan:
				// 初期はリセットまで読み飛ばし
				if currentSequence == 0 && mes.MessageType != common.MessageTypeReset {
					log.Println("waiting reset", mes.Sequence)
					continue
				}
				// 通り過ぎたシークエンスは捨てる
				if currentSequence > mes.Sequence {
					log.Println("old", mes.Sequence)
					continue
				}
				// Resetが来たらとにかくシークエンスを揃える
				if mes.MessageType == common.MessageTypeReset {
					currentSequence = mes.Sequence
					ss.Input(mes.Arms, 3*time.Second)
					for _, mes := range deferedMessage {
						mesChan <- mes
					}
					log.Println("reset", mes.Sequence)
					continue
				}

				// sequenceが想定外ならあとで
				if currentSequence+1 != mes.Sequence {
					deferedMessage = append(deferedMessage, mes)
					log.Println("defer", mes.Sequence)
					continue
				}

				// 処理
				switch mes.MessageType {
				case common.MessageTypeArms:
					ss.Input(mes.Arms, 1*time.Second)
				case common.MessageTypeRotate:
					dd := [][]float64{
						{-10, -10, 0, 0, 0, 0, 0, 0}, // 初期状態

						{10, -10, -30, 0, -30, 0, -30, 0},  // 1
						{-10, -10, -30, 0, -30, 0, -30, 0}, // 1
						{-10, 10, 0, -30, 0, -30, 0, -30},  // 2
						{-10, -10, 0, -30, 0, -30, 0, -30}, // 2

						{10, -10, -30, 0, -30, 0, -30, 0},  // 1
						{-10, -10, -30, 0, -30, 0, -30, 0}, // 1
						{-10, 10, 0, -30, 0, -30, 0, -30},  // 2
						{-10, -10, 0, -30, 0, -30, 0, -30}, // 2

						{10, -10, -30, 0, -30, 0, -30, 0},  // 1
						{-10, -10, -30, 0, -30, 0, -30, 0}, // 1
						{-10, 10, 0, -30, 0, -30, 0, -30},  // 2
						{-10, -10, 0, -30, 0, -30, 0, -30}, // 2

						{10, -10, -30, 0, -30, 0, -30, 0},  // 1
						{-10, -10, -30, 0, -30, 0, -30, 0}, // 1
						{-10, 10, 0, -30, 0, -30, 0, -30},  // 2
						{-10, -10, 0, -30, 0, -30, 0, -30}, // 2

						{-10, -10, 0, 0, 0, 0, 0, 0}, // 初期状態
					}
					for _, d := range dd {
						arms := []*common.Arm{
							{
								Degrees: []float64{d[2], d[1]},
							},
							{
								Degrees: []float64{d[3], d[0]},
							},
							{
								Degrees: []float64{d[4], d[1]},
							},
							{
								Degrees: []float64{d[5], d[0]},
							},
							{
								Degrees: []float64{d[6], d[1]},
							},
							{
								Degrees: []float64{d[7], d[0]},
							},
						}
						ss.Input(arms, 1*time.Second)
					}

				case common.MessageTypeWalk:
					dd := [][]float64{
						{-10, -10, -30, 0, 30, -30, 0, 30}, // 初期状態

						{10, -10, -30, 30, 30, -60, 0, 0},  // 足上げ1で前進1
						{-10, -10, -30, 30, 30, -60, 0, 0}, // 足下げ1で前進1
						{-10, 10, -30, 30, 30, -60, 0, 0},  // 足上げ2で前進1
						{-10, 10, 0, 0, 60, -30, -30, 30},  // 足上げ2で前進2
						{-10, -10, 0, 0, 60, -30, -30, 30}, // 足下げで前進2
						{10, -10, 0, 0, 60, -30, -30, 30},  // 足上げ1で前進2

						{10, -10, -30, 30, 30, -60, 0, 0},  // 足上げ1で前進1
						{-10, -10, -30, 30, 30, -60, 0, 0}, // 足下げ1で前進1
						{-10, 10, -30, 30, 30, -60, 0, 0},  // 足上げ2で前進1
						{-10, 10, 0, 0, 60, -30, -30, 30},  // 足上げ2で前進2
						{-10, -10, 0, 0, 60, -30, -30, 30}, // 足下げで前進2
						{10, -10, 0, 0, 60, -30, -30, 30},  // 足上げ1で前進2

						{10, -10, -30, 30, 30, -60, 0, 0},  // 足上げ1で前進1
						{-10, -10, -30, 30, 30, -60, 0, 0}, // 足下げ1で前進1
						{-10, 10, -30, 30, 30, -60, 0, 0},  // 足上げ2で前進1
						{-10, 10, 0, 0, 60, -30, -30, 30},  // 足上げ2で前進2
						{-10, -10, 0, 0, 60, -30, -30, 30}, // 足下げで前進2
						{10, -10, 0, 0, 60, -30, -30, 30},  // 足上げ1で前進2

						{10, -10, -30, 30, 30, -60, 0, 0},  // 足上げ1で前進1
						{-10, -10, -30, 30, 30, -60, 0, 0}, // 足下げ1で前進1
						{-10, 10, -30, 30, 30, -60, 0, 0},  // 足上げ2で前進1
						{-10, 10, 0, 0, 60, -30, -30, 30},  // 足上げ2で前進2
						{-10, -10, 0, 0, 60, -30, -30, 30}, // 足下げで前進2
						{10, -10, 0, 0, 60, -30, -30, 30},  // 足上げ1で前進2

						{-10, -10, -30, 0, 30, -30, 0, 30}, // 初期状態
					}
					for _, d := range dd {
						arms := []*common.Arm{
							{
								Degrees: []float64{d[2], d[1]},
							},
							{
								Degrees: []float64{d[3], d[0]},
							},
							{
								Degrees: []float64{d[4], d[1]},
							},
							{
								Degrees: []float64{d[5], d[0]},
							},
							{
								Degrees: []float64{d[6], d[1]},
							},
							{
								Degrees: []float64{d[7], d[0]},
							},
						}
						ss.Input(arms, 1*time.Second)
					}
				case common.MessageTypeDemo:
					// do nothing
				default:
					// do nothing
				}

				// Deferedなのを戻す
				for _, mes := range deferedMessage {
					mesChan <- mes
				}
				currentSequence = mes.Sequence
			}
		}
	}()

	// heartbeat
	go func(){
		ctx := context.Background()

		client, err := pubsub.NewClient(ctx, "hexapod", option.WithCredentialsFile("cred.json"))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
			return
		}
		t := client.Topic("heartbeat")
		if err != nil {
			log.Fatalf("Failed to create channel: %v", err)
			return
		}


		ticker := time.NewTicker(5* time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				t.Publish(ctx, &pubsub.Message{
					Data:        []byte(time.Now().String()),
				})
			}
		}
	}()

	// PubSub
	ctx := context.Background()
	projectID := "hexapod"
	subscription := "servo"
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("cred.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	sub := client.Subscription(subscription)
	_ = sub.Receive(ctx, func(c context.Context, m *pubsub.Message) {
		mes, err := common.UnmarshalMessage(m.Data)
		if err != nil {
			return
		}

		mesChan <- mes
		fmt.Println(mes.Sequence, time.Now().Sub(mes.Now))

		m.Ack()
	})
}

func mixArms(base, future []*common.Arm, ratio float64) []*common.Arm {
	if ratio > 1.0 {
		ratio = 1.0
	}

	if base == nil {
		return future
	}

	ret := make([]*common.Arm, 6)
	for i := range ret {
		ret[i] = &common.Arm{
			Degrees: make([]float64, 2),
		}
	}

	for i := range base {
		for j := range base[i].Degrees {
			diff := future[i].Degrees[j] - base[i].Degrees[j]
			ret[i].Degrees[j] = base[i].Degrees[j] + ratio*diff
		}
	}

	return ret
}

type HexaServo struct {
	servos  []*Servo
	input   chan []*common.Arm
	current []*common.Arm
}

func NewHexaServo() *HexaServo {
	d := initializeDriver()
	ss := make([]*Servo, 12)
	for i := range ss {
		ss[i] = NewServo(d, i)
	}

	c := make(chan []*common.Arm, 0)
	return &HexaServo{
		servos: ss,
		input:  c,
	}
}

func (h *HexaServo) Control() {
	for {
		select {
		case arms := <-h.input:
			i := 0
			for ai, a := range arms {
				for di, d := range a.Degrees {
					if h.current == nil || d != h.current[ai].Degrees[di] {
						h.servos[i].setDegree(d)
					}
					i++
				}
			}
		}
	}
}

func (h *HexaServo) Input(val []*common.Arm, duration time.Duration) {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	start := time.Now()
	end := start.Add(duration)

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			h.input <- mixArms(h.current, val, now.Sub(start).Seconds()/end.Sub(start).Seconds())
			if now.After(end) {
				h.current = val
				return
			}
		}
	}

}
