package main

import (
	"../common"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	ctx := context.Background()
	c, err := NewOrderChannel(ctx, "hexapod", "servo")
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
		return
	}

	if len(os.Args) < 2 {
		log.Fatalf("param")
		return
	}

	var sequence int

	data, err := ioutil.ReadFile(`seq.txt`)
	if err != nil {
		log.Fatalf("failed to read seq.txt: %v", err)
		return
	}

	sequence, err = strconv.Atoi(string(data))
	if err != nil {
		log.Fatalf("failed to read seq.txt %v", err)
		return
	}

	switch os.Args[1] {
	case "heartbeat":
		// PubSub
		ctx := context.Background()
		projectID := "hexapod"
		subscription := "heartbeat"
		client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("cred.json"))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		sub := client.Subscription(subscription)
		_ = sub.Receive(ctx, func(c context.Context, m *pubsub.Message) {
			fmt.Println(string(m.Data))
			m.Ack()
		})
	case "walk":
		err = c.Publish(ctx, &common.Message{
			Now:         time.Now(),
			Sequence:    sequence,
			MessageType: common.MessageTypeWalk,
			Arms:        nil,
		})
	case "reset":
		sequence += 100
		err = c.Publish(ctx, &common.Message{
			Now:         time.Now(),
			Sequence:    sequence,
			MessageType: common.MessageTypeReset,
			Arms: []*common.Arm{
				{
					Degrees: []float64{0,0},
				},
				{
					Degrees: []float64{0,0},
				},
				{
					Degrees: []float64{0,0},
				},
				{
					Degrees: []float64{0,0},
				},
				{
					Degrees: []float64{0,0},
				},
				{
					Degrees: []float64{0,0},
				},
			},
		})
	case "rotate":
		err = c.Publish(ctx, &common.Message{
			Now:         time.Now(),
			Sequence:    sequence,
			MessageType: common.MessageTypeRotate,
			Arms:        nil,
		})
	case "arms":
		d := make([]float64, 8)
		d[0], _ = strconv.ParseFloat(os.Args[2], 10)
		d[1], _ = strconv.ParseFloat(os.Args[3], 10)
		d[2], _ = strconv.ParseFloat(os.Args[4], 10)
		d[3], _ = strconv.ParseFloat(os.Args[5], 10)
		d[4], _ = strconv.ParseFloat(os.Args[6], 10)
		d[5], _ = strconv.ParseFloat(os.Args[7], 10)
		d[6], _ = strconv.ParseFloat(os.Args[8], 10)
		d[7], _ = strconv.ParseFloat(os.Args[9], 10)

		err = c.Publish(ctx, &common.Message{
			Now:         time.Now(),
			Sequence:    sequence,
			MessageType: common.MessageTypeArms,
			Arms: []*common.Arm{
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
			},
		})
	default:
		log.Fatalf("param")
	}
	if err != nil {
		log.Fatalf("Failed to publish: %v", err)
		return
	}
	sequence++

	err = ioutil.WriteFile("seq.txt", []byte(fmt.Sprintf("%d", sequence)), 0644)
	if err != nil {
		log.Fatalf("Failed to write seq: %v", err)
		return
	}

	return
}

type OrderChannel struct {
	index int
	topic *pubsub.Topic
}

func NewOrderChannel(ctx context.Context, projectID, topicName string) (*OrderChannel, error) {
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("cred.json"))
	if err != nil {
		return nil, err
	}
	t := client.Topic(topicName)
	return &OrderChannel{
		index: 0,
		topic: t,
	}, nil
}

func (o *OrderChannel) Publish(ctx context.Context, message *common.Message) error {
	res := o.topic.Publish(ctx, &pubsub.Message{
		Data: message.Message(),
	})
	_, err := res.Get(ctx)
	if err != nil {
		return err
	}

	return nil
}
