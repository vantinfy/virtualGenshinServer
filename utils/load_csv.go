package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type CsvUtilMgr struct{}

var csvUtilMgr *CsvUtilMgr

func GetCsvUtilMgr() *CsvUtilMgr {
	if csvUtilMgr == nil {
		csvUtilMgr = new(CsvUtilMgr)
	}
	return csvUtilMgr
}

// LoadCsv 加载csv文件配置
func (m *CsvUtilMgr) LoadCsv(fileName string, slicePtr interface{}) error {

	filePath := "conf/csv/" + fileName + ".csv"

	res, err := m.readCsv(filePath)
	if err != nil {
		return err
	}

	return m.parseData(res, slicePtr, fileName)
}

func (m *CsvUtilMgr) readCsv(filePath string) ([][]string, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("read file err", err)
		return nil, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))

	var res [][]string
	var readFileLine int // 读取到文件的第几行

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		// 对csv文件第一行特殊处理 \ufeff表示零宽非连接符
		if readFileLine == 0 {
			for index := range line {
				if index == 0 && strings.Contains(line[index], "\ufeff") {
					line[index] = strings.Replace(line[index], "\ufeff", "", 1)
					break
				}
			}
			readFileLine++
		}
		res = append(res, line)
	}
	return res, nil
}

func (m *CsvUtilMgr) parseData(csvData [][]string, dataPtr interface{}, fileName string) error {
	outInnerType := m.getValueType(dataPtr) // outInnerType: *csvs.ConfigPlayerLevel
	//fmt.Println("out inner type", outInnerType.Elem())
	data := reflect.New(outInnerType.Elem()) // outInnerType.elem: csvs.ConfigPlayerLevel

	value := reflect.Indirect(data)   // value {0 0 0 0}
	tagMap := m.getTagMap(csvData[0]) // map[0:PlayerLevel 1:PlayerExp 2:WorldLevel 3:ChapterId]
	fieldMap, trimFlag, keyTag := m.getFieldMapSimple(value.Interface(), tagMap)
	//fmt.Println("field map:", fieldMap, "///") // 忘给结构体字段加json tag：map[:[ChapterId int]] --- 正常：map[ChapterId:[ChapterId int] PlayerExp:[PlayerExp int] PlayerLevel:[PlayerLevel int] WorldLevel:[WorldLevel int]]
	//fmt.Println("trim flag:", trimFlag, "///") // map[]	--- map[0:true 1:true 2:true 3:true]
	//fmt.Println("key tag:", keyTag, "///")     // nil --- PlayerLevel
	return m.genConfig(dataPtr, csvData, tagMap, fieldMap, trimFlag, fileName, keyTag)
}

func (m *CsvUtilMgr) getValueType(ptr interface{}) reflect.Type {
	value := reflect.ValueOf(ptr)
	if value.Kind() == reflect.Ptr {
		// 返回指针指向的值,这里是slice
		value = value.Elem()
	}

	// eg. type []*csvs.ConfigPlayerLevel 返回当前值对应的类型
	outType := value.Type()
	// eg. type *csvs.ConfigPlayerLevel 返回这些类型包含的一个值
	outInnerType := outType.Elem()

	return outInnerType
}

func (m *CsvUtilMgr) getTagMap(fields []string) map[int]string {
	tagMap := make(map[int]string, len(fields))
	for index, v := range fields {
		tagMap[index] = v
	}

	return tagMap
}

func (m *CsvUtilMgr) getFieldMapSimple(config interface{}, tagMap map[int]string) (map[string][]string, map[int]bool, string) {
	// t: csvs.ConfigPlayerLevel fields: 4
	t := reflect.TypeOf(config)
	fieldMap := make(map[string][]string, t.NumField())
	trimFlag := make(map[int]bool)
	keyTag := ""
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		trim := field.Tag.Get("trim")

		if trim == "1" || trim == "" {
			tag = m.trimNumber(tag)
			//fmt.Println("trim tag", tag)
		}

		for key, v := range tagMap {
			x := m.trimNumber(v)
			if (trim == "" || trim == "1") && tag == x {
				trimFlag[key] = true
			}
		}

		// 设置表数据结构
		for key, v := range tagMap {
			if (trim == "" || trim == "1") && tag == m.trimNumber(v) {
				tagMap[key] = tag
				break
			}
		}

		data := make([]string, 2, 2)
		data[0] = field.Name
		data[1] = fmt.Sprintf("%v", field.Type)
		fieldMap[tag] = data
		if i == 0 {
			keyTag = tag
		}
	}
	return fieldMap, trimFlag, keyTag
}

// 剪掉数字
func (m *CsvUtilMgr) trimNumber(tag string) string {
	substring := strings.TrimFunc(tag, func(r rune) bool {
		return unicode.IsNumber(r)
	})
	return substring
}

// 指针只能取interface{}, 值得取Elem,否则会蹦
// tagMap: csv第一行
// filedMap: key: tag
func (m *CsvUtilMgr) genConfig(dataPtr interface{}, csvData [][]string, tagMap map[int]string,
	fieldMap map[string][]string, trimFlag map[int]bool, fileName string, keyTag string) (err error) {
	dataVal := reflect.Indirect(reflect.ValueOf(dataPtr))
	outInnerType := m.getValueType(dataPtr)

	//if fileName == "expspeedup" {
	//litter.Dump(tagMap)
	//litter.Dump(fieldMap)
	//}

	for r := 1; r < len(csvData); r++ {
		data := reflect.New(outInnerType.Elem())

		key := 0
		for c := 0; c < len(csvData[r]); c++ {
			tag := tagMap[c]
			if _, ok := trimFlag[c]; ok {
				tag = m.trimNumber(tag)
			}

			fieldInfo, ok := fieldMap[tag]
			if !ok {
				continue
			}

			if len(fieldInfo) != 2 {
				continue
			}
			fieldName := fieldInfo[0]
			filedType := fieldInfo[1]
			cellValue := csvData[r][c]
			switch filedType {
			case "int":
				v, err := strconv.Atoi(cellValue)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", index:", r)
					break
				}
				reflect.Indirect(data).FieldByName(fieldName).SetInt(int64(v))
				if tag == keyTag && key == 0 {
					key = v
					//fmt.Println("fileName:", fileName, ", tag:", tag)
				}
			case "int64":
				v, err := strconv.ParseInt(cellValue, 10, 64)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName)
					break
				}
				reflect.Indirect(data).FieldByName(fieldName).SetInt(v)
			case "string":
				reflect.Indirect(data).FieldByName(fieldName).SetString(cellValue)
			case "[]int":
				v, err := strconv.Atoi(cellValue)
				if err != nil {
					fv, err := strconv.ParseFloat(cellValue, 64)
					if err != nil {
						fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", row:", r, ", col:", c, ", ", strings.Join(csvData[r], ","))
					} else {
						c := reflect.Indirect(data).FieldByName(fieldName)
						newSlice := reflect.Append(c, reflect.ValueOf(int(fv*10)))
						reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
					}
					break
				} else {
					c := reflect.Indirect(data).FieldByName(fieldName)
					newSlice := reflect.Append(c, reflect.ValueOf(v))
					reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
				}
			case "[]int64":
				v, err := strconv.ParseInt(cellValue, 10, 0)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName)
					break
				}
				c := reflect.Indirect(data).FieldByName(fieldName)
				newSlice := reflect.Append(c, reflect.ValueOf(v))
				reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
			case "[]string":
				c := reflect.Indirect(data).FieldByName(fieldName)
				newSlice := reflect.Append(c, reflect.ValueOf(cellValue))
				reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
			case "float32":
				fv, err := strconv.ParseFloat(cellValue, 32)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", row:", r, ", col:", c, ", ", strings.Join(csvData[r], ","))
				} else {
					reflect.Indirect(data).FieldByName(fieldName).SetFloat(fv)
				}
			case "float64":
				fv, err := strconv.ParseFloat(cellValue, 64)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", row:", r, ", col:", c, ", ", strings.Join(csvData[r], ","))
				} else {
					reflect.Indirect(data).FieldByName(fieldName).SetFloat(fv)
				}
			case "[]float32":
				fv, err := strconv.ParseFloat(cellValue, 32)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", row:", r, ", col:", c, ", ", strings.Join(csvData[r], ","))
				} else {
					c := reflect.Indirect(data).FieldByName(fieldName)
					newSlice := reflect.Append(c, reflect.ValueOf(float32(fv)))
					reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
				}
			case "[]float64":
				fv, err := strconv.ParseFloat(cellValue, 64)
				if err != nil {
					fmt.Println(err.Error(), ", fileName:"+fileName, ", fieldName:", fieldName, ", row:", r, ", col:", c, ", ", strings.Join(csvData[r], ","))
				} else {
					c := reflect.Indirect(data).FieldByName(fieldName)
					newSlice := reflect.Append(c, reflect.ValueOf(fv))
					reflect.Indirect(data).FieldByName(fieldName).Set(newSlice)
				}
			}
		}

		// [] slice
		//fmt.Println(r, "data", dataVal.Interface(), reflect.TypeOf(dataVal.Interface()).Kind())

		// dataVal[] dataVal.Type: []*csvs.ConfigPlayerLevel
		kind := reflect.TypeOf(dataVal.Interface()).Kind()
		if kind == reflect.Slice {
			dataVal.Set(reflect.Append(dataVal, data))
		} else if kind == reflect.Map {
			dataVal.SetMapIndex(reflect.ValueOf(key), data)
		}
	}

	return nil
}
