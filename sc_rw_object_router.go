package main

import (
	log "github.com/corgi-kx/logcustom"
	"io"
)
import "fmt"

//负责将合适的 ScFlag为 SC_RW_FLAG的路由至自定义的各个JavaSerializer 实现类
//2018-02-01 15:41:51 davidwang2006@aliyun.com

//DeserializeScRwObject
//反序列化 SC_FLAG为 SC_RW_OBJECT 0x03的
//我们从0x78, 0x70 之后真正开始数据的地方读取
func DeserializeScRwObject(reader io.Reader, refs []*JavaReferenceObject, className string) (JavaSerializer, error) {
	StdLogger.LevelUp()
	defer StdLogger.LevelDown()
	log.Debugf("[DeserializeScRwObject] >>\n")
	defer log.Debugf("[DeserializeScRwObject] <<\n")
	//
	switch className {
	case "java.util.HashMap", "java.util.LinkedHashMap":
		mp := &JavaHashMap{}
		log.Debugf("[DeserializeScRwObject] 读取java.util.HashMap\n")
		if err := mp.Deserialize(reader, refs); err != nil {
			return nil, err
		} else {
			return mp, nil
		}
	case "java.util.ArrayList":
		lst := &JavaArrayList{}
		log.Debugf("[DeserializeScRwObject] 读取java.util.ArrayList\n")
		if err := lst.Deserialize(reader, refs); err != nil {
			return nil, err
		} else {
			return lst, nil
		}
	case "java.util.LinkedList":
		lst := &JavaLinkedList{}
		if err := lst.Deserialize(reader, refs); err != nil {
			return nil, err
		} else {
			return lst, nil
		}
	default:
		return nil, fmt.Errorf("[DeserializeScRwObject] unexpected className %s, not be supported", className)
	}
}

//ReadNextEle
//read next map entry or list element
func ReadNextEle(reader io.Reader, refs []*JavaReferenceObject) (JavaSerializer, error) {
	StdLogger.LevelUp()
	defer StdLogger.LevelDown()
	log.Debugf("[ReadNextEle] >>\n")
	defer log.Debugf("[ReadNextEle] <<\n")
	var tp byte //type
	var err error

	if tp, err = ReadNextByte(reader); err != nil {
		return nil, err
	}

	//0x73
	log.Debugf("[ReadNextEle] type is 0x%x\n", tp)
	var js JavaSerializer
	switch tp {
	case TC_STRING:
		Log(fmt.Sprintf("%2x:	 TC_STRING.代表一个new String.用String来引用对象(ReadNextEle)。\n", tp))
		js = new(JavaTcString)
	case TC_ARRAY:
		js = &JavaTcArray{}
	case TC_OBJECT:
		log.Debugf("%2x:	 TC_OBJECT.声明这是一个新的对象(在list里面)\n", tp)
		Log(fmt.Sprintf("%2x:	 TC_OBJECT.声明这是一个新的对象(在list里面)\n", tp))
		js = &JavaTcObject{}
	case TC_REFERENCE:
		Log(fmt.Sprintf("%2x:	 TC_REFERENCE\n", tp))
		if refIndex, err := ReadUint32(reader); err != nil {
			return nil, err
		} else {
			Log(fmt.Sprintf("%2x:	 TC_REFERENCE引用序号\n", refIndex))
			ref := refs[refIndex-INTBASE_WIRE_HANDLE]
			switch ref.RefType {
			case TC_STRING:
				Log(fmt.Sprintf("%2x:	 TC_STRING.代表一个new String.用String来引用对象。\n", tp))
				if str, ok := ref.Val.(string); !ok {
					return nil, fmt.Errorf("[JavaHashMap] ref [%v] value should be string type", ref.Val)
				} else {
					Log(fmt.Sprintf("%2x:	 TC_STRING.代表一个new String.用String来引用对象。\n", tp))
					tcStr := new(JavaTcString)
					*tcStr = JavaTcString(str)
					return tcStr, nil
				}
			case TC_ARRAY, TC_OBJECT:
				if tempJs, ok := ref.Val.(JavaSerializer); !ok {
					return nil, fmt.Errorf("[JavaHashMap] ref [%v] value should be JavaSerializer type", ref.Val)
				} else {
					return tempJs, nil
				}
			default:
				return nil, fmt.Errorf("[JavaHashMap] unexpected refType 0x%x", ref.RefType)

			}

		}
	default:
		return nil, fmt.Errorf("Unexpected type 0x%x for map entry", tp)
	}

	//数组成员 对象类的描述
	Log("\n\n新对象\n\n")
	if err = js.Deserialize(reader, refs); err != nil {
		return nil, err
	}
	return js, nil

}

//SerializeScRwObject
//序列化 SC_FLAG为 SC_RW_OBJECT 0x03的
//我们从0x78, 0x70 之后真正开始数据的地方写入
func SerializeScRwObject(writer io.Writer, refs []*JavaReferenceObject, classDesc *JavaTcClassDesc) error {
	className := classDesc.ClassName
	StdLogger.LevelUp()
	defer StdLogger.LevelDown()
	log.Debugf("[SerializeScRwObject] >>\n")
	defer log.Debugf("[SerializeScRwObject] <<\n")
	//
	switch className {
	case "java.util.HashMap", "java.util.LinkedHashMap":
		mp := &JavaHashMap{
			ClassDesc: classDesc,
		}
		if err := mp.Serialize(writer, refs); err != nil {
			return err
		} else {
			return nil
		}
	case "java.util.ArrayList":
		lst := &JavaArrayList{}
		if err := lst.Serialize(writer, refs); err != nil {
			return err
		} else {
			return nil
		}
	case "java.util.LinkedList":
		lst := &JavaLinkedList{}
		if err := lst.Serialize(writer, refs); err != nil {
			return err
		} else {
			return nil
		}
	default:
		return fmt.Errorf("[SerializeScRwObject] unexpected className %s, not be supported", className)
	}
}
