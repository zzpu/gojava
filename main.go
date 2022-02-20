package main

//java serialize & deserialize

import (
	"bytes"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	log "github.com/corgi-kx/logcustom"
	"github.com/go-redis/redis"
	"os"

	//jserial "github.com/jkeys089/jserial"
	"time"
)

var (
	f *os.File
)

type SonEntity struct {
	SonName string
}

type BaseEntity struct {
	BaseName string
}

type UserEntity struct {
	BaseEntity
	Id       int64
	UserName string
	UserSex  string
	Gifts    ArrayList
	Son      SonEntity
}
type ArrayList struct {
	Eles []interface{}
}

func main() {
	f, _ = os.Create("./log.txt") //创建文件

	defer f.Close()

	redisPool := redis.NewClient(&redis.Options{
		Addr:         "47.115.166.195:6379",
		Password:     "",
		DB:           0,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     0,
		PoolTimeout:  30 * time.Second,
		MaxRetries:   2,
		IdleTimeout:  5 * time.Minute,
	})
	aa, _ := redisPool.Get("cccccc").Bytes()

	for _, v := range aa {
		fmt.Printf("%0x ", v)
	}
	log.SetLogDiscardLevel(log.Levelwarn)
	var f = bytes.NewBuffer(aa)
	//jo := &JavaArrayList{}
	jo, err := DeserializeStream(f)
	if err != nil {
		log.Errorf("err=%v", err)
	}

	j, ok := jo.(*JavaTcObject)
	if ok {
		aaa, _ := json.Marshal(j.Classes[1].RwDatas)
		log.Debugf("----------->%v", string(aaa))
		ccc, _ := json.Marshal(j.JsonData)

		var ttt UserEntity
		json.Unmarshal(ccc, &ttt)
		log.Debugf("----------->%v", string(ccc))

		dddd, _ := json.Marshal(ttt)

		log.Debugf("----------->%v", string(dddd))
		//rw := j.Classes[0].RwDatas
		//jj,ok := rw[0].(*JavaArrayList)
		//
		//if ok{
		//
		//	ccc,_ := json.Marshal(jj.Eles)
		//
		//	var ttt Test
		//	json.Unmarshal(ccc,&ttt)
		//
		//	fmt.Println(string(ccc))
		//}

	}

	//ReadUint16(f)
	//var refs  = make([]*JavaReferenceObject,10)

	//err := jo.Deserialize(f,refs)
	//if err != nil{
	//	fmt.Println("err=",err)
	//}

	return

	//
	//jo := NewJavaTcObject(1)
	//clz := NewJavaTcClassDesc("com.gauzz.lession2.model.Test", 1, 0x02)
	//jfa := NewJavaField(TC_PRIM_INTEGER, "Value", 222)
	//jfb := NewJavaField(TC_OBJ_OBJECT, "Next", "@846s4d6f4654f4a5s4f64asdf45ds")
	//jfb.FieldObjectClassName = "java.lang.String"
	//clz.AddField(jfa)
	//clz.AddField(jfb)
	////clz.SortFields()
	//
	//jo.AddClassDesc(clz)
	//var f bytes.Buffer
	//
	//_ = SerializeJavaEntity(&f, jo)
	//redisPool.Set("bbbbbb",f.String(),0)
	//fmt.Println(f.String())
}

func Log(s string) {

	f.WriteString(s) //写入文件(字节数组)

	f.Sync()

}
