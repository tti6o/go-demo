package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//审批存盘记录
type ApprovalRecord struct {
	SiteId         string           `json:"site_id,omitempty" bson:"site_id"`                 // 站点id
	HistoryRecords *[]HistoryRecord `json:"history_records,omitempty" bson:"history_records"` // 审批历史进度
	CreateTime     int64            `json:"create_time,omitempty" bson:"create_time"`
	UpdateTime     int64            `json:"update_time,omitempty" bson:"update_time"`
}

//审批历史进度
type HistoryRecord struct {
	Activity       string     `bson:"activity,omitempty"`         // 审批节点
	ProcessInstId  string     `bson:"process_inst_id,omitempty"`  // 流程实例名称
	Status         int        `bson:"activity_status,omitempty"`  // 该节点的审批状态
	Handlers       []string   `bson:"handlers,omitempty"`         // 可以审批的人列表
	CreateTime     *time.Time `bson:"create_time,omitempty"`      // 该节点创建单据的时间
	ActivityItemId string     `bson:"activity_item_id,omitempty"` // 单据id
	SubmitHandler  string     `bson:"submit_handler,omitempty"`   // 执行审批的人
	SubmitAction   string     ` bson:"submit_action,omitempty"`   // 审批动作
	SubmitOpinion  string     `bson:"submit_opinion,omitempty"`   // 审批意见
	SubmitTime     *time.Time `bson:"submit_time,omitempty"`      // 审批时间
}

type Class struct {
	Name     interface{} `json:"name"`
	Number   int         `json:"number"`
	Students *[]Student  `json:"students"`
}

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return
	}
	//class := Class{
	//	"class2",
	//	2,
	//	&[]Student{
	//		{
	//			"zhangsan",
	//			1},
	//		{
	//			"lisi",
	//			2,
	//		},
	//	},
	//}
	//insertOneResult, err := client.Database("easp-ass-test").Collection("approval_record").InsertOne(ctx, &class)
	//fmt.Println(insertOneResult, err)

	//query := bson.M{"name": "class2"}
	//update := bson.M{"$addToSet": bson.M{"students": bson.M{"name": "wangwu", "age": 3}}}
	//updateResult, err := client.Database("easp-ass-test").Collection("approval_record").UpdateOne(ctx, query, update)
	//fmt.Println(updateResult, err)

	// 创建 审批记录
	now := time.Now()
	approvalRecord := ApprovalRecord{
		SiteId: "111",
		HistoryRecords: &[]HistoryRecord{
			{
				Activity:      "leader 审批",
				ProcessInstId: "222",
				Status:        1, //待审批
				Handlers:      []string{"nickyluo"},
				CreateTime:    &now,
			},
		},
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	insertOneResult, err := client.Database("easp-ass-test").Collection("approval_record").InsertOne(ctx, &approvalRecord)
	fmt.Printf("insertOneResult:%v,err:%v \n", insertOneResult, err)

	//query := bson.M{"process_inst_id": "222", "handler": "nickyluo"}
	//update := bson.M{"$addToSet": bson.M{"history_records": bson.M{"submit_time": time.Now(), "step": 2}}}
	//updateResult, err := client.Database("easp-ass-test").Collection("approval_record").UpdateOne(ctx, query, update)
	//fmt.Println(updateResult, err)

	// 数组增加元素
	//modify := HistoryRecord{
	//	SubmitTime: time.Now(),
	//}
	//jsonStr, err := json.Marshal(modify)
	//if err != nil {
	//	return
	//}
	//modifyMap := map[string]interface{}{}
	//json.Unmarshal(jsonStr, &modifyMap)
	//query := bson.M{"handler": "nickyluo", "process_inst_id": "222"}
	//update := bson.M{"$addToSet": modifyMap}
	//fmt.Printf("query = %v,update=%v \n", query, update)
	//updateResult, err := client.Database("easp-ass-test").Collection("approval_record").UpdateOne(ctx, query, update)
	//fmt.Printf("updateResult1 =%+v,err=%v \n", updateResult, err)

	//更新数组里某个元素整体

	modify := &HistoryRecord{
		SubmitTime: &now,
	}
	bsonStr, err := bson.Marshal(modify)
	if err != nil {
		return
	}
	modifyMap := map[string]interface{}{}
	bson.Unmarshal(bsonStr, &modifyMap)
	bMap := bson.M{}
	for key, value := range modifyMap {
		bKey := "history_records.$." + key
		bMap[bKey] = value
	}
	query := bson.M{
		"site_id": "111",
		"history_records": bson.M{
			"$elemMatch": bson.M{
				"process_inst_id": "222",
			},
		},
	}
	update := bson.M{
		"$set": bMap,
	}
	fmt.Printf("query = %+v,update=%+v \n", query, update)
	updateResult, err := client.Database("easp-ass-test").Collection("approval_record").UpdateOne(ctx, query, update)
	fmt.Printf("updateResult1 =%+v,err=%v \n", updateResult, err)

	//更新数组里某个元素的单个字段

	//query := bson.M{
	//	"site_id": "111",
	//	"history_records": bson.M{
	//		"$elemMatch": bson.M{
	//			"process_inst_id": "222",
	//		},
	//	},
	//}
	//updateFiled := fmt.Sprintf("history_records.$.%s", "submit_opinion")
	//update := bson.M{
	//	"$set": bson.M{
	//		updateFiled: "1234567",
	//	},
	//}
	//fmt.Printf("query = %v,update=%v \n", query, update)
	//updateResult, err := client.Database("easp-ass-test").Collection("approval_record").UpdateOne(ctx, query, update)
	//fmt.Printf("updateResult2 =%+v,err=%v \n", updateResult, err)
}
