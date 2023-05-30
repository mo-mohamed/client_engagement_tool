package groupBroadcastProcessor

import (
	dbc "customer_engagement/data_store/config"
	"customer_engagement/queue"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type GroupMessageProcessor struct {
	bSize int
}

type ProfileMessage struct {
	Mdn     string
	GroupId int
	Id      int
	Body    string
}

func NewGroupMessageProcessor(batchSize int) *GroupMessageProcessor {
	return &GroupMessageProcessor{
		bSize: batchSize,
	}
}

func (g *GroupMessageProcessor) Process(message *queue.Message) error {
	groupId, _ := strconv.Atoi(message.Attributes["GroupId"].Value)
	dataEnqueued, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", message.Attributes["DateEnqueued"].Value)
	numProfiles := getNumberOfProfiles(groupId, dataEnqueued)

	numOfBatches := numProfiles / g.bSize
	startBatch := 0

	for i := 0; i < numOfBatches; i++ {
		// fmt.Println("STARTING A BATCH")
		profilesData := getProfilesBatch(groupId, g.bSize, startBatch, dataEnqueued)

		var wg sync.WaitGroup
		for _, p := range profilesData {
			p.Body = message.Body
			wg.Add(1)
			go process(&p, &wg)
		}
		wg.Wait()
		startBatch = startBatch + g.bSize
	}

	return nil
}

func process(p *ProfileMessage, wg *sync.WaitGroup) {
	// TODO Send the message here
	fmt.Println("Profile processed")
	fmt.Printf("%+v \n", p)
	wg.Done()
}

func getNumberOfProfiles(groupId int, dateEnqueued time.Time) int {
	query := `
		SELECT COUNT(*) FROM group_profile gp
		INNER JOIN profile p on gp.profile_id = p.id
		WHERE gp.group_id = ? AND gp.created_at <= ?;
	`
	var count int
	dbc.DB.Raw(query, groupId, dateEnqueued).Scan(&count)
	return count
}

func getProfilesBatch(groupId, limit, offest int, dateEnqueued time.Time) []ProfileMessage {
	query := `
		SELECT gp.group_id as GroupId,
			   p.id as Id,
			   p.mdn as Mdn
		FROM group_profile gp 
		INNER JOIN profile p on gp.profile_id = p.id
		WHERE gp.group_id = ? AND gp.created_at <= ? limit ?, ?;
	`

	var profiles []ProfileMessage
	dbc.DB.Raw(query, groupId, dateEnqueued, offest, limit).Scan(&profiles)
	return profiles
}
