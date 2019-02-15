package skynetproxy

import (
    "encoding/binary"
    "fmt"
    "reflect"
    "bytes"
)

const typeNil       uint8 = 0
const typeBoolean   uint8 = 1
const typeNumber    uint8 = 2

// hibits 0 : 0 , 1: byte, 2:word, 4: dword, 6: qword, 8 : double
const typeNUmberZero    uint8 = 0
const typeNUmberByte    uint8 = 1
const typeNumberWord    uint8 = 2
const typeNumberDword   uint8 = 4
const typeNumberQword   uint8 = 6
const typeNumberReal    uint8 = 8

const typeUserdata      uint8 = 3
const typeShortString   uint8 = 4
const typeLongString    uint8 = 5
const typeTable         uint8 = 6

const maxDepth          uint8 = 32
const maxCookie         int = 32

const packHeadLen      int = 4

func PackProxyMsg(pack []byte) []byte {
    packLen := len(pack)
    buf := make([]byte, packLen + packHeadLen)
    binary.LittleEndian.PutUint32(buf, uint32(packLen))
    copy(buf[packHeadLen:], pack)
    return buf
}

func UnPackProxyMsg(msg []byte) []byte {
    msgLen := len(msg)
    buf := make([]byte, msgLen - packHeadLen)
    copy(buf, msg[packHeadLen:])
    return buf
}

func DataSize(msg []byte) uint32 {
    return binary.LittleEndian.Uint32(msg)
}

// comm
func combineType(t uint8, v uint8) uint8 {
    return t | v << 3
}

// pack msg
func PackMsg(params ...interface{}) (bool,[]byte) {
    var b bytes.Buffer
    for _, param := range params {
        ok := pack(param,&b)
        if !ok {
            fmt.Println("pack error! please check your params")
            return ok,b.Bytes()
        }
    }

    return true,b.Bytes()
}

func pack(param interface{}, b *bytes.Buffer) (bool)  {
    paramType := reflect.TypeOf(param)
    if paramType.Kind() == reflect.Ptr {
        // 暂时不支持指针参数
        fmt.Println("error! param have ptr")
        return false
    } else {
        return packUnPtr(param,b)
    }

    return false
}

func packUnPtr(param interface{}, b *bytes.Buffer) (bool) {
    paramType := reflect.TypeOf(param)
    paramKind := paramType.Kind()
    switch paramKind {
    case reflect.Bool:
        return packOneBool(param.(bool),b)
    case reflect.Struct:
        return packOneStruct(param,b)
    case reflect.Int:
        return packOneInter(int64(param.(int)),b)
    case reflect.Int8:
        return packOneInter(int64(param.(int8)),b)
    case reflect.Int16:
        return packOneInter(int64(param.(int16)),b)
    case reflect.Int32:
        return packOneInter(int64(param.(int32)),b)
    case reflect.Int64:
        return packOneInter(param.(int64),b)
    case reflect.Uint:
        return packOneInter(int64(param.(uint)),b)
    case reflect.Uint8:
        return packOneInter(int64(param.(uint8)),b)
    case reflect.Uint16:
        return packOneInter(int64(param.(uint16)),b)
    case reflect.Uint32:
        return packOneInter(int64(param.(uint32)),b)
    case reflect.Uint64:
        return packOneInter(int64(param.(uint64)),b)
    case reflect.Float32:
        return packOneReal(float64(param.(float32)),b)
    case reflect.Float64:
        return packOneReal(param.(float64),b)
    case reflect.String:
        return packOneString(param.(string),b)
    default:
        fmt.Printf("Unknown kind = %d.\n",paramKind)
    }

    return false
}

func writeBytrReal(num float64, b *bytes.Buffer) {
    binary.Write(b,binary.LittleEndian,num)
}

func writeBytr64(num int64, b *bytes.Buffer) {
    binary.Write(b,binary.LittleEndian,num)
}

func writeBytr32(num int32, b *bytes.Buffer) {
    binary.Write(b,binary.LittleEndian,num)
}

func writeBytrU8(num uint8, b *bytes.Buffer) {
    binary.Write(b,binary.LittleEndian,num)
}

func writeBytrU16(num uint16, b *bytes.Buffer) {
    binary.Write(b,binary.LittleEndian,num)
}

func writeBytrU32(num uint32, b *bytes.Buffer) {
    binary.Write(b,binary.LittleEndian,num)
}

func packOneNil(b *bytes.Buffer) (bool) {
    var n uint8 = typeNil;
    b.WriteByte(n)
    return true
}

func packOneBool(r bool, b *bytes.Buffer) (bool) {
    var n uint8;
    if r {
        n = combineType(typeBoolean, 1)
    } else {
        n = combineType(typeBoolean, 0)
    }

    b.WriteByte(n)
    fmt.Printf("packOneBool:n[%d]\n",n)
    return true
}

func packOneInter(num int64, b *bytes.Buffer) (bool) {
    if num == 0 {
        var n uint8 = combineType(typeNumber, typeNUmberZero);
        b.WriteByte(n)
        //fmt.Printf("packOneInter:%d\n",num)
    } else if (num > (1 << 31 - 1)) || (num < 0 - (1 << 31 - 1)) {
        //64bit
        var n uint8 = combineType(typeNumber, typeNumberQword)
        b.WriteByte(n)
        writeBytr64(num,b)
        //fmt.Printf("packOneInter(64bit):num[%d],n[%d]\n",num,n)
    } else if num < 0 {
        var n uint8 = combineType(typeNumber, typeNumberDword)
        b.WriteByte(n)
        writeBytr32(int32(num),b)
        //fmt.Printf("packOneInter(num<0):num[%d],n[%d]\n",num,n)
    } else if num < 0x100 {
        var n uint8 = combineType(typeNumber, typeNUmberByte)
        b.WriteByte(n)
        writeBytrU8(uint8(num),b)
        //fmt.Printf("packOneInter(num<0x100):num[%d],n[%d]\n",num,n)
    } else if num < 0x10000 {
        var n uint8 = combineType(typeNumber, typeNumberWord)
        b.WriteByte(n)
        writeBytrU16(uint16(num),b)
        //fmt.Printf("packOneInter(num<0x10000):num[%d],n[%d]\n",num,n)
    } else {
        var n uint8 = combineType(typeNumber, typeNumberDword)
        b.WriteByte(n)
        writeBytrU32(uint32(num),b)
        //fmt.Printf("packOneInter(num32):num[%d],n[%d]\n",num,n)
    }

    return true
}

func packOneReal(num float64, b *bytes.Buffer) (bool) {
    var n uint8 = combineType(typeNumber , typeNumberReal);
    b.WriteByte(n)
    writeBytrReal(num,b)
    return true
}

func packOneString(s string, b *bytes.Buffer) (bool) {
    nLen := len(s) 
    if nLen < maxCookie {
        var n uint8 = combineType(typeShortString, uint8(nLen))
        b.WriteByte(n)
        if nLen > 0 {
            b.WriteString(s)
        }
    } else {
        var n uint8
        if nLen < 0x10000 {
            n = combineType(typeLongString, 2)
            b.WriteByte(n)
            writeBytrU16(uint16(nLen),b)
        } else {
            n = combineType(typeLongString, 4)
            b.WriteByte(n)
            writeBytrU32(uint32(nLen),b)
        }
        b.WriteString(s)
    }
    return true
}

func packOneStruct(param interface{},b *bytes.Buffer) (bool) {
    t := reflect.TypeOf(param)
    v := reflect.ValueOf(param)
    if t.Kind() == reflect.Ptr { 
        t = t.Elem() 
    }

    fieldNum := t.NumField() 
    
    var n uint8 = combineType(typeTable, 0)
    b.WriteByte(n)

    for i:= 0;i<fieldNum;i++ {
        packUnPtr(string(t.Field(i).Tag),b)
        packUnPtr(v.Field(i).Interface(),b)
    }

    packOneNil(b)
    return true
}

// unpack msg
func readType(b *bytes.Buffer) (uint8,error) {
    c ,err := b.ReadByte()
    if err != nil {
        fmt.Println("bytes.Buffer read error")
        return 0,err
    }

    return c,nil
}

func UnPackMsg(msg []byte) []interface{} {
	var b bytes.Buffer
	b.Write(msg)
    var retlist [] interface{}
    for i:= 1;;i++ {
        if b.Len() <= 0 {
            break
        }
        addInterface(&retlist,unpackOne(&b))
	}
	
	return retlist
}

func addInterface(ret *[]interface{}, param interface{}) {
    if param == nil {
		fmt.Println(*ret)
        fmt.Println("addInterface param is nil")
        return
    }

    *ret = append(*ret, param)
}

func unpackOne(b *bytes.Buffer) interface{} {
    t,err := readType(b)
    if err != nil {
        fmt.Println("unpackOne err! data len is not enough")
        return nil
    }

    return unpackValue(b, t & 0x7, t >> 3)
}

func getInter(b *bytes.Buffer, cookie uint8) (int) {
    param := unpackInter(b,cookie)
    paramType := reflect.TypeOf(param)
    paramKind := paramType.Kind()
    switch paramKind {
    case reflect.Int:
        return param.(int)
    case reflect.Int32:
        return int(param.(int32))
    case reflect.Int64:
        return int(param.(int64))
    case reflect.Uint8:
        return int(param.(uint8))
    case reflect.Uint16:
        return int(param.(uint16))
    default:
        fmt.Printf("(getInter)Unknown kind = %d.\n",paramKind)
    }

    return 0
}

func unpackInter(b *bytes.Buffer, cookie uint8) interface{} {
    switch cookie {
    case typeNUmberZero:
        //b.ReadByte();
        //fmt.Println("unpackInter(typeNUmberZero)")
        return 0
    case typeNumberQword:
        var num int64
        binary.Read(b,binary.LittleEndian,&num)
        //fmt.Printf("unpackInter(typeNumberQword) num = %d.\n",num)
        return num
    case typeNumberDword:
        var num int32
        binary.Read(b,binary.LittleEndian,&num)
        //fmt.Printf("unpackInter(typeNumberDword) num = %d.\n",num)
        return num;
    case typeNUmberByte:
        var num uint8
        binary.Read(b,binary.LittleEndian,&num)
        //fmt.Printf("unpackInter(typeNUmberByte) num = %d.\n",num)
        return num;
    case typeNumberWord:
        var num uint16
        binary.Read(b,binary.LittleEndian,&num)
        //fmt.Printf("unpackInter(typeNumberWord) num = %d.\n",num)
        return num
    default:
        //fmt.Println("unKnown cookie = %d.\n", cookie)
        return 0
    }
}

func unpackString(b *bytes.Buffer, len uint32) string {
    strBuf := make([]byte, len)
    n, err := b.Read(strBuf)
    if err != nil || n != int(len) {
        fmt.Printf("bytes.Buffer readstring err! len = %d, readN = %d.\n",len,n)
    }

    return string(strBuf)
}

func unpackTable(b *bytes.Buffer, arryLen int) interface{} {
    if (arryLen == maxCookie - 1) {
        t,err := readType(b)
        if err != nil {
            fmt.Println("unpackTable err! data len is not enough.")
            return nil
        }
        
        cookie := t >> 3
        if (t & 7) != typeNumber || cookie == typeNumberReal {
            fmt.Println("unpackTable err! type is match.")
            return nil
        }

        arryLen = getInter(b,cookie)
    }

    tb := make(map[interface{}]interface{})
    for i := 1; i <= arryLen; i++ {
        tb[i] = unpackOne(b)
    }

    for i :=1 ;;i++ {
        key := unpackOne(b)
        if key == nil {
            break
        }
		value := unpackOne(b)
        tb[key] = value
    }

    return tb
}

func unpackValue(b *bytes.Buffer, t uint8, cookie uint8) interface{} {
    //fmt.Printf("unpackValue type = %d,cookie = %d.\n", t,cookie)
    switch t {
    case typeNil:
        fmt.Printf("unpackValue typeNil.\n")
    case typeBoolean:
        if cookie == 1 {
            return true
            //*retList = append(*retList, true)
        }else {
            return false
        }
    case typeNumber:
        if cookie == typeNumberReal {
            var num float64
            binary.Read(b,binary.LittleEndian,&num)
            return num
        } else {
            return unpackInter(b,cookie)
        }
    case typeShortString:
        return unpackString(b,uint32(cookie))
        fmt.Println("unpackValue typeShortString.")
    case typeLongString:
        if cookie == 2 {
            var num uint16
            binary.Read(b,binary.LittleEndian,&num)
            return unpackString(b,uint32(num))
            //*retList = append(*retList, unpackString(b,uint32(num)))
        } else {
            var num uint32
            binary.Read(b,binary.LittleEndian,&num)
            return num
            //*retList = append(*retList, unpackString(b,num))
        }
    case typeTable:
        return unpackTable(b, int(cookie))
        //unpackTable(retList, b, int(cookie))
        fmt.Printf("table type = %d.\n",t)
    default:
        fmt.Printf("unpackValue Unknown type = %d.\n",t)
    }

    return nil
}