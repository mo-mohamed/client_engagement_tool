package groupBroadcastProcessor

import (
	dbc "customer_engagement/data_store/config"
	"customer_engagement/queue"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type GroupMessageProcessor struct{}

type profile struct {
	MDN  string
	GID  int
	PID  int
	Body string
}

func (g *GroupMessageProcessor) Process(message *queue.Message) error {

	// Write down group logic here
	gId, _ := strconv.Atoi(message.Attributes["GroupId"].Value)
	dateEn, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", message.Attributes["DateEnqueued"].Value)
	numProfiles := `
			SELECT COUNT(*) FROM group_profile gp JOIN profile p on gp.profile_id = p.id where gp.group_id = ? AND gp.created_at <= ?;
		`
	var count int
	dbc.DB.Raw(numProfiles, gId, dateEn).Scan(&count)
	fmt.Println("number is: ", count)

	batchSize := 2
	numOfBatches := count / batchSize
	startBatch := 0

	for i := 0; i < numOfBatches; i++ {
		fmt.Println("STARTING A BATCH")
		profiles := `
			SELECT gp.group_id as GID, p.id as PID, p.mdn as MDN FROM group_profile gp JOIN profile p on gp.profile_id = p.id where gp.group_id = ? AND gp.created_at <= ? limit ?, ?;
		`
		var profilesData []profile
		dbc.DB.Raw(profiles, gId, dateEn, startBatch, batchSize).Scan(&profilesData)
		fmt.Printf("%+v \n", profilesData)

		var wg sync.WaitGroup
		for _, p := range profilesData {
			p.Body = message.Body
			wg.Add(1)
			go pro(&p, &wg)
		}
		wg.Wait()
		startBatch = startBatch + batchSize
		fmt.Println("BATCH ENDED")
	}

	return nil
}

func pro(p *profile, wg *sync.WaitGroup) {
	fmt.Println("Profile processed")
	fmt.Printf("%+v \n", p)
	wg.Done()
}
