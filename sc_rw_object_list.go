package main

import "io"
import "fmt"

//JavaArrayList
type JavaArrayList struct {
	Size int
	Eles []interface{}
}

//Deserialize
func (arrList *JavaArrayList) Deserialize(reader io.Reader, refs []*JavaReferenceObject) error {
	StdLogger.LevelUp()
	defer StdLogger.LevelDown()
	StdLogger.Debug("[JavaArrayList] >>\n")
	defer StdLogger.Debug("[JavaArrayList] <<\n")
	//start with size
	if ui, err := ReadUint32(reader); err != nil {
		return err
	} else {
		arrList.Size = int(ui)
	}
	Log("\n\n=================Blockdata数据块起始位置=================\n\n")
	Log(fmt.Sprintf("%2x:	Blockdata数据块大小\n", arrList.Size))
	StdLogger.Debug("[JavaArrayList] >> Size:%d\n", arrList.Size)
	//TC_BLOCKDATA
	if b, err := ReadNextByte(reader); err != nil {
		return err
	} else if b != TC_BLOCKDATA {
		return fmt.Errorf("There should be TC_BLOCKDATA, but got 0x%x", b)
	} else {
		Log(fmt.Sprintf("%2x:	TcBlockdata数据块标记\n", b))
	}

	//should follow by 0x04, 表示4字节后是所有的elements --> header长度
	if b, err := ReadNextByte(reader); err != nil {
		return err
	} else if b != 0x04 {
		return fmt.Errorf("There should be 0x04, but got 0x%x", b)
	} else {
		Log(fmt.Sprintf("%2x:	固定为0x04\n", b))
	}

	//数组元素的个数
	if ui, err := ReadUint32(reader); err != nil {
		return err
	} else if arrList.Size != int(ui) {
		return fmt.Errorf("Size should be %d, but got %d", arrList.Size, ui)
	} else {

		Log(fmt.Sprintf("%2x:	数组元素的个数\n", ui))
	}

	//now it's the data
	arrList.Eles = make([]interface{}, 0, arrList.Size)

	for i := 0; i < arrList.Size; i += 1 {
		StdLogger.Debug("[JavaArrayList] >> 读取次数:%d\n", i+1)

		Log(fmt.Sprintf("\n\n======================================================读取第%d个数据块======================================================\n\n", i+1))
		if ele, err := ReadNextEle(reader, refs); err != nil {
			StdLogger.Error("[JavaArrayList] Error when read %d element: %v\n", i, err)
			return err
		} else {
			arrList.Eles = append(arrList.Eles, ele.JsonMap())
		}
	}
	//TC_ENDBLOCKDATA
	//must be 0x78 TC_ENDBLOCKDATA
	if b, err := ReadNextByte(reader); err != nil {
		return err
	} else if b != TC_ENDBLOCKDATA {
		return fmt.Errorf("[Deserialize] There should be TC_ENDBLOCKDATA, but got 0x%x", b)
	} else {
		Log(fmt.Sprintf("%2x:	TC_ENDBLOCKDATA,对象数据块结束的标志\n", b))
	}

	return nil
}

//JsonMap just return list's elements
func (arrList *JavaArrayList) JsonMap() interface{} {
	return arrList.Eles
}

//JavaLinkedList
type JavaLinkedList struct {
	Size int
	Eles []interface{}
}

//Deserialize
func (linkedList *JavaLinkedList) Deserialize(reader io.Reader, refs []*JavaReferenceObject) error {
	StdLogger.LevelUp()
	defer StdLogger.LevelDown()
	StdLogger.Debug("[JavaLinkedList] >>\n")
	defer StdLogger.Debug("[JavaLinkedList] <<\n")
	//TC_BLOCKDATA
	if b, err := ReadNextByte(reader); err != nil {
		return err
	} else if b != TC_BLOCKDATA {
		return fmt.Errorf("There should be TC_BLOCKDATA, but got 0x%x", b)
	}
	//should follow by 0x04, 表示4字节后是所有的elements
	if b, err := ReadNextByte(reader); err != nil {
		return err
	} else if b != 0x04 {
		return fmt.Errorf("There should be 0x04, but got 0x%x", b)
	}
	if ui, err := ReadUint32(reader); err != nil {
		return err
	} else {
		linkedList.Size = int(ui)
	}
	//now it's the data
	linkedList.Eles = make([]interface{}, 0, linkedList.Size)

	for i := 0; i < linkedList.Size; i += 1 {

		if ele, err := ReadNextEle(reader, refs); err != nil {
			StdLogger.Error("[JavaLinkedList] Error when read %d element: %v\n", i, err)
			return err
		} else {
			linkedList.Eles = append(linkedList.Eles, ele.JsonMap())
		}
	}
	//TC_ENDBLOCKDATA
	//must be 0x78 TC_ENDBLOCKDATA
	if b, err := ReadNextByte(reader); err != nil {
		return err
	} else if b != TC_ENDBLOCKDATA {
		return fmt.Errorf("There should be TC_ENDBLOCKDATA, but got 0x%x", b)
	} else {
		Log(fmt.Sprintf("%2x:	TC_ENDBLOCKDATA,对象数据块结束的标志\n", b))
	}

	return nil
}

//JsonMap just return list's elements
func (linkedList *JavaLinkedList) JsonMap() interface{} {
	return linkedList.Eles
}

func (linkedList *JavaLinkedList) Serialize(writer io.Writer, refs []*JavaReferenceObject) error {
	return fmt.Errorf("to be continued....")
}
func (arrayList *JavaArrayList) Serialize(writer io.Writer, refs []*JavaReferenceObject) error {

	return fmt.Errorf("to be continued....")
}
