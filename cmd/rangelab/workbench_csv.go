package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type rangeWorkbenchCSVField struct {
	name  string
	index []int
}

func writeRangeWorkbenchRowsCSV(path string, rows any) error {
	value := reflect.ValueOf(rows)
	if value.Kind() != reflect.Slice {
		return fmt.Errorf("range workbench csv rows must be a slice")
	}
	fields := rangeWorkbenchCSVFields(value.Type().Elem())
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	header := []string{}
	for _, field := range fields {
		header = append(header, field.name)
	}
	if err := writer.Write(header); err != nil {
		return err
	}
	for i := 0; i < value.Len(); i++ {
		record := []string{}
		rowValue := value.Index(i)
		for _, field := range fields {
			record = append(record, fmt.Sprint(rowValue.FieldByIndex(field.index).Interface()))
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func rangeWorkbenchCSVFields(t reflect.Type) []rangeWorkbenchCSVField {
	fields := []rangeWorkbenchCSVField{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			for _, sub := range rangeWorkbenchCSVFields(field.Type) {
				fields = append(fields, rangeWorkbenchCSVField{name: sub.name, index: append(field.Index, sub.index...)})
			}
			continue
		}
		name := field.Name
		if tag := field.Tag.Get("json"); tag != "" {
			name = strings.Split(tag, ",")[0]
		}
		if name == "" || name == "-" {
			continue
		}
		fields = append(fields, rangeWorkbenchCSVField{name: name, index: field.Index})
	}
	return fields
}
