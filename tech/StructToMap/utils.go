
func metadataToMap(n report.Metadata) map[string]interface{} {
	// convert struct to map
	return utils.StructToMap(n)
}

func StructToMap[T any](c T) map[string]interface{} {
	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)

	numFields := 0
	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Tag.Get("json")
		if strings.HasSuffix(key, ",omitempty") {
			if v.Field(i).IsZero() {
				continue
			}
		}
		numFields += 1
	}

	bb := make(map[string]interface{}, numFields)

	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Tag.Get("json")
		omitempty := strings.HasSuffix(key, ",omitempty")
		if omitempty {
			if v.Field(i).IsZero() {
				continue
			}
			key = strings.TrimSuffix(key, ",omitempty")
		}
		bb[key] = v.Field(i).Interface()
	}

	return bb
}
